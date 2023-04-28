package perftest

import (
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
)

type tpmRunner struct {
	config        *testRunnerConfig
	starter       VmStarter
	inforetriever *vmdata.VMDataRetriever
	numtests      int
	tpmalloc      tpmallocator
}

func NewTpmRunner(
	config *testRunnerConfig, starter VmStarter,
	inforetriever *vmdata.VMDataRetriever, numtests int, tpmalloc tpmallocator,
) *tpmRunner {
	return &tpmRunner{
		config:        config,
		starter:       starter,
		inforetriever: inforetriever,
		numtests:      numtests,
		tpmalloc:      tpmalloc,
	}
}

func (r *tpmRunner) RunInstance() error {
	path, err := dirutil.CreateTempDir(r.config.tempPath)
	if err != nil {
		return err
	}
	defer dirutil.RemoveTempDir(path)
	configPath := dirutil.JoinPath(path, r.config.templateName+".json")
	// copy values of struct
	templateData := *r.config.templateData

	time := time.Now()
	tpmInstance, err := r.tpmalloc.Allocate()
	if err != nil {
		return err
	}
	defer r.tpmalloc.Return(tpmInstance)
	templateData.TpmSocket = tpmInstance.SocketPath
	// generate config file
	err = firecracker.NewFirecrackerConfig(r.config.templateName, templateData, path)
	if err != nil {
		return err
	}
	return startvmBlocking(r.starter, configPath, time)
}

func (r *tpmRunner) Finish() error {
	return finish(r.inforetriever, &r.config.resultPath, r.numtests)
}
