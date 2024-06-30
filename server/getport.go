package server

import (
	"fmt"
	"os"
)

func GetPort() string {
	port := ""
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(0)
	}
	if len(os.Args) == 1 {
		port = "8989"
	} else {
		port = os.Args[1]
	}
	return port
}
