package firecracker

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type FirecrackerClient struct {
	binaryPath string
	timeout    time.Duration
}

func NewFirecrackerClient(binaryPath string) *FirecrackerClient {
	return &FirecrackerClient{
		binaryPath: binaryPath,
		timeout:    4 * time.Second,
	}
}

func (c *FirecrackerClient) Start(configPath string) (*FirecrackerExecution, error) {
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
				fmt.Println("Error reading logs")
				log.Fatal(e)
			}
		}
	}()

	outpc := make(chan error, 1)
	go func() {

		waitc := make(chan error, 1)
		go func() {
			waitc <- fcCmd.Wait()
		}()
		select {
		case ret := <-waitc:
			outpc <- ret
		case <-time.After(c.timeout):
			outpc <- fmt.Errorf("Timeout while running Firecracker")
		}
	}()

	return newFirecrackerExecution(&sb, outpc, fcCmd.Process), nil
}
