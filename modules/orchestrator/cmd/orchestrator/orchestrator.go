package main

import (
	"flag"
	"os"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmstarter"
	handlers "github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/web"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func genDefaultFCBinPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return wd + "/../firecracker/bin/firecracker"
}

func genDefaultFCConfigPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return wd + "/../vm-start/fc-config.json"
}
func main() {
	apiSock := flag.String("api-sock", "/tmp/fctpm-orchestrator", "File path to the socket that should be listened to")
	dbPath := flag.String("db-path", "test.db", "File path to the db file to be used for sqlite")
	fcPath := flag.String("firecracker-bin", genDefaultFCBinPath(), "File path to the firecracker binary that should be used")
	flag.Parse()

	err := removeFileIfExists(*apiSock)
	if err != nil {
		panic("Could not remove socket file")
	}

	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to db")
	}
	db.AutoMigrate(&vminfo.VMInfo{})

	fcClient := firecracker.NewFirecrackerClient(*fcPath)
	vminfoRepo := vminfo.NewRepository(db)
	vmExecRepo := vmexecution.NewRepository()
	vmstarterService := vmstarter.NewVMStarterService(*fcClient, vminfoRepo, vmExecRepo)

	r := gin.Default()
	v1 := r.Group("/v1")
	handlers.AttachWebHandlers(v1, vmstarterService)

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
