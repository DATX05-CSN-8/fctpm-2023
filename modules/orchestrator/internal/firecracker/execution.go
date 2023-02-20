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
	sb          *strings.Builder
	statusp     *vminfo.Status
	subscribers *[]func(vminfo.Status)
}

func newFirecrackerExecution(sb *strings.Builder, outpc chan error) *firecrackerExecution {
	var status vminfo.Status = vminfo.Running
	subscribers := make([]func(vminfo.Status), 0)
	go func() {
		err := <-outpc
		if err != nil {
			status = vminfo.Error
		} else {
			status = vminfo.Stopped
		}
		for _, s := range subscribers {
			s(status)
		}
	}()
	return &firecrackerExecution{
		sb:          sb,
		statusp:     &status,
		subscribers: &subscribers,
	}
}

func (f *firecrackerExecution) Status() vminfo.Status {
	return *f.statusp
}

func (f *firecrackerExecution) Logs() string {
	return f.sb.String()
}

func (f *firecrackerExecution) Subscribe(cb func(vminfo.Status)) {
	*f.subscribers = append(*f.subscribers, cb)
}
