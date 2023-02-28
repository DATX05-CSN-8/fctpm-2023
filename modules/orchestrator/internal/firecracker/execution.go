package firecracker

import (
	"strings"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

type StatusCallback func(vminfo.Status)

type FirecrackerExecution struct {
	sb          *strings.Builder
	statusp     *vminfo.Status
	subscribers *[]StatusCallback
}

func newFirecrackerExecution(sb *strings.Builder, outpc chan error) *FirecrackerExecution {
	var status vminfo.Status = vminfo.Running
	subscribers := make([]StatusCallback, 0)
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
	return &FirecrackerExecution{
		sb:          sb,
		statusp:     &status,
		subscribers: &subscribers,
	}
}

func (f *FirecrackerExecution) Status() vminfo.Status {
	return *f.statusp
}

func (f *FirecrackerExecution) Logs() string {
	return f.sb.String()
}

func (f *FirecrackerExecution) Subscribe(cb func(vminfo.Status)) {
	*f.subscribers = append(*f.subscribers, cb)
}
