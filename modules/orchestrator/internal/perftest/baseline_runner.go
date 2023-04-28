package perftest

import (
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
)

type baselineRunner struct {
	config        *testRunnerConfig
	starter       VmStarter
	inforetriever *vmdata.VMDataRetriever
	numtests      int
}

func NewBaselineRunner(config *testRunnerConfig, starter VmStarter, infoRetriever *vmdata.VMDataRetriever, numtests int) *baselineRunner {
	return &baselineRunner{
		config:        config,
		starter:       starter,
		inforetriever: infoRetriever,
		numtests:      numtests,
	}
}

func (br *baselineRunner) RunInstance() error {
	path, err := dirutil.CreateTempDir(br.config.tempPath)
	if err != nil {
		return err
	}
	defer dirutil.RemoveTempDir(path)
	configPath := dirutil.JoinPath(path, br.config.templateName+".json")

	time := time.Now()
	// generate config file
	err = firecracker.NewFirecrackerConfig(br.config.templateName, br.config.templateData, path)
	if err != nil {
		return err
	}
	return startvmBlocking(br.starter, configPath, time)
}

func (br *baselineRunner) Finish() error {
	return finish(br.inforetriever, &br.config.resultPath, br.numtests)
}
