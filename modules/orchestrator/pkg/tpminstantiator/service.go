package tpminstantiator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

type tpmInstantiatorService struct {
}

type TpmInstance struct {
	SocketPath string
	DirPath    string
	proc       *os.Process
}

func NewTpmInstantiatorService() *tpmInstantiatorService {
	return &tpmInstantiatorService{}
}

func joinPath(paths ...string) string {
	return strings.Join(paths, string(os.PathSeparator))
}
func ensureDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	return err
}

func (s *tpmInstantiatorService) setupState(path string) error {
	cmd := exec.Command("swtpm_setup", "--tpm-state",
		path, "--createek", "--tpm2",
		"--create-ek-cert", "--create-platform-cert", "--lock-nvram", "--logfile",
		joinPath(path, "swtpm_setup.log"),
	)
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Error executing swtpm_setup %s", err)
	}
	err = cmd.Wait()
	if err != nil {
		log.Fatalf("Error executing swtpm_setup %s", err)
	}
	return nil
}

func (s *tpmInstantiatorService) Create(path string) (*TpmInstance, error) {
	// setup swtpm state
	err := s.setupState(path)
	if err != nil {
		return nil, err
	}
	socketPath := joinPath(path, "socket")
	cmd := exec.Command(
		"swtpm", "socket",
		"--tpmstate", "dir="+path,
		"--tpm2", "--ctrl", "type=unixio,path="+socketPath,
		"--flags", "startup-clear",
		"--log", "file="+joinPath(path, "swtpm.log"),
	)

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	go func() {
		cmd.Wait()
	}()

	return &TpmInstance{
		SocketPath: socketPath,
		proc:       cmd.Process,
		DirPath:    path,
	}, nil
}

func (s *tpmInstantiatorService) Destroy(instance *TpmInstance) error {
	err := instance.proc.Kill()
	if err == os.ErrProcessDone {
		// Do nothing
	} else if err != nil {
		fmt.Println("Error occured when stopping swtpm process.", err)
		return err
	}
	return nil
}
