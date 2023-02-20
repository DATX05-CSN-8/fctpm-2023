package firecracker

import (
	"strings"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

type FirecrackerExecution interface {
	Status() vminfo.Status
	Logs() string
}

type firecrackerExecution struct {
	sb      *strings.Builder
	statusp *vminfo.Status
}

func newFirecrackerExecution(sb *strings.Builder, outpc chan error) *firecrackerExecution {
	var status vminfo.Status = vminfo.Running
	go func() {
		err := <-outpc
		if err != nil {
			status = vminfo.Error
		} else {
			status = vminfo.Stopped
		}
	}()
	return &firecrackerExecution{
		sb:      sb,
		statusp: &status,
	}
}

func (f *firecrackerExecution) Status() vminfo.Status {
	return *f.statusp
}

func (f *firecrackerExecution) Logs() string {
	return f.sb.String()
}
