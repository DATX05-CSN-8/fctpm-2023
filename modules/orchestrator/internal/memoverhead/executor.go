package memoverhead

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/pmap"
)

type memoryOverheadExecutor struct {
	runner     MemoryOverheadRunner
	resultpath *string
}

type MemoryOverheadRunner interface {
	Run(int) (instance, error)
	Stop(instance) error
}

type instance interface {
	Processes() map[string]*os.Process
	Cleanup() error
}

func NewMemoryOverheadExecutor(runner MemoryOverheadRunner, resultpath *string) (*memoryOverheadExecutor, error) {
	return &memoryOverheadExecutor{
		runner:     runner,
		resultpath: resultpath,
	}, nil
}

func (e *memoryOverheadExecutor) RunWithMems(sizes []int) error {
	if len(sizes) < 1 {
		return fmt.Errorf("Number of sizes need to be greater than 0")
	}
	var data []map[string]float64
	for _, size := range sizes {
		result, err := e.runMem(size)
		if err != nil {
			return err
		}
		data = append(data, result)
	}
	var keys []string
	for k := range data[0] {
		keys = append(keys, k)
	}

	file, err := os.Create(*e.resultpath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	withExtrakeys := append([]string{"mem"}, append(keys, "total")...)
	err = writer.Write(withExtrakeys)
	if err != nil {
		return err
	}
	for i, d := range data {
		var towrite []string
		towrite = append(towrite, fmt.Sprintf("%d", sizes[i]))
		sum := 0.0
		for _, key := range keys {
			val := fmt.Sprintf("%.0f", d[key])
			towrite = append(towrite, val)
			sum += d[key]
		}
		towrite = append(towrite, fmt.Sprintf("%.0f", sum))
		err = writer.Write(towrite)
	}
	return nil
}

func (e *memoryOverheadExecutor) runMem(size int) (map[string]float64, error) {
	instance, err := e.runner.Run(size)
	defer e.runner.Stop(instance)
	if err != nil {
		return nil, err
	}

	// sleep a little bit to wait for startup
	fmt.Println("Sleeping to wait for the process to come up")
	<-time.After(2 * time.Second)
	data := make(map[string]float64)
	for s, p := range instance.Processes() {
		data[s] = float64(p.Pid)
		_, err := pmap.Run(p.Pid)
		if err != nil {
			return nil, err
		}

		// TODO parse pmap output
	}
	return data, nil
}
