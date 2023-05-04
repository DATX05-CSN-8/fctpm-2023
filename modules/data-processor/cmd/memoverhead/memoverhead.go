package main

import (
	"flag"
	"image/color"
	"math/rand"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/data-processor/internal/csvreader"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/data-processor/internal/inputdataparser"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)

func main() {
	var inputDatas inputdataparser.InputDataList
	flag.Var(&inputDatas, "input-data", "Input data to be used for the calculation. Specify as [Name of series],[File path0]")
	outputtype := flag.String("output-type", "png", "Output file type (png, tex)")
	flag.Parse()

	p := plot.New()
	rand.Seed(int64(0))
	p.X.Label.Text = "VM Memory (MB)"
	p.Y.Label.Text = "Overhead (MB)"
	var min float64 = 0
	var max float64 = 8192
	extraSize := 10

	baseline, err := csvreader.ReadCsvFile(inputDatas[0].File, csvreader.CreateMemOverHeadParser(1))
	if err != nil {
		panic(err)
	}
	noSwtpm, err := csvreader.ReadCsvFile(inputDatas[1].File, csvreader.CreateMemOverHeadParser(1))
	if err != nil {
		panic(err)
	}
	wstpm, err := csvreader.ReadCsvFile(inputDatas[1].File, csvreader.CreateMemOverHeadParser(3))
	if err != nil {
		panic(err)
	}

	plotutil.AddLinePoints(p, "baseline", toXys(baseline), "Modified without swtpm", toXys(noSwtpm), "Modified with swtpm", toXys(wstpm))
	grid := plotter.NewGrid()
	grid.Vertical.Color = color.RGBA{R: 0, G: 0, B: 0, A: 35}
	grid.Horizontal.Color = grid.Vertical.Color
	p.Add(grid)
	p.Legend.YOffs = vg.Inch
	// p.X.Tick.Marker = MyTicker{width: 10, labelwidth: 100, labelfmt: "%.f"}
	// p.Y.Tick.Marker = MyTicker{width: 0.05, labelwidth: 0.2, labelfmt: "%.1f"}
	p.X.Min = min - float64(extraSize)
	p.X.Max = max + float64(extraSize)
	p.Y.Min = 4
	p.Y.Max = 15

	err = p.Save(4*vg.Inch, 4*vg.Inch, "out."+*outputtype)
	if err != nil {
		panic(err)
	}
}

func toXys(data []*csvreader.MemoverheadRowData) plotter.XYs {
	pts := make(plotter.XYs, len(data))
	for i := range pts {
		pts[i].X = data[i].MemSpec
		pts[i].Y = data[i].MemOverhead
	}
	return pts
}
