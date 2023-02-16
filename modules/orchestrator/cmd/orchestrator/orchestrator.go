package main

import (
	"flag"
	"fmt"
)

func main() {
	var apiSock = flag.String("api-sock", "/tmp/fctpm-orchestrator", "File path to the socket that should be listened to")
	flag.Parse()
	fmt.Println("api socket has value ", *apiSock)
}
