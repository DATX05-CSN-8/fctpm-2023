package vminfo

import "time"

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
