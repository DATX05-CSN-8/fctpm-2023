package perftest

import (
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

type baselineRunner struct {
	config        *testRunnerConfig
	starter       vmStarter
	inforetriever *vmdata.VMDataRetriever
}

func NewBaselineRunner(config *testRunnerConfig, starter vmStarter, infoRetriever *vmdata.VMDataRetriever) *baselineRunner {
	return &baselineRunner{
		config:        config,
		starter:       starter,
		inforetriever: infoRetriever,
	}
}

func (br *baselineRunner) RunInstance() error {
	path, err := dirutil.CreateTempDir(br.config.tempPath)
	if err != nil {
		return err
	}
	defer dirutil.RemoveTempDir(path)
	configPath := dirutil.JoinPath(path, br.config.templateName+".json")
	outchan := make(chan int)

	time := time.Now()
	// generate config file
	err = firecracker.NewFirecrackerConfig(br.config.templateName, br.config.templateData, path)
	if err != nil {
		return err
	}
	exec, err := br.starter.StartVMWithStartTime(configPath, time)
	if err != nil {
		return err
	}
	exec.Exec.Subscribe(func(status vminfo.Status) {
		outchan <- 1
	})
	// wait for completion
	<-outchan
	return nil
}

func (br *baselineRunner) Finish() error {
	data, err := br.inforetriever.GetAllInfo()
	if err != nil {
		return err
	}
	return writeDataToCsv(data, br.config.resultPath)
}
