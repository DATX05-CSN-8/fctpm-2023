package csvwriter

import (
	"encoding/csv"
	"os"
)

func WriteCsvFile(path string, data *[][]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	csvw := csv.NewWriter(f)
	err = csvw.WriteAll(*data)
	if err != nil {
		return err
	}
	return nil
}
