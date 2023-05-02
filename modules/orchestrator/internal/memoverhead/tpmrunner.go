package memoverhead

import (
	"fmt"
	"os"
	"strings"
	"syscall"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/tpminstantiator"
)

type MemoverheadTpmRunner struct {
	fcclient     *firecracker.FirecrackerClient
	tpmallocator tpmallocator
	config       *testRunnerConfig
}

type tpmallocator interface {
	Allocate() (*tpminstantiator.TpmInstance, error)
	Return(*tpminstantiator.TpmInstance) error
}

type tpmInstance struct {
	fcProcess   *os.Process
	tpmInstance *tpminstantiator.TpmInstance
	path        string
}

func NewTpmRunner(fc *firecracker.FirecrackerClient, config *testRunnerConfig, tpmalloc tpmallocator) (MemoverheadTpmRunner, error) {
	return MemoverheadTpmRunner{
		fcclient:     fc,
		tpmallocator: tpmalloc,
		config:       config,
	}, nil
}

func (r *MemoverheadTpmRunner) Run(memsize int) (instance, error) {
	path, err := dirutil.CreateTempDir(r.config.tempPath)
	if err != nil {
		return nil, err
	}

	// allocate tpm
	tpm, err := r.tpmallocator.Allocate()
	if err != nil {
		r.tpmallocator.Return(tpm)
		return nil, err
	}
	configPath := dirutil.JoinPath(path, r.config.templateName+".json")

	// copy values of struct
	templateData := *r.config.templateData
	templateData.MemSize = memsize
	templateData.TpmSocket = tpm.SocketPath

	err = firecracker.NewFirecrackerConfig(r.config.templateName, templateData, path)
	if err != nil {
		return nil, err
	}

	execution, err := r.fcclient.Start(configPath)
	if err != nil {
		return nil, err
	}
	execution.Subscribe(func(status vminfo.Status) {
		badExit := !strings.HasSuffix(execution.Logs(), "Error occurred: signal: interrupt")
		if status != vminfo.Stopped && badExit {
			fmt.Printf("Logs received on error exit\n%s\n", execution.Logs())
		}
		dirutil.RemoveTempDir(path)
	})
	if execution.Status() == vminfo.Error {
		defer dirutil.RemoveTempDir(path)
		return nil, fmt.Errorf("Error occurred early when starting firecracker. Logs will be shown on the next lines\n%s", execution.Logs())
	}

	return &tpmInstance{
		fcProcess:   execution.Process(),
		tpmInstance: tpm,
		path:        path,
	}, nil
}

func (r *MemoverheadTpmRunner) Stop(inst instance) error {
	processes := inst.Processes()

	err := processes["firecracker"].Signal(syscall.SIGINT)
	if err != nil {
		fmt.Println("Error occurred signaling firecracker")
	}
	tpmInst, ok := inst.(*tpmInstance)
	if !ok {
		return fmt.Errorf("Instance received was of wrong type, tpm runner")
	}

	err = r.tpmallocator.Return(tpmInst.tpmInstance)
	if err != nil {
		fmt.Println("Error occured returning tpm", err)
	}

	defer inst.Cleanup()
	return err
}

func (i *tpmInstance) Processes() map[string]*os.Process {
	m := map[string]*os.Process{
		"firecracker": i.fcProcess,
		"swtpm":       i.tpmInstance.Process(),
	}
	return m
}

func (i *tpmInstance) Cleanup() error {
	return dirutil.RemoveDirIfExists(i.path)
}
