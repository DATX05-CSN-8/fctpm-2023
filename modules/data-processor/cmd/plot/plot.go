package main

import (
	"flag"
	"image/color"
	"math/rand"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/data-processor/internal/csvreader"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/data-processor/internal/inputdataparser"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	var inputDatas inputdataparser.InputDataList
	flag.Var(&inputDatas, "input-data", "Input data to be used for the calculation. Specify as [Name of series],[File path0]")
	outputtype := flag.String("output-type", "png", "Output file type (png, tex)")
	flag.Parse()

	p := plot.New()
	rand.Seed(int64(0))
	p.X.Label.Text = "Boot time (ms)"
	p.Y.Label.Text = "CDF"
	var min float64 = 9 * 10_000
	var max float64 = 0
	extraSize := 10
	dashstyles := [][]vg.Length{{}, {vg.Points(7), vg.Points(2)}, {vg.Points(2)}}
	colors := []color.Color{
		color.RGBA{R: 47, G: 134, B: 142, A: 255},
		color.RGBA{R: 95, G: 134, B: 45, A: 255},
		color.RGBA{R: 157, G: 63, B: 52, A: 255},
	}

	for i, d := range inputDatas {
		parsed, err := csvreader.ReadCsvFile(d.File, csvreader.OutputDataRowParser)
		if err != nil {
			panic(err)
		}

		deltas := csvreader.CalculateDeltas(parsed)

		fn := plotter.NewFunction(func(x float64) float64 {
			var sz float64 = 0
			for _, d := range deltas {
				if d > x {
					break
				}
				sz += 1
			}
			return sz / float64(len(deltas))
		})
		if deltas[0] < min {
			min = deltas[0]
		}
		last := deltas[len(deltas)-1]
		if last > max {
			max = last
		}
		fn.LineStyle.Dashes = dashstyles[i%len(dashstyles)]
		fn.LineStyle.Width = vg.Points(2)
		fn.Color = colors[i%len(colors)]
		p.Add(fn)
		p.Legend.Add(d.Name, fn)
		fn.XMin = deltas[0]
		fn.XMax = last
	}
	grid := plotter.NewGrid()
	grid.Vertical.Color = color.RGBA{R: 0, G: 0, B: 0, A: 35}
	grid.Horizontal.Color = grid.Vertical.Color
	p.Add(grid)
	p.X.Tick.Marker = MyTicker{width: 10, labelwidth: 100, labelfmt: "%.f"}
	p.Y.Tick.Marker = MyTicker{width: 0.05, labelwidth: 0.2, labelfmt: "%.1f"}
	p.X.Min = min - float64(extraSize)
	p.X.Max = max + float64(extraSize)
	p.Y.Min = 0
	p.Y.Max = 1.1

	err := p.Save(4*vg.Inch, 4*vg.Inch, "out."+*outputtype)
	if err != nil {
		panic(err)
	}
}
