package main

import (
	"flag"
	"os"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	apiSock := flag.String("api-sock", "/tmp/fctpm-orchestrator", "File path to the socket that should be listened to")
	flag.Parse()

	err := removeFileIfExists(*apiSock)
	if err != nil {
		panic("Could not remove socket file")
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to db")
	}
	db.AutoMigrate(&vminfo.VMInfo{})

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "world",
		})
	})

	err = r.RunUnix(*apiSock)
	if err != nil {
		panic(err)
	}
}

func removeFileIfExists(filename string) error {
	_, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return nil
	}
	return os.Remove(filename)
}
