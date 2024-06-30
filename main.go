package main

import "server"

func main() {
	port := server.GetPort()
	server := server.New(server.Config{
		Host: "localhost",
		Port: port,
	})
	server.Run()
}
