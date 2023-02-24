package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DeleteVMInfoRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

func (s *webhandlers) deleteVMInfo(c *gin.Context) {
	var r DeleteVMInfoRequest
	if err := c.ShouldBindUri(&r); err != nil {
		c.Status(400)
		fmt.Println(err)
		return
	}
	err := s.dataService.Delete(r.Id)
	if err == gorm.ErrRecordNotFound {
		c.Status(404)
		return
	}
	if err != nil {
		c.Status(500)
		fmt.Println(err)
		return
	}

	c.Status(200)
}
