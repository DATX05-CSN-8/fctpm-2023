package handlers

import (
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/gin-gonic/gin"
)

type webhandlers struct {
	starterService vmStarterService
	dataService    vmDataRetrieverService
}

type vmStarterService interface {
	StartVM(config string) (string, error)
}

type vmDataRetrieverService interface {
	GetLogs(id string) (string, error)
	GetInfo(id string) (vminfo.VMInfo, error)
}

func AttachWebHandlers(rh gin.IRoutes, starterService vmStarterService, dataService vmDataRetrieverService) {
	whs := webhandlers{
		starterService: starterService,
		dataService:    dataService,
	}
	rh.POST("/vm", whs.startVMEndpoint)
	rh.GET("/vm/:id/logs", whs.getVMLogs)
	rh.GET("/vm/:id", whs.getVMInfo)
}
