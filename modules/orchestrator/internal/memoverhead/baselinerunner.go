package memoverhead

import (
	"fmt"
	"os"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

type MemoverheadBaselineRunner struct {
	fcclient *firecracker.FirecrackerClient
	config   *testRunnerConfig
}

type baselineInstance struct {
	fcProcess *os.Process
	path      string
}

type memTemplateData struct {
	KernelImagePath string
	InitRdPath      string
	MemSize         int
}

func NewBaselineRunner(fc *firecracker.FirecrackerClient, config *testRunnerConfig) (MemoverheadBaselineRunner, error) {
	return MemoverheadBaselineRunner{
		fcclient: fc,
		config:   config,
	}, nil
}

func (r *MemoverheadBaselineRunner) Run(memsize int) (instance, error) {
	path, err := dirutil.CreateTempDir(r.config.tempPath)
	if err != nil {
		return nil, err
	}

	configPath := dirutil.JoinPath(path, r.config.templateName+".json")

	// copy values of struct
	templateData := *r.config.templateData
	templateData.MemSize = memsize

	err = firecracker.NewFirecrackerConfig(r.config.templateName, templateData, path)
	if err != nil {
		return nil, err
	}

	execution, err := r.fcclient.Start(configPath)
	if err != nil {
		return nil, err
	}
	execution.Subscribe(func(status vminfo.Status) {
		fmt.Printf("Logs received on exit\n%s\n", execution.Logs())
		dirutil.RemoveTempDir(path)
	})
	if execution.Status() == vminfo.Error {
		defer dirutil.RemoveTempDir(path)
		return nil, fmt.Errorf("Error occurred early when starting firecracker. Logs will be shown on the next lines\n%s", execution.Logs())
	}

	return &baselineInstance{
		fcProcess: execution.Process(),
		path:      path,
	}, nil
}

func (r *MemoverheadBaselineRunner) Stop(inst instance) error {
	processes := inst.Processes()
	err := processes["firecracker"].Kill()
	defer inst.Cleanup()
	return err
}

func (b *baselineInstance) Processes() map[string]*os.Process {
	m := map[string]*os.Process{
		"firecracker": b.fcProcess,
	}
	return m
}

func (b *baselineInstance) Cleanup() error {
	return dirutil.RemoveDirIfExists(b.path)
}
