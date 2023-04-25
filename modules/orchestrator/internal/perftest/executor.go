package perftest

import (
	"fmt"
)

type PerftestRunner interface {
	// Assumed to be blocking
	RunInstance() error
	Finish() error
}

type perftestExecutor struct {
	total      int
	concurrent int
}

func NewPerftestExecutor(total int, concurrent int) *perftestExecutor {
	return &perftestExecutor{
		total:      total,
		concurrent: concurrent,
	}
}

func (pte *perftestExecutor) RunPerftest(runner PerftestRunner) error {
	sem := make(chan int, pte.concurrent)
	errMutex := make(chan int, 1)
	errs := make([]error, 0)
	for i := 0; i < pte.total; i++ {
		sem <- 1
		go func() {
			err := runner.RunInstance()
			<-sem
			if err != nil {
				errMutex <- 1
				errs = append(errs, err)
				<-errMutex
			}
		}()
	}
	// wait for all semaphores to finish
	for i := 0; i < pte.concurrent; i++ {
		sem <- 1
	}

	err := runner.Finish()
	if err != nil {
		errs = append(errs, err)
	}
	if len(errs) != 0 {
		for _, v := range errs {
			fmt.Println(v)
		}
		return fmt.Errorf("%d errors occurred. ", len(errs))
	}
	return nil
}
