package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/perftest"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmstarter"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/tpminstantiator"
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

func getLogWriterFn(logPath *string) (func(*string), error) {
	err := dirutil.RemoveDirIfExists(*logPath)
	if err != nil {
		return nil, err
	}
	err = dirutil.EnsureDirectory(*logPath)
	if err != nil {
		return nil, err
	}
	idx := 0
	idxmutex := make(chan int, 1)
	fn := func(logs *string) {
		idxmutex <- 1
		i := idx
		idx += 1
		<-idxmutex
		filename := fmt.Sprint(i) + "-bootlog"
		err := os.WriteFile(dirutil.JoinPath(*logPath, filename), []byte(*logs), 0644)
		if err != nil {
			fmt.Printf("An error occurred writing logs: %s", err)
		}
	}
	return fn, nil
}

func main() {

	dbPath := flag.String("db-path", "test.db", "File path to the db file to be used for sqlite")
	fcPath := flag.String("firecracker-bin", genDefaultFCBinPath(), "File path to the firecracker binary that should be used")
	resultPath := flag.String("result-path", "output.csv", "Path to CSV file to create")
	tempPath := flag.String("temp-path", "/tmp/firecracker-perftest", "Path to temporary data directory")
	clean := flag.Bool("clean", false, "Clean the output database and csv")
	bootLogPath := flag.String("boot-log-path", "", "(optional) path to output boot logs")
	rtype := flag.String("type", "baseline", "Type of performance test to run. Either 'baseline' or 'tpm'")
	inum := flag.Int("num", 1, "Number of VM instances to run. Value between '1' and '1000'")
	flag.Parse()

	if *clean {
		err := dirutil.RemoveFileIfExists(*dbPath)
		if err != nil {
			fmt.Println("Error occurred removing db")
			panic(err)
		}
		err = dirutil.RemoveFileIfExists(*resultPath)
		if err != nil {
			fmt.Println("Error occurred removing result file")
			panic(err)
		}
	}

	db, err := gorm.Open(sqlite.Open(*dbPath), &gorm.Config{})
	if err != nil {
		panic("Failed to connect to db")
	}
	db.AutoMigrate(&vminfo.VMInfo{})

	fcClient := firecracker.NewFirecrackerClient(*fcPath)
	vminfoRepo := vminfo.NewRepository(db)
	vmExecRepo := vmexecution.NewRepository()
	var vmstarterService perftest.VmStarter
	if *bootLogPath != "" {
		logWriterFn, err := getLogWriterFn(bootLogPath)
		if err != nil {
			fmt.Println("Error occurred creating log writer fn")
			panic(err)
		}
		vmstarterService = vmstarter.NewVMStarterServiceWithLogs(*fcClient, vminfoRepo, vmExecRepo, logWriterFn)
	} else {
		vmstarterService = vmstarter.NewVMStarterService(*fcClient, vminfoRepo, vmExecRepo)
	}

	dataRetrieverService := vmdata.NewVMDataRetriever(vminfoRepo, vmExecRepo)

	perftestExecutor := perftest.NewPerftestExecutor(5, 1)
	baseTemplateData := firecracker.SimpleTemplateData{
		KernelImagePath: "/home/melker/fctpm-2023/vm-image/out/fc-image-kernel",
		InitRdPath:      "/home/melker/fctpm-2023/vm-image/out/fc-image-initrd.img",
	}

	var runner perftest.PerftestRunner
	if *rtype == "baseline" {
		templateName := "256-no-tpm"
		runnerCfg := perftest.NewTestRunnerConfig(&baseTemplateData, templateName, *tempPath, *resultPath)
		runner = perftest.NewBaselineRunner(runnerCfg, vmstarterService, dataRetrieverService)
	} else if *rtype == "tpm" {
		templateName := "256-tpm"
		runnerCfg := perftest.NewTestRunnerConfig(&baseTemplateData, templateName, *tempPath, *resultPath)
		tpmPath := dirutil.JoinPath(*tempPath, "tpm")
		err = dirutil.EnsureDirectory(tpmPath)
		if err != nil {
			fmt.Println("Could not create temp tpm directory")
			panic(err)
		}
		tpminst := tpminstantiator.NewTpmInstantiatorServiceWithBasePath(tpmPath)
		runner = perftest.NewTpmRunner(runnerCfg, vmstarterService, dataRetrieverService, tpminst)
	} else if *rtype == "pool" {
		templateName := "256-tpm"
		runnerCfg := perftest.NewTestRunnerConfig(&baseTemplateData, templateName, *tempPath, *resultPath)
		tpmPath := dirutil.JoinPath(*tempPath, "tpm")
		err = dirutil.EnsureDirectory(tpmPath)
		if err != nil {
			fmt.Println("Could not create temp tpm directory")
			panic(err)
		}

		// TODO organise handle of input in other way
		if *inum < 0 || *inum > 1001 {
			panic("Invalid number of VM instances: '" + strconv.Itoa(*inum) + "'.")
		}
		pool, err := perftest.NewTpmPool(*inum, tpmPath)
		if err != nil {
			fmt.Println("Could not create temp tpm pool")
			panic(err)
		}
		runner = perftest.NewTpmPoolRunner(runnerCfg, vmstarterService, dataRetrieverService, pool, *inum)
	} else {
		panic("Invalid performance test type: '" + *rtype + "'.")
	}

	// execute
	fmt.Println("Running perf test")
	err = perftestExecutor.RunPerftest(runner)
	if err != nil {
		fmt.Println("Error occurred running perftest")
		panic(err)
	}

}
