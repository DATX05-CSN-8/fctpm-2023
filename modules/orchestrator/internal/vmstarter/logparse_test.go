package vmstarter

import (
	"os"
	"testing"
	"time"
)

func testWithData(logfile string, expectedBoot time.Duration, expectedCPU time.Duration, t *testing.T) {
	// given
	logBytes, err := os.ReadFile(logfile)
	if err != nil {
		t.Error("Could not open fixture", err)
	}
	logs := string(logBytes)
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

func TestParseLogs1(t *testing.T) {
	// given
	logfile := "testdata/log1"
	expectedBoot, err := time.ParseDuration("886437us")
	if err != nil {
		t.Error("Error in creating expected boot value", err)
	}
	expectedCPU, err := time.ParseDuration("133461us")
	if err != nil {
		t.Error("Error in creating expected CPU value", err)
	}
	testWithData(logfile, expectedBoot, expectedCPU, t)
}

func TestParseLogs2(t *testing.T) {
	// given
	logfile := "testdata/log2"
	expectedBoot, err := time.ParseDuration("314218us")
	if err != nil {
		t.Error("Error in creating expected boot value", err)
	}
	expectedCPU, err := time.ParseDuration("76741us")
	if err != nil {
		t.Error("Error in creating expected CPU value", err)
	}
	testWithData(logfile, expectedBoot, expectedCPU, t)
}
