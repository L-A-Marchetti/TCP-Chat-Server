package server

import (
	"log"
	"os"
)

func SaveLogs(s string) {
	// Open logs.txt
	file, err := os.OpenFile("server/logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	// Append the logs.
	_, err = file.WriteString(s)
	if err != nil {
		log.Fatal(err)
	}
}
