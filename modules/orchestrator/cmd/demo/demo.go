package main

import (
	"fmt"
	"os"
	"time"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
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
	time.Sleep(5 * time.Second)

	fmt.Printf("Current status %d\n", fce.Status())
	fmt.Printf("Current logs:\n%s", fce.Logs())
}
