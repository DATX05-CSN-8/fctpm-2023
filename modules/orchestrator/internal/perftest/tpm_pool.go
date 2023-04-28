package perftest

import (
	"fmt"
	"strconv"
)

type BufferPool struct {
	newq []string
}

func NewTpmPool(resultpath *string, size int) *BufferPool {

	s := BufferPool{
		newq: make([]string, size),
	}

	for i := 0; i < len(s.newq); i++ {
		str := "test buffer " + strconv.Itoa(i+1)
		s.newq[i] = str
		fmt.Println(str)
	}
	return &s
}
