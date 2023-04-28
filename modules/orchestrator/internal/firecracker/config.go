package firecracker

import (
	"os"
	"text/template"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
)

type SimpleTemplateData struct {
	KernelImagePath string
	InitRdPath      string
	TpmSocket       string
	MemSize         int
	BootArgs        string
}

func NewFirecrackerConfig(templatename string, data any, outpath string) error {
	// read template file with template name
	err := dirutil.RemoveFileIfExists(outpath)
	if err != nil {
		return err
	}
	err = dirutil.EnsureDirectory(outpath)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFiles(dirutil.JoinPath("resources", "firecracker-configs", templatename+".json.tmpl"))
	if err != nil {
		return err
	}
	outPath := dirutil.JoinPath(outpath, templatename+".json")
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, data)
	if err != nil {
		return err
	}
	return nil
}
