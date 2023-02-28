package main

import (
	"fmt"
	"os"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/vminfo"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fcPath := wd + "/../firecracker/bin/firecracker"
	configPath := wd + "/../vm-start/fc-config.json"
	fc := firecracker.NewFirecrackerClient(fcPath)
	fce, err := fc.Start(configPath)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current status %d\n", fce.Status())
	fmt.Printf("Current logs:\n%s", fce.Logs())
	c := make(chan int)
	fce.Subscribe(func(status vminfo.Status) {
		fmt.Printf("Status: %d\n", status)
		c <- 1
	})
	<-c
}
