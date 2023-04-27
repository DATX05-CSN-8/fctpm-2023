package csvreader

import (
	"fmt"
	"strconv"
	"time"
)

type RowData struct {
	Idx   int
	Start time.Time
	Exec  time.Duration
	End   time.Time
}

func parseTime(val string) (*time.Time, error) {
	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		return nil, err
	}
	t := time.UnixMicro(num)
	return &t, nil
}

func OutputDataRowParser(row []string) (*RowData, error) {
	if len(row) != 4 {
		return nil, fmt.Errorf("Row length was not 4")
	}
	idx, err := strconv.Atoi(row[0])
	if err != nil {
		return nil, err
	}

	start, err := parseTime(row[1])
	if err != nil {
		return nil, err
	}
	var exec time.Duration
	execInt, err := strconv.ParseInt(row[2], 10, 64)
	exec = time.Duration(execInt) * time.Microsecond
	if err != nil {
		return nil, err
	}
	end, err := parseTime(row[3])
	if err != nil {
		return nil, err
	}
	val := RowData{
		Idx:   idx,
		Start: *start,
		Exec:  exec,
		End:   *end,
	}
	return &val, nil
}
