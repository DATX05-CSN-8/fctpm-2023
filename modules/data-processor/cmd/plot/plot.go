package main

import (
	"flag"
	"image/color"
	"math/rand"
	"sort"

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

	for _, d := range inputDatas {
		parsed, err := csvreader.ReadCsvFile(d.File, csvreader.OutputDataRowParser)
		if err != nil {
			panic(err)
		}

		deltas := make([]float64, len(parsed))
		for i, d := range parsed {
			delta := float64(d.End.Sub(d.Start).Microseconds()) / 1000
			deltas[i] = delta
		}
		sort.Float64s(deltas)

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
		fn.Color = color.YCbCr{
			Y:  127,
			Cb: uint8(rand.Int()),
			Cr: uint8(rand.Int()),
		}
		p.Add(fn)
		p.Legend.Add(d.Name, fn)
		fn.XMin = deltas[0]
		fn.XMax = last
	}
	grid := plotter.NewGrid()
	p.Add(grid)
	p.X.Tick.Marker = MyTicker{width: 5, labelwidth: 20, labelfmt: "%.f"}
	p.Y.Tick.Marker = MyTicker{width: 0.05, labelwidth: 0.2, labelfmt: "%.1f"}
	p.X.Min = min - float64(extraSize)
	p.X.Max = max + float64(extraSize)
	p.Y.Min = 0
	p.Y.Max = 1

	err := p.Save(4*vg.Inch, 4*vg.Inch, "out."+*outputtype)
	if err != nil {
		panic(err)
	}
}
