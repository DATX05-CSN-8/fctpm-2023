package main

import (
	"flag"
	"fmt"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/data-processor/internal/csvreader"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/data-processor/internal/csvwriter"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/data-processor/internal/inputdataparser"
	"gonum.org/v1/gonum/stat"
)

func main() {
	var inputDatas inputdataparser.InputDataList
	flag.Var(&inputDatas, "input-data", "Input data to be used for the calculation. Specify as [Name of series],[File path0]")
	outpath := flag.String("outpath", "./out.csv", "Path to output file")
	flag.Parse()

	tabledata := make([][]string, 1)
	tabledata[0] = []string{"Type", "Min", "Max", "Average", "p95", "p95 Index", "Standard Deviation"}
	var baselineP95 float64
	for i, d := range inputDatas {
		parsed, err := csvreader.ReadCsvFile(d.File, csvreader.OutputDataRowParser)
		if err != nil {
			panic(err)
		}

		deltas := csvreader.CalculateDeltas(parsed)
		min := deltas[0]
		max := deltas[len(deltas)-1]
		avg := stat.Mean(deltas, nil)
		p95 := stat.Quantile(0.95, stat.Empirical, deltas, nil)
		var p95idx float64
		if i == 0 {
			p95idx = 1
			baselineP95 = p95
		} else {
			p95idx = p95 / baselineP95
		}
		std := stat.StdDev(deltas, nil)
		format := "%.2f"
		tabledata = append(tabledata, []string{
			d.Name,
			fmt.Sprintf(format, min),
			fmt.Sprintf(format, max),
			fmt.Sprintf(format, avg),
			fmt.Sprintf(format, p95),
			fmt.Sprintf(format, p95idx),
			fmt.Sprintf(format, std),
		})
	}
	if err := csvwriter.WriteCsvFile(*outpath, &tabledata); err != nil {
		panic(err)
	}
}
