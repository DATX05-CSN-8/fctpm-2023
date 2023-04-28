package memoverhead

import (
	"os"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
)

type MemoverheadBaselineRunner struct {
	fcclient *firecracker.FirecrackerClient
	config   *testRunnerConfig
}

type baselineInstance struct {
	fcProcess *os.Process
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
	defer dirutil.RemoveTempDir(path)
	configPath := dirutil.JoinPath(path, r.config.templateName+".json")

	// copy values of struct
	templateData := *r.config.templateData
	templateData.MemSize = memsize

	err = firecracker.NewFirecrackerConfig(r.config.templateName, templateData, configPath)
	if err != nil {
		return nil, err
	}

	execution, err := r.fcclient.Start(configPath)
	if err != nil {
		return nil, err
	}
	return &baselineInstance{
		fcProcess: execution.Process(),
	}, nil
}

func (r *MemoverheadBaselineRunner) Stop(inst instance) error {
	processes := inst.Processes()
	err := processes["firecracker"].Kill()
	return err
}

func (b *baselineInstance) Processes() map[string]*os.Process {
	m := map[string]*os.Process{
		"firecracker": b.fcProcess,
	}
	return m
}
