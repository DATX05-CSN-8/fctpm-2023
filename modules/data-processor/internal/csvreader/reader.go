package csvreader

import (
	"encoding/csv"
	"io"
	"os"
)

type RowParser[R any] func([]string) (*R, error)

func ReadCsvFile[R any](path string, rowparser RowParser[R]) ([]*R, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)
	var d []*R
	// skip title row
	reader.Read()
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}
		parsed, err := rowparser(row)
		if err != nil {
			return nil, err
		}
		d = append(d, parsed)
	}
	return d, nil
}
