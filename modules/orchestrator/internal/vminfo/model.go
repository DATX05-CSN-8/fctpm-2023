package vminfo

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Status int64

const (
	Running Status = iota
	Stopped
	Error
)

type VMInfo struct {
	Id        string
	StartTime time.Time
	EndTime   time.Time
	Status    Status
}

func NewVMInfo(startTime time.Time) VMInfo {
	return VMInfo{
		Id:        uuid.NewString(),
		StartTime: startTime,
		EndTime:   time.UnixMilli(0),
		Status:    Running,
	}
}

func (m *VMInfo) String() string {
	return fmt.Sprintf("Id: %s, Status: %d, Start: %s, End: %s", m.Id, m.Status, m.StartTime.String(), m.EndTime.String())
}
