package vminfo

import (
	"fmt"
	"time"
)

type Status int64

const (
	Running Status = iota
	Stopped
	Error
)

type VMInfo struct {
	Id        string `gorm:"primary_key"`
	StartTime time.Time
	ExecTime  time.Duration
	EndTime   time.Time
	Status    Status
}

func NewVMInfo(id string, startTime time.Time) VMInfo {
	zeroDuration, _ := time.ParseDuration("0ms")
	return VMInfo{
		Id:        id,
		StartTime: startTime,
		ExecTime:  zeroDuration,
		EndTime:   startTime,
		Status:    Running,
	}
}

func (m *VMInfo) String() string {
	return fmt.Sprintf("Id: %s, Status: %d, Start: %s, Exec: %s", m.Id, m.Status, m.StartTime.String(), m.ExecTime.String())
}
