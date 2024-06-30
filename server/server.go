package server

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func (server *Server) Run() {
	// Start listening for incoming TCP connections on the specified host and port.
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	log.Printf("Listening on port %s", server.port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		// Accept incoming connection from listener.
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		// Create a new Client instance for the accepted connection
		client := Client{
			conn: conn,
		}
		// Welcome message.
		welcome := PrintWelcome()
		client.conn.Write([]byte(welcome))
		// Ask for pseudo and erase the end of line.
		client.conn.Write([]byte("[ENTER YOUR PSEUDO]: "))
		reader := bufio.NewReader(client.conn)
		pseudo, _ := reader.ReadString('\n')
		client.pseudo = strings.TrimRight(pseudo, "\n")
		// Lock the mutex to synchronize access to the server's clients list.
		server.mutex.Lock()
		// Append the new client to the server's list of clients.
		server.clients = append(server.clients, client)
		// Unlock the mutex to allow other goroutines to access the clients list.
		server.mutex.Unlock()
		// Handle client messages in a goroutine.
		go server.handleClient(client)
	}
}
