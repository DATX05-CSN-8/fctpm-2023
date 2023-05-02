package pmap

import (
	"fmt"
	"os/exec"
)

func Run(pid int) (*string, error) {
	pidStr := fmt.Sprint(pid)

	cmd := exec.Command("pmap", "-x", pidStr)
	outp, err := cmd.CombinedOutput()
	str := string(outp[:])
	if err != nil {
		return nil, fmt.Errorf("Error occurred running pmap, output: %s, err: %v", str, err)
	}

	return &str, nil
}
