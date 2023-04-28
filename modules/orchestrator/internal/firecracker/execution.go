package firecracker

import (
	"fmt"
	"os"
	"strings"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

type StatusCallback func(vminfo.Status)

type FirecrackerExecution struct {
	sb          *strings.Builder
	statusp     *vminfo.Status
	subscribers *[]StatusCallback
	process     *os.Process
}

func newFirecrackerExecution(sb *strings.Builder, outpc chan error, process *os.Process) *FirecrackerExecution {
	var status vminfo.Status = vminfo.Running
	subscribers := make([]StatusCallback, 0)
	go func() {
		err := <-outpc
		if err != nil {
			sb.WriteString(fmt.Sprintf("Error occurred: %s", err.Error()))
			status = vminfo.Error
		} else {
			status = vminfo.Stopped
		}
		for _, s := range subscribers {
			s(status)
		}
	}()
	return &FirecrackerExecution{
		sb:          sb,
		statusp:     &status,
		subscribers: &subscribers,
		process:     process,
	}
}

func (f *FirecrackerExecution) Status() vminfo.Status {
	return *f.statusp
}

func (f *FirecrackerExecution) Logs() string {
	return f.sb.String()
}

func (f *FirecrackerExecution) Process() *os.Process {
	return f.process
}

func (f *FirecrackerExecution) Subscribe(cb func(vminfo.Status)) {
	*f.subscribers = append(*f.subscribers, cb)
}
