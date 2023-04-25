package perftest

import (
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
	// wait for completion
	<-outchan
	return nil
}
