package vmdata

import (
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

type ReadRepo[V vminfo.VMInfo | vmexecution.VMExecution] interface {
	FindAll() ([]V, error)
	FindById(id string) (V, error)
}

type VMDataRetriever struct {
	infos      ReadRepo[vminfo.VMInfo]
	executions ReadRepo[vmexecution.VMExecution]
}

func NewVMDataRetriever(infos ReadRepo[vminfo.VMInfo], executions ReadRepo[vmexecution.VMExecution]) *VMDataRetriever {
	return &VMDataRetriever{
		infos:      infos,
		executions: executions,
	}
}

func (vdr *VMDataRetriever) GetLogs(id string) (string, error) {
	found, err := vdr.executions.FindById(id)
	if err != nil {
		return "", err
	}
	return found.Exec.Logs(), nil
}

func (vdr *VMDataRetriever) GetInfo(id string) (vminfo.VMInfo, error) {
	return vdr.infos.FindById(id)
}
