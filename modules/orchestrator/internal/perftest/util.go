package perftest

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

type testRunnerConfig struct {
	templateData *firecracker.SimpleTemplateData
	templateName string
	tempPath     string
	resultPath   string
}

type vmStarter interface {
	StartVMWithStartTime(config string, started time.Time) (*vmexecution.VMExecution, error)
}

func NewTestRunnerConfig(templateData *firecracker.SimpleTemplateData, templateName string, tempPath string, resultPath string) *testRunnerConfig {
	return &testRunnerConfig{
		templateData: templateData,
		templateName: templateName,
		tempPath:     tempPath,
		resultPath:   resultPath,
	}
}

func writeDataToCsv(data []vminfo.VMInfo, outpath string) error {
	file, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	err = writer.Write([]string{"idx", "Start Time", "Exec Time", "End Time"})
	if err != nil {
		return err
	}
	for i, info := range data {
		row := []string{
			fmt.Sprint(i),
			fmt.Sprint(info.StartTime.UnixMicro()),
			fmt.Sprint(info.ExecTime.Microseconds()),
			fmt.Sprint(info.EndTime.UnixMicro()),
		}
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}
	return nil
}
