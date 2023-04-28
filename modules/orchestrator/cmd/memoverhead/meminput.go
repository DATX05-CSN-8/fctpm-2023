package main

import (
	"fmt"
	"strconv"
	"strings"
)

type MemInput []int

func (i *MemInput) String() string {
	return fmt.Sprintf("Size: %d", len(*i))
}

func (i *MemInput) Set(value string) error {
	splitted := strings.Split(value, ",")
	for _, s := range splitted {
		parsed, err := strconv.Atoi(s)
		if err != nil {
			return err
		}
		*i = append(*i, parsed)
	}
	return nil
}
