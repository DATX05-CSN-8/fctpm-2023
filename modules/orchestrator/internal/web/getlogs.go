package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type GetVMLogsRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

func (s *webhandlers) getVMLogs(c *gin.Context) {
	var r GetVMLogsRequest
	if err := c.ShouldBindUri(&r); err != nil {
		c.Status(400)
		fmt.Println(err)
		return
	}

	logs, err := s.dataService.GetLogs(r.Id)
	if err == fmt.Errorf("NOT_FOUND") {
		c.Status(404)
		return
	}
	if err != nil {
		c.Status(500)
		fmt.Println(err)
	}
	c.JSON(200, gin.H{
		"log": logs,
	})
}
