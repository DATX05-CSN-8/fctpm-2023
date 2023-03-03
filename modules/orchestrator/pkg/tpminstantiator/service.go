package tpminstantiator

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/google/uuid"
)

type tpmInstantiatorService struct {
	swtpmPath string
	basePath  string
}

type TpmInstance struct {
	Id         string
	SocketPath string
	proc       *os.Process
	sb         *strings.Builder
}

func NewTpmInstantiatorService(swtpmPath string, basePath string) *tpmInstantiatorService {
	return &tpmInstantiatorService{
		swtpmPath: swtpmPath,
		basePath:  basePath,
	}
}

func joinPath(paths []string) string {
	return strings.Join(paths, string(os.PathSeparator))
}
func ensureDirectory(paths ...string) (string, error) {
	path := joinPath(paths)
	println("Path", path)
	err := os.MkdirAll(path, os.ModePerm)
	return path, err
}

func (s *tpmInstantiatorService) Create() (*TpmInstance, error) {
	// TODO instantiate TPM
	id := uuid.NewString()
	path, err := ensureDirectory(s.basePath, id)
	if err != nil {
		return nil, err
	}
	socketPath := joinPath([]string{path, "socket"})
	cmd := exec.Command(s.swtpmPath, "socket", "--tpmstate", "dir="+path, "--tpm2", "--ctrl", "type=unixio,path="+socketPath, "--log", "level=20")

	out, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	cmd.Stderr = cmd.Stdout

	err = cmd.Start()
	if err != nil {
		return nil, err
	}
	var sb strings.Builder
	go func() {
		b := make([]byte, 32)
		for {
			c, e := out.Read(b)
			if c < 0 {
				log.Fatal("Negative read swtpm", id)
			}
			sb.Write(b[:c])
			fmt.Printf("swtpm-%s: %s\n", id, b[:c])
			if e == io.EOF {
				return
			}
			if e != nil {
				log.Fatal(e)
			}
		}
	}()

	go func() {
		cmd.Wait()
	}()

	return &TpmInstance{
		Id:         id,
		SocketPath: socketPath,
		proc:       cmd.Process,
		sb:         &sb,
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
	path := joinPath([]string{s.basePath, instance.Id})
	err = os.RemoveAll(path)
	if err != nil {
		fmt.Println("Error occurred when removing swtpm directory", err)
		return err
	}
	return nil
}
