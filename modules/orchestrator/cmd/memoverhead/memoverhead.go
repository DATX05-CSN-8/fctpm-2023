package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/memoverhead"
)

func genDefaultFCBinPath() string {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return wd + "/../firecracker/bin/firecracker"
}

func main() {
	fcPath := flag.String("firecracker-bin", genDefaultFCBinPath(), "File path to the firecracker binary that should be used")
	var memsizes MemInput
	flag.Var(&memsizes, "mem-sizes", "Comma-separated list of memory sizes to run the memory overhead measurement for")
	rtype := flag.String("type", "baseline", "Type of performance test to run. Either 'baseline' or 'tpm'")
	resultPath := flag.String("result-path", "output.csv", "Path to CSV file to create")
	tempPath := flag.String("temp-path", "/tmp/firecracker-memoverhead", "Path to temporary data directory")
	kernel := flag.String("kernel", "NONE", "Path to kernel to run")
	initrd := flag.String("init-rd", "NONE", "Path to init rd to run")
	flag.Parse()

	if *kernel == "NONE" {
		panic("You need to specify the kernel command line argument")
	}
	if *initrd == "NONE" {
		panic("You need to specify the init-rd command line argument")
	}

	fcClient := firecracker.NewFirecrackerClientWithTimeout(*fcPath, 1*time.Hour)
	baseTemplateData := firecracker.SimpleTemplateData{
		KernelImagePath: *kernel,
		InitRdPath:      *initrd,
		MemSize:         128,
		TpmSocket:       "",
		BootArgs:        "panic=-1",
	}

	var runner memoverhead.MemoryOverheadRunner
	if *rtype == "baseline" {
		templateName := "default"
		cfg := memoverhead.NewMemoryOverheadConfig(&baseTemplateData, templateName, *tempPath)
		r, err := memoverhead.NewBaselineRunner(fcClient, cfg)
		if err != nil {
			panic(err)
		}
		runner = &r
	} else {
		// TODO add handling for tpm runner
		panic(fmt.Errorf("Invalid rtype specified"))
	}
	executor, err := memoverhead.NewMemoryOverheadExecutor(runner, resultPath)
	if err != nil {
		panic(err)
	}

	err = executor.RunWithMems(memsizes)
	if err != nil {
		panic(err)
	}
	fmt.Println("Sleeping to finish all processes")
	<-time.After(5 * time.Second)

}
