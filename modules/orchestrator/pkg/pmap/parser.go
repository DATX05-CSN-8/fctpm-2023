package pmap

import (
	"bufio"
	"strconv"
	"strings"
)

type ParsedRow struct {
	kbytes int
	mode   int
}
type Parsed = []ParsedRow

const (
	MODE_READ int = iota
	MODE_WRITE
	MODE_EXEC
)

func ParseOutput(output *string) (*Parsed, error) {
	scanner := bufio.NewScanner(strings.NewReader(*output))
	var data []ParsedRow
	for scanner.Scan() {
		text := scanner.Text()
		if strings.HasPrefix(text, "Address") {
			continue
		}
		if strings.HasPrefix(text, "total") {
			continue
		}
		if strings.HasPrefix(text, "-----") {
			continue
		}
		if strings.Contains(text, ":") {
			continue
		}
		// now it is a valid row
		splitted := strings.Fields(text)
		mode := splitted[4]
		kbytes := splitted[1]
		kbytesVal, err := strconv.ParseInt(kbytes, 10, 32)
		if err != nil {
			return nil, err
		}
		modeVal, err := parseMode(&mode)
		if err != nil {
			return nil, err
		}
		data = append(data, ParsedRow{
			kbytes: int(kbytesVal),
			mode:   modeVal,
		})
	}
	return &data, nil
}

func parseMode(modeStr *string) (int, error) {
	var mode int = 0
	if strings.HasPrefix(*modeStr, "r") {
		mode = 1 << MODE_READ
	}
	if strings.HasPrefix((*modeStr)[1:], "w") {
		mode |= 1 << MODE_WRITE
	}
	if strings.HasPrefix((*modeStr)[2:], "x") {
		mode |= 1 << MODE_EXEC
	}
	return mode, nil
}

func (r *ParsedRow) Kbytes() int {
	return r.kbytes
}

func (r *ParsedRow) Writeable() bool {
	return (r.mode & (1 << MODE_WRITE)) > 0
}
