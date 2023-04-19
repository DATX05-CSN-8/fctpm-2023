package vmstarter

import (
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/google/uuid"
)

type infoRepo interface {
	Create(data vminfo.VMInfo) (vminfo.VMInfo, error)
	Update(data vminfo.VMInfo) (vminfo.VMInfo, error)
}

type executionRepo interface {
	Create(data vmexecution.VMExecution) (vmexecution.VMExecution, error)
	Delete(e vmexecution.VMExecution) error
}

type vmStarterService struct {
	fcc           firecracker.FirecrackerClient
	infoRepo      infoRepo
	executionRepo executionRepo
}

func NewVMStarterService(fcc firecracker.FirecrackerClient, infoRepo infoRepo, executionRepo executionRepo) *vmStarterService {
	return &vmStarterService{
		fcc:           fcc,
		infoRepo:      infoRepo,
		executionRepo: executionRepo,
	}
}

func (s *vmStarterService) StartVM(config string) (string, error) {
	return s.StartVMWithStartTime(config, time.Now())
}

func (s *vmStarterService) StartVMWithStartTime(config string, started time.Time) (string, error) {
	// generate ID
	id := uuid.NewString()
	// start execution
	exec, err := s.fcc.Start(config)
	if err != nil {
		return "", err
	}
	vmexec := vmexecution.NewVMExecution(id, exec)
	vi := vminfo.NewVMInfo(id, started)
	// persist info
	vi, err = s.infoRepo.Create(vi)
	if err != nil {
		return "", err
	}
	// persist execution
	vmexec, err = s.executionRepo.Create(vmexec)
	if err != nil {
		return "", err
	}
	// add subscriber to update info
	(*exec).Subscribe(func(status vminfo.Status) {
		vi.EndTime = time.Now()
		// get logs
		logs := (*exec).Logs()
		// parse logs
		bootTime, err := parseLogsForBootTime(logs)
		if err != nil {
			panic(err)
		}
		// persist in inforepo
		vi.ExecTime = bootTime.BootTime
		vi.Status = status
		s.infoRepo.Update(vi)
		// delete from execution repo
		s.executionRepo.Delete(vmexec)
	})
	return id, nil
}
