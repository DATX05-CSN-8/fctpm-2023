package firecracker

import (
	"os"
	"strings"
	"text/template"
)

func joinPath(paths []string) string {
	return strings.Join(paths, string(os.PathSeparator))
}

func removeFileIfExists(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}
	return os.RemoveAll(filename)
}

func ensureDirectory(paths ...string) (string, error) {
	path := joinPath(paths)
	err := os.MkdirAll(path, os.ModePerm)
	return path, err
}

type SimpleTemplateData struct {
	KernelImagePath string
	InitRdPath      string
	TpmSocket       string
}

func NewFirecrackerConfig(templatename string, data any, outpath string) error {
	// read template file with template name
	err := removeFileIfExists(outpath)
	if err != nil {
		return err
	}
	_, err = ensureDirectory(outpath)
	if err != nil {
		return err
	}
	tmpl, err := template.ParseFiles(joinPath([]string{"resources", "firecracker-configs", templatename + ".json.tmpl"}))
	if err != nil {
		return err
	}
	outPath := joinPath([]string{outpath, templatename + ".json"})
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
