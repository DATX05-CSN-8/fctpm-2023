package pmap

import (
	"os"
	"reflect"
	"testing"
)

func TestParsePmap(t *testing.T) {
	// given
	inputfilepath := "testdata/pmapout1"
	inputbytes, err := os.ReadFile(inputfilepath)
	if err != nil {
		t.Error("Could not open fixture", err)
	}
	inputdata := string(inputbytes)
	expected := []ParsedRow{
		{
			kbytes: 184,
			mode:   1 << MODE_READ,
		},
		{
			kbytes: 748,
			mode:   1<<MODE_READ | 1<<MODE_EXEC,
		},
		{
			kbytes: 224,
			mode:   1 << MODE_READ,
		},
		{
			kbytes: 12,
			mode:   1 << MODE_READ,
		},
		{
			kbytes: 36,
			mode:   1<<MODE_READ | 1<<MODE_WRITE,
		},
		{
			kbytes: 44,
			mode:   1<<MODE_READ | 1<<MODE_WRITE,
		},
		{
			kbytes: 2644,
			mode:   1<<MODE_READ | 1<<MODE_WRITE,
		},
	}
	// when
	parsed, err := ParseOutput(&inputdata)
	if err != nil {
		t.Error("Error parsing data", err)
	}
	// then
	if !reflect.DeepEqual(*parsed, expected) {
		t.Error("The outputs were not equal", *parsed, expected)
	}
}
