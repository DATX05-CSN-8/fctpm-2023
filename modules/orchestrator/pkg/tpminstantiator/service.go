package tpminstantiator

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/google/uuid"
)

type tpmInstantiatorService struct {
	basePath *string
}

type TpmInstance struct {
	SocketPath string
	DirPath    string
	proc       *os.Process
}

func NewTpmInstantiatorService() *tpmInstantiatorService {
	return &tpmInstantiatorService{}
}

func NewTpmInstantiatorServiceWithBasePath(basepath string) *tpmInstantiatorService {
	return &tpmInstantiatorService{
		basePath: &basepath,
	}
}

func (s *tpmInstantiatorService) setupState(path string) error {
	cmd := exec.Command("swtpm_setup", "--tpm-state",
		path, "--createek", "--tpm2",
		"--create-ek-cert", "--create-platform-cert", "--lock-nvram", "--logfile",
		dirutil.JoinPath(path, "swtpm_setup.log"),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("Error starting swtpm_setup %s, output was %s", err, out)
	}
	return nil
}

func (s *tpmInstantiatorService) Create(path string) (*TpmInstance, error) {
	// setup swtpm state
	err := s.setupState(path)
	if err != nil {
		return nil, err
	}
	socketPath := dirutil.JoinPath(path, "socket")
	cmd := exec.Command(
		"swtpm", "socket",
		"--tpmstate", "dir="+path,
		"--tpm2", "--ctrl", "type=unixio,path="+socketPath,
		"--flags", "startup-clear",
		"--log", "file="+dirutil.JoinPath(path, "swtpm.log"),
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

func (s *tpmInstantiatorService) Allocate() (*TpmInstance, error) {
	id := uuid.NewString()
	path := dirutil.JoinPath(*s.basePath, id)
	err := dirutil.EnsureDirectory(path)
	if err != nil {
		return nil, err
	}
	instance, err := s.Create(path)
	if err != nil {
		return nil, err
	}
	// TODO use fsnotify here to wait until socket file exists
	time.Sleep(time.Millisecond * 100)
	return instance, nil
}

func (s *tpmInstantiatorService) Return(instance *TpmInstance) error {
	err := s.Destroy(instance)
	if err != nil {
		return err
	}
	return dirutil.RemoveDirIfExists(instance.DirPath)
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
