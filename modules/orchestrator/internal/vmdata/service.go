package vmdata

import (
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

type ReadDeleteRepo[V vminfo.VMInfo | vmexecution.VMExecution] interface {
	FindAll() ([]V, error)
	FindById(id string) (V, error)
	Delete(data V) error
}

type VMDataRetriever struct {
	infos      ReadDeleteRepo[vminfo.VMInfo]
	executions ReadDeleteRepo[vmexecution.VMExecution]
}

func NewVMDataRetriever(infos ReadDeleteRepo[vminfo.VMInfo], executions ReadDeleteRepo[vmexecution.VMExecution]) *VMDataRetriever {
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

func (vdr *VMDataRetriever) GetAllInfo() ([]vminfo.VMInfo, error) {
	return vdr.infos.FindAll()
}

func (vdr *VMDataRetriever) Delete(id string) error {
	found, err := vdr.GetInfo(id)
	if err != nil {
		return err
	}
	return vdr.infos.Delete(found)
}
