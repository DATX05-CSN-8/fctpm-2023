package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/perftest"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmdata"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmexecution"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vmstarter"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/resourcepool"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/tpminstantiator"
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
	fcPath := flag.String("firecracker-bin", genDefaultFCBinPath(), "File path to the firecracker binary that should be used")
	resultPath := flag.String("result-path", "output.csv", "Path to CSV file to create")
	tempPath := flag.String("temp-path", "/tmp/firecracker-perftest", "Path to temporary data directory")
	clean := flag.Bool("clean", false, "Clean the output database and csv")
	bootLogPath := flag.String("boot-log-path", "", "(optional) path to output boot logs")
	rtype := flag.String("type", "baseline", "Type of performance test to run. Either 'baseline' or 'tpm'")
	totalVms := flag.Int("total-vms", 20, "The total number of VMs to run as part of the perf test scenario.")
	parallelism := flag.Int("parallelism", 1, "The number of VMs to run in parallel as part of the scenario")
	kernelPath := flag.String("kernel-path", "/home/melker/fctpm-2023/vm-image/out/fc-image-kernel", "Path to Firecracker kernel")
	initPath := flag.String("init-path", "/home/melker/fctpm-2023/vm-image/out/fc-image-initrd.img", "Path to Firecracker init")
	flag.Parse()
	var err error
	if *clean {
		err = dirutil.RemoveFileIfExists(*resultPath)
		if err != nil {
			fmt.Println("Error occurred removing result file")
			panic(err)
		}
	}

	fcClient := firecracker.NewFirecrackerClient(*fcPath)
	vminfoRepo := vminfo.NewMapRepository()
	vmExecRepo := vmexecution.NewMapRepository()
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

	perftestExecutor := perftest.NewPerftestExecutor(*totalVms, *parallelism)
	baseTemplateData := firecracker.SimpleTemplateData{
		KernelImagePath: *kernelPath,
		InitRdPath:      *initPath,
	}

	var runner perftest.PerftestRunner
	templateName := "default"
	if *rtype == "baseline" {
		runnerCfg := perftest.NewTestRunnerConfig(&baseTemplateData, templateName, *tempPath, *resultPath)
		runner = perftest.NewBaselineRunner(runnerCfg, vmstarterService, dataRetrieverService)
	} else if *rtype == "tpm" {
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
		runnerCfg := perftest.NewTestRunnerConfig(&baseTemplateData, templateName, *tempPath, *resultPath)
		tpmPath := dirutil.JoinPath(*tempPath, "tpm")
		err = dirutil.EnsureDirectory(tpmPath)
		if err != nil {
			fmt.Println("Could not create temp tpm directory")
			panic(err)
		}
		tpminst := tpminstantiator.NewTpmInstantiatorServiceWithBasePath(tpmPath)
		pool, err := resourcepool.NewResourcePool[tpminstantiator.TpmInstance](*totalVms, tpminst)
		if err != nil {
			fmt.Println("Could not create temp tpm pool")
			panic(err)
		}
		runner = perftest.NewTpmRunner(runnerCfg, vmstarterService, dataRetrieverService, pool)
	} else if *rtype == "pool-alt" {
		runnerCfg := perftest.NewTestRunnerConfig(&baseTemplateData, templateName, *tempPath, *resultPath)
		tpmPath := dirutil.JoinPath(*tempPath, "tpm")
		err = dirutil.EnsureDirectory(tpmPath)
		if err != nil {
			fmt.Println("Could not create temp tpm directory")
			panic(err)
		}
		tpminst := tpminstantiator.NewTpmInstantiatorServiceWithBasePath(tpmPath)
		pool, err := resourcepool.NewAltResourcePool[tpminstantiator.TpmInstance](*totalVms, tpminst)
		if err != nil {
			fmt.Println("Could not create temp tpm pool")
			panic(err)
		}
		runner = perftest.NewTpmRunner(runnerCfg, vmstarterService, dataRetrieverService, pool)
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
	fmt.Println("Stopping perf test execution")

	<-time.After(10 * time.Second)

}
