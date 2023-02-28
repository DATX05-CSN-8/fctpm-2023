package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GetVMInfoRequest struct {
	Id string `uri:"id" binding:"required,uuid"`
}

func (s *webhandlers) getVMInfo(c *gin.Context) {
	var r GetVMInfoRequest
	if err := c.ShouldBindUri(&r); err != nil {
		c.Status(400)
		fmt.Println(err)
		return
	}

	info, err := s.dataService.GetInfo(r.Id)
	if err == gorm.ErrRecordNotFound {
		c.Status(404)
		return
	}
	if err != nil {
		c.Status(500)
		fmt.Println(err)
		return
	}
	c.JSON(200, info)
}

func (s *webhandlers) getVMInfos(c *gin.Context) {
	infos, err := s.dataService.GetAllInfo()
	if err != nil {
		c.Status(500)
		fmt.Println(err)
		return
	}
	c.JSON(200, infos)
}
