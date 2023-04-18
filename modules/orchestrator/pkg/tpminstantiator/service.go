package tpminstantiator

import (
	"fmt"
	"os"
	"strings"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/vtpm"
	"github.com/google/uuid"
)

type tpmInstantiatorService struct {
	swtpmPath     string
	swtpmBiosPath string
	basePath      string
}

type TpmInstance struct {
	Id         string
	DevicePath string
	device     *vtpm.VTPM
}

func NewTpmInstantiatorService(swtpmPath string, swtpmBiosPath string, basePath string) *tpmInstantiatorService {
	return &tpmInstantiatorService{
		swtpmPath:     swtpmPath,
		basePath:      basePath,
		swtpmBiosPath: swtpmBiosPath,
	}
}

func joinPath(paths []string) string {
	return strings.Join(paths, string(os.PathSeparator))
}
func ensureDirectory(paths ...string) (string, error) {
	path := joinPath(paths)
	err := os.MkdirAll(path, os.ModePerm)
	return path, err
}

func (s *tpmInstantiatorService) Create() (*TpmInstance, error) {
	id := uuid.NewString()
	// remove directory first
	err := os.RemoveAll(joinPath([]string{s.basePath, id}))
	if err != nil {
		return nil, err
	}
	statePath, err := ensureDirectory(s.basePath, id)
	if err != nil {
		return nil, err
	}
	vtpm, err := vtpm.NewVTPM(statePath, false, "", false, "melker", "", []byte{})
	if err != nil {
		return nil, err
	}

	vtpm.CreatedStatepath, err = vtpm.Start()
	if err != nil {
		return nil, err
	}

	instance := &TpmInstance{
		device:     vtpm,
		Id:         id,
		DevicePath: vtpm.GetTPMDevname(),
	}
	return instance, nil

	// fmt.Printf("Hostdev %s, major %d, minor %d, devpath %s\n", hostdev, major, minor, devpath)
	// r := bufio.NewReader(os.Stdin)
	// fmt.Print("Press enter to stop swtpm...")
	// _, _ = r.ReadString('\n')

	// return nil, fmt.Errorf("Mock error")
}

func (s *tpmInstantiatorService) Destroy(instance *TpmInstance) error {
	err := instance.device.Stop(false)
	if err != nil {
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

// func (s *tpmInstantiatorService) _Create() (*TpmInstance, error) {
// 	// TODO instantiate TPM
// 	id := uuid.NewString()
// 	path, err := ensureDirectory(s.basePath, id)
// 	if err != nil {
// 		return nil, err
// 	}
// 	socketPath := joinPath([]string{path, "socket"})
// 	cmd := exec.Command(
// 		s.swtpmPath, "socket",
// 		"--tpmstate", "dir="+path,
// 		"--tpm2", "--ctrl", "type=unixio,path="+socketPath,
// 		"--log", "level=20",
// 		"--flags", "not-need-init,startup-clear",
// 		"--locality", "reject-locality-4,allow-set-locality",
// 	)

// 	out, err := cmd.StdoutPipe()
// 	if err != nil {
// 		return nil, err
// 	}
// 	cmd.Stderr = cmd.Stdout

// 	err = cmd.Start()
// 	if err != nil {
// 		return nil, err
// 	}
// 	var sb strings.Builder
// 	go func() {
// 		b := make([]byte, 32)
// 		for {
// 			c, e := out.Read(b)
// 			if c < 0 {
// 				log.Fatal("Negative read swtpm", id)
// 			}
// 			sb.Write(b[:c])
// 			fmt.Printf("%s\n", b[:c])
// 			// fmt.Printf("swtpm-%s: %s\n", id, b[:c])
// 			if e == io.EOF {
// 				return
// 			}
// 			if e != nil {
// 				log.Fatal(e)
// 			}
// 		}
// 	}()

// 	go func() {
// 		cmd.Wait()
// 	}()

// 	instance := &TpmInstance{
// 		Id:         id,
// 		SocketPath: socketPath,
// 		proc:       cmd.Process,
// 		sb:         &sb,
// 	}
// 	return instance, nil
// 	// err = s.startupTpm(instance)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// return instance, nil
// }

// func (s *tpmInstantiatorService) Destroy(instance *TpmInstance) error {
// 	err := instance.proc.Kill()
// 	if err == os.ErrProcessDone {
// 		// Do nothing
// 	} else if err != nil {
// 		fmt.Println("Error occured when stopping swtpm process.", err)
// 		return err
// 	}
// 	path := joinPath([]string{s.basePath, instance.Id})
// 	err = os.RemoveAll(path)
// 	if err != nil {
// 		fmt.Println("Error occurred when removing swtpm directory", err)
// 		return err
// 	}
// 	return nil
// }
