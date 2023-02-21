package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type startVMRequest struct {
	Config string `binding:"required"`
}

func (r *startVMRequest) String() string {
	return fmt.Sprintf("Config: %s", r.Config)
}

func (s *webhandlers) startVMEndpoint(c *gin.Context) {
	var r startVMRequest
	err := c.BindJSON(&r)
	if err != nil {
		c.Status(400)
		fmt.Println(err)
		return
	}
	id, err := s.starterService.StartVM(r.Config)
	if err != nil {
		c.Status(500)
		fmt.Println(err)
		return
	}

	c.JSON(200, gin.H{
		"id": id,
	})
}
