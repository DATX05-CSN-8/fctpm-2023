package firecracker

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
)

type FirecrackerClient interface {
	Start(config string) (FirecrackerExecution, error)
}

type firecrackerClient struct {
	binaryPath string
}

func NewFirecrackerClient(binaryPath string) *firecrackerClient {
	return &firecrackerClient{
		binaryPath: binaryPath,
	}
}

func (c *firecrackerClient) Start(configPath string) (*firecrackerExecution, error) {
	_, err := os.Stat(configPath)
	if os.IsNotExist(err) {
		return nil, err
	}
	fcCmd := exec.Command(c.binaryPath, "--no-api", "--config-file", configPath, "--boot-timer")
	out, err := fcCmd.StdoutPipe()
	fcCmd.Stderr = fcCmd.Stdout
	if err != nil {
		return nil, err
	}
	var sb strings.Builder

	err = fcCmd.Start()
	if err != nil {
		return nil, err
	}
	go func() {
		b := make([]byte, 8)
		for {
			c, e := out.Read(b)

			if c < 0 {
				panic("Negative read")
			}
			sb.Write(b[:c])

			if e == io.EOF {
				return
			}
			if e != nil {
				log.Fatal(e)
			}
		}
	}()

	outpc := make(chan error, 1)
	go func() {
		outpc <- fcCmd.Wait()
	}()
	return newFirecrackerExecution(&sb, outpc), nil
}
