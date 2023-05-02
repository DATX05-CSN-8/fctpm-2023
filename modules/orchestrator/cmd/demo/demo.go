package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/dirutil"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/internal/firecracker"
	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/tpminstantiator"
)

func main() {
	tpmPath := os.Getenv("TPM_PATH")
	if len(tpmPath) == 0 {
		fmt.Println("TPM_PATH environmnent variable needs to be specified")
		return
	}
	fcTmplPath := os.Getenv("FIRECRACKER_TEMPLATE_PATH")
	if len(fcTmplPath) == 0 {
		fmt.Println("FIRECRACKER_TEMPLATE_PATH variable needs to be specified")
		return
	}

	service := tpminstantiator.NewTpmInstantiatorService()
	tempDirPath, err := dirutil.CreateTempDir(tpmPath)
	if err != nil {
		panic(err)
	}
	// create swtpm
	instance, err := service.Create(tempDirPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer service.Destroy(instance)

	fmt.Printf("SWTPM socket path: %s\n", instance.SocketPath)

	// create firecracker config
	data := firecracker.SimpleTemplateData{
		KernelImagePath: "/home/alex/fctpm-2023/vm-image/out/fc-image-kernel",
		InitRdPath:      "/home/alex/fctpm-2023/vm-image/out/fc-image-initrd.img",
		TpmSocket:       instance.SocketPath,
	}
	err = firecracker.NewFirecrackerConfig("simple", data, fcTmplPath)
	if err != nil {
		fmt.Println(err)
		return
	}
	// wait for input
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to stop swtpm...")
	_, _ = r.ReadString('\n')

}
