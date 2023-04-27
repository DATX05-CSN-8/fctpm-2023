package inputdataparser

import (
	"fmt"
	"strings"
)

type InputData struct {
	Name string
	File string
}

type InputDataList []InputData

func (i *InputDataList) String() string {
	var sb strings.Builder
	for _, id := range *i {
		sb.WriteString(fmt.Sprintf("Name: %s, File: %s\n", id.Name, id.File))
	}
	return sb.String()
}

func (i *InputDataList) Set(value string) error {
	splitted := strings.Split(value, ",")
	if len(splitted) != 2 {
		return fmt.Errorf("Please put two values separated by a comma")
	}
	rawName := splitted[0]
	var name string
	if rawName[0] == '"' {
		name = rawName[1 : len(rawName)-1]
	} else {
		name = rawName
	}
	path := splitted[1]
	*i = append(*i, InputData{
		Name: name,
		File: path,
	})
	return nil
}
