package vmstarter

import (
	"os"
	"testing"
	"time"
)

func TestParseLogs(t *testing.T) {
	// given
	logBytes, err := os.ReadFile("testdata/log1")
	if err != nil {
		t.Error("Could not open fixture", err)
	}
	logs := string(logBytes)
	expectedBoot, err := time.ParseDuration("886437us")
	if err != nil {
		t.Error("Error in creating expected boot value", err)
	}
	expectedCPU, err := time.ParseDuration("133461us")
	if err != nil {
		t.Error("Error in creating expected CPU value", err)
	}
	// when
	ds, err := parseLogsForBootTime(logs)
	// then
	if err != nil {
		t.Error("Error occurred parsing boot time", err)
		return
	}
	if ds.BootTime != expectedBoot {
		t.Error("Actual boot time was not same as expected", ds.BootTime, expectedBoot)
	}
	if ds.CpuTime != expectedCPU {
		t.Error("Actual CPU time was not same as expected", ds.CpuTime, expectedCPU)
	}
}
