package csvreader

import "strconv"

type MemoverheadRowData struct {
	MemSpec     float64
	MemOverhead float64
}

func CreateMemOverHeadParser(actualColumnIdx int) RowParser[MemoverheadRowData] {
	return func(row []string) (*MemoverheadRowData, error) {
		spec, err := strconv.ParseFloat(row[0], 64)
		if err != nil {
			panic(err)
		}
		actual, err := strconv.ParseFloat(row[actualColumnIdx], 64)
		if err != nil {
			panic(err)
		}
		memspec := spec / 1024
		overhead := (actual / 1024) - (spec / 1024)
		return &MemoverheadRowData{
			MemSpec:     memspec,
			MemOverhead: overhead,
		}, nil
	}
}
