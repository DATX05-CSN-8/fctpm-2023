package perftest

import (
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
)

type tpmPoolRunner struct {
	config        *testRunnerConfig
	starter       VmStarter
	inforetriever *vmdata.VMDataRetriever
	numtests      int
	tpmalloc      tpmallocator
	numinstance   int
}

func NewTpmPoolRunner(
	config *testRunnerConfig, starter VmStarter,
	inforetriever *vmdata.VMDataRetriever, numtests int, tpmalloc tpmallocator, numinstances int,
) *tpmPoolRunner {
	return &tpmPoolRunner{
		config:        config,
		starter:       starter,
		inforetriever: inforetriever,
		numtests:      numtests,
		tpmalloc:      tpmalloc,
		numinstance:   numinstances,
	}
}

func (r *tpmPoolRunner) RunInstance() error {
	for i := 0; i < r.numinstance; i++ {
		// AAA todo concurrent
		// AAA todo semaphore? atleast to write to runner

		path, err := dirutil.CreateTempDir(r.config.tempPath)
		if err != nil {
			return err
		}
		configPath := dirutil.JoinPath(path, r.config.templateName+".json")
		// copy values of struct
		templateData := *r.config.templateData

		starttime := time.Now()

		tpmInstance, err := r.tpmalloc.Allocate()
		if err != nil {
			return err
		}
		templateData.TpmSocket = tpmInstance.SocketPath
		// generate config file
		err = firecracker.NewFirecrackerConfig(r.config.templateName, templateData, path)
		if err != nil {
			return err
		}
		err = startvmBlocking(r.starter, configPath, starttime)
		if err != nil {
			return err
		}
		err = r.tpmalloc.Return(tpmInstance)
		if err != nil {
			return err
		}
		err = dirutil.RemoveTempDir(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *tpmPoolRunner) Finish() error {
	return finish(r.inforetriever, &r.config.resultPath, r.numtests)
}
