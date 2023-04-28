package perftest

import (
	"fmt"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
)

type tpmPoolRunner struct {
	config        *testRunnerConfig
	starter       VmStarter
	inforetriever *vmdata.VMDataRetriever
	tpmpoolalloc  *TpmPool
	instancenum   int
}

func NewTpmPoolRunner(
	config *testRunnerConfig, starter VmStarter,
	inforetriever *vmdata.VMDataRetriever, tpmpoolalloc *TpmPool, num int,
) *tpmPoolRunner {
	return &tpmPoolRunner{
		config:        config,
		starter:       starter,
		inforetriever: inforetriever,
		tpmpoolalloc:  tpmpoolalloc,
		instancenum:   num,
	}
}

func (r *tpmPoolRunner) RunInstance() error {
	for i := 0; i < r.instancenum; i++ {
		// AAA todo concurrent
		// AAA todo semaphore? atleast to write to runner

		path, err := dirutil.CreateTempDir(r.config.tempPath)
		if err != nil {
			return err
		}
		dirutil.JoinPath(path, fmt.Sprint(i))
		configPath := dirutil.JoinPath(path, r.config.templateName+".json")
		// copy values of struct
		templateData := *r.config.templateData

		starttime := time.Now()
		tpmInstance := r.tpmpoolalloc.tpmq[i]

		templateData.TpmSocket = tpmInstance.instance.SocketPath
		// generate config file
		err = firecracker.NewFirecrackerConfig(r.config.templateName, templateData, path)
		if err != nil {
			return err
		}
		err = startvmBlocking(r.starter, configPath, starttime)
		if err != nil {
			return err
		}
		err = r.tpmpoolalloc.tpmq[i].alloc.Return(r.tpmpoolalloc.tpmq[i].instance)

		err = dirutil.RemoveTempDir(path)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *tpmPoolRunner) Finish() error {
	return finish(r.inforetriever, &r.config.resultPath)
}
