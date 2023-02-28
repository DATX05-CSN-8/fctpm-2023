package vmexecution

import (
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
)

type VMExecution struct {
	Id   string
	Exec *firecracker.FirecrackerExecution
}

func NewVMExecution(id string, e *firecracker.FirecrackerExecution) VMExecution {
	return VMExecution{
		Id:   id,
		Exec: e,
	}
}
