package csvreader

import "sort"

func CalculateDeltas(data []*RowData) []float64 {
	deltas := make([]float64, len(data))
	for i, d := range data {
		delta := float64(d.End.Sub(d.Start).Microseconds()) / 1000
		deltas[i] = delta
	}
	sort.Float64s(deltas)
	return deltas
}
