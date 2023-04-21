package vmstarter

import (
	"fmt"
	"regexp"
	"time"
)

type BootTime struct {
	BootTime time.Duration
	CpuTime  time.Duration
}

var exp = regexp.MustCompile(`(?im)Guest-boot-time = (\d+) us \d+ ms, +(\d+) CPU us \d+ CPU ms$`)

func parseLogsForBootTime(logs string) (*BootTime, error) {
	found := exp.FindStringSubmatch(logs)
	if found == nil {
		return nil, fmt.Errorf("None returned by regexp")
	}
	bttime, err := time.ParseDuration(found[1] + "us")
	if err != nil {
		return nil, fmt.Errorf("Could not parse boot time number %s", string(found[1]))
	}
	cputime, err := time.ParseDuration(found[2] + "us")
	if err != nil {
		return nil, fmt.Errorf("Could not parse CPU time number %s", string(found[2]))
	}
	bt := &BootTime{
		BootTime: bttime,
		CpuTime:  cputime,
	}
	return bt, nil
}
