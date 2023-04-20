package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/perftest"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmstarter"
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

func main() {

	dbPath := flag.String("db-path", "test.db", "File path to the db file to be used for sqlite")
	fcPath := flag.String("firecracker-bin", genDefaultFCBinPath(), "File path to the firecracker binary that should be used")
	resultPath := flag.String("result-path", "output.csv", "Path to CSV file to create")
	tempPath := flag.String("temp-path", "/tmp/firecracker-perftest", "Path to temporary data directory")
	flag.Parse()

	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to db")
	}
	db.AutoMigrate(&vminfo.VMInfo{})

	fcClient := firecracker.NewFirecrackerClient(*fcPath)
	vminfoRepo := vminfo.NewRepository(db)
	vmExecRepo := vmexecution.NewRepository()
	vmstarterService := vmstarter.NewVMStarterService(*fcClient, vminfoRepo, vmExecRepo)
	dataRetrieverService := vmdata.NewVMDataRetriever(vminfoRepo, vmExecRepo)

	perftestExecutor := perftest.NewPerftestExecutor(5, 1)
	baseTemplateData := firecracker.SimpleTemplateData{
		KernelImagePath: "/home/melker/fctpm-2023/vm-image/out/fc-image-kernel",
		InitRdPath:      "/home/melker/fctpm-2023/vm-image/out/fc-image-initrd.img",
	}
	runnerCfg := perftest.NewTestRunnerConfig(&baseTemplateData, "256-no-tpm", *tempPath, *resultPath)
	runner := perftest.NewBaselineRunner(runnerCfg, vmstarterService, dataRetrieverService)

	// execute
	fmt.Println("Running perf test")
	err = perftestExecutor.RunPerftest(runner)
	if err != nil {
		panic(err)
	}

}
