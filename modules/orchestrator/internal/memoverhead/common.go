package memoverhead

import "github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"

type testRunnerConfig struct {
	templateData *firecracker.SimpleTemplateData
	templateName string
	tempPath     string
}

type memTemplateData struct {
	KernelImagePath string
	InitRdPath      string
	MemSize         int
}

func NewMemoryOverheadConfig(templatedata *firecracker.SimpleTemplateData, templateName string, tempPath string) *testRunnerConfig {
	return &testRunnerConfig{
		templateData: templatedata,
		templateName: templateName,
		tempPath:     tempPath,
	}
}
