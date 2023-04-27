package vmstarter

import (
	"fmt"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/google/uuid"
)

type infoRepo interface {
	Create(data vminfo.VMInfo) (vminfo.VMInfo, error)
	Update(data *vminfo.VMInfo) (*vminfo.VMInfo, error)
}

type executionRepo interface {
	Create(data vmexecution.VMExecution) (vmexecution.VMExecution, error)
	Delete(e vmexecution.VMExecution) error
}

type vmStarterService struct {
	fcc           firecracker.FirecrackerClient
	infoRepo      infoRepo
	executionRepo executionRepo
	logWriterFn   func(*string)
}

func NewVMStarterService(fcc firecracker.FirecrackerClient, infoRepo infoRepo, executionRepo executionRepo) *vmStarterService {
	return &vmStarterService{
		fcc:           fcc,
		infoRepo:      infoRepo,
		executionRepo: executionRepo,
	}
}

func NewVMStarterServiceWithLogs(
	fcc firecracker.FirecrackerClient,
	infoRepo infoRepo,
	executionRepo executionRepo,
	logWriterFn func(*string),
) *vmStarterService {
	return &vmStarterService{
		fcc:           fcc,
		infoRepo:      infoRepo,
		executionRepo: executionRepo,
		logWriterFn:   logWriterFn,
	}
}

func (s *vmStarterService) StartVM(config string) (string, error) {
	exec, err := s.StartVMWithStartTime(config, time.Now())
	if err != nil {
		return "", err
	}
	return exec.Id, nil
}

func (s *vmStarterService) StartVMWithStartTime(config string, started time.Time) (*vmexecution.VMExecution, error) {
	// generate ID
	id := uuid.NewString()
	// start execution
	exec, err := s.fcc.Start(config)
	if err != nil {
		return nil, err
	}
	vmexec := vmexecution.NewVMExecution(id, exec)
	vi := vminfo.NewVMInfo(id, started)
	// persist info
	vi, err = s.infoRepo.Create(vi)
	if err != nil {
		return nil, err
	}
	// persist execution
	vmexec, err = s.executionRepo.Create(vmexec)
	if err != nil {
		return nil, err
	}
	// add subscriber to update info
	(*exec).Subscribe(func(status vminfo.Status) {
		vi.EndTime = time.Now()
		logs := (*exec).Logs()
		vi.Status = status
		defer s.executionRepo.Delete(vmexec)
		defer s.infoRepo.Update(&vi)
		defer func() {
			if s.logWriterFn != nil {
				s.logWriterFn(&logs)
			}
		}()

		if status == vminfo.Error {
			fmt.Printf("VM Stopped with error\n%s\n", logs)
			return
		}
		// parse logs
		bootTime, err := parseLogsForBootTime(logs)
		if err != nil {
			fmt.Printf("Error occurred parsing logs for boot time. Logs were:\n%s\n", logs)
			panic(err)
		}
		vi.ExecTime = bootTime.BootTime
	})
	return &vmexec, nil
}
