package main

import (
	"flag"

	"github.com/gin-gonic/gin"
)

func main() {
	apiSock := flag.String("api-sock", "/tmp/fctpm-orchestrator", "File path to the socket that should be listened to")
	flag.Parse()

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "world",
		})
	})

	err := r.RunUnix(*apiSock)
	if err != nil {
		panic(err)
	}
}
