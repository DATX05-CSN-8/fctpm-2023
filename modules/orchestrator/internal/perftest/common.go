package perftest

import (
	"fmt"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

func finish(inforetriever *vmdata.VMDataRetriever, resultpath *string) error {
	data, err := inforetriever.GetAllInfo()
	if err != nil {
		return err
	}
	return writeDataToCsv(data, *resultpath)
}

func startvmBlocking(starter VmStarter, configpath string, time time.Time) error {
	outchan := make(chan int)
	exec, err := starter.StartVMWithStartTime(configpath, time)
	if err != nil {
		return err
	}
	exec.Exec.Subscribe(func(status vminfo.Status) {
		outchan <- 1
	})
	// Need not wait as we did not have time to subscribe.
	if exec.Exec.Status() == vminfo.Error {
		return fmt.Errorf(exec.Exec.Logs())
	} else if exec.Exec.Status() == vminfo.Stopped {
		return nil
	}
	// wait for completion
	<-outchan
	return nil
}
