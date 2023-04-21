package perftest

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

type testPerftestRunner struct {
	sb         *strings.Builder
	writemutex chan int
	idx        int
	sleepdur   time.Duration
}

func (r *testPerftestRunner) RunInstance() error {
	r.writemutex <- 1
	curr_idx := r.idx
	r.idx = r.idx + 1
	r.sb.WriteString(fmt.Sprintf("%d starting\n", curr_idx))
	<-r.writemutex
	time.Sleep(r.sleepdur + time.Duration(curr_idx*2*1000*1000))
	r.writemutex <- 1
	r.sb.Write([]byte(fmt.Sprintf("%d finished\n", curr_idx)))
	<-r.writemutex
	return nil
}

func (r *testPerftestRunner) Finish() error {
	return nil
}

func (r *testPerftestRunner) GetOutput() string {
	return r.sb.String()
}

func NewPerfTestRunner() *testPerftestRunner {
	var sb strings.Builder
	dur, err := time.ParseDuration("20ms")
	if err != nil {
		panic(err)
	}
	return &testPerftestRunner{
		sb:         &sb,
		writemutex: make(chan int, 1),
		idx:        0,
		sleepdur:   dur,
	}
}

func TestRunPerfTestSequential(t *testing.T) {
	// given
	total := 5
	pte := NewPerftestExecutor(5, 1)
	runner := NewPerfTestRunner()
	expected := ""
	for i := 0; i < total; i++ {
		expected += fmt.Sprintf("%d starting\n", i)
		expected += fmt.Sprintf("%d finished\n", i)
	}
	// expectedOutput := ""
	// when
	err := pte.RunPerftest(runner)
	if err != nil {
		t.Error("Error running perf test", err)
	}
	// then
	output := runner.GetOutput()
	if output != expected {
		t.Error("Actual output was not same as expected", output, expected)
	}
}

func TestRunPerfTestConcurrent(t *testing.T) {
	// given
	total := 10
	concurrent := 3
	pte := NewPerftestExecutor(total, concurrent)
	runner := NewPerfTestRunner()
	expected := ""
	for i := 0; i < concurrent; i++ {
		expected += fmt.Sprintf("%d starting\n", i)
	}
	for i := concurrent; i < total; i++ {
		expected += fmt.Sprintf("%d finished\n", i-concurrent)
		expected += fmt.Sprintf("%d starting\n", i)
	}
	for i := total - concurrent; i < total; i++ {
		expected += fmt.Sprintf("%d finished\n", i)
	}
	// expectedOutput := ""
	// when
	err := pte.RunPerftest(runner)
	if err != nil {
		t.Error("Error running perf test", err)
	}
	// then
	output := runner.GetOutput()
	if output != expected {
		t.Error("Actual output was not same as expected", output, expected)
	}
}
