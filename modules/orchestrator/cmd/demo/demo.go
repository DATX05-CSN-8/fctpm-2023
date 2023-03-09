package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/DATX05-CSN-8/fctpm-2023/modules/orchestrator/pkg/tpminstantiator"
)

func main() {
	swtpmBin := os.Getenv("SWTPM_BIN")
	if len(swtpmBin) == 0 {
		fmt.Println("SWTPM_BIN environment variable needs to be specified")
		return
	}
	tpmPath := os.Getenv("TPM_PATH")
	if len(tpmPath) == 0 {
		fmt.Println("TPM_PATH environmnent variable needs to be specified")
		return
	}

	service := tpminstantiator.NewTpmInstantiatorService(swtpmBin, tpmPath)
	instance, err := service.Create()
	if err != nil {
		fmt.Println(err)
		return
	}
	r := bufio.NewReader(os.Stdin)
	fmt.Print("Press enter to stop swtpm...")
	_, _ = r.ReadString('\n')

	service.Destroy(instance)

}
