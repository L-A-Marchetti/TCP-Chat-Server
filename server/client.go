package server

import (
	"bufio"
	"fmt"
	"log"
)

func (server *Server) handleClient(client Client) {
	// Host log.
	log.Printf("Client %s connected\n", client.pseudo)
	// Notify other clients about new connection
	server.broadcastMessage(Client{}, fmt.Sprintf("%s has joind the chat...\n", client.pseudo))
	// Display historic of the conversations.
	for _, log := range historic {
		client.conn.Write([]byte("[" + log.Time + "][" + log.Pseudo + "]: " + log.Message))
	}
	// Remove client from server's list when client disconnects.
	defer func() {
		server.mutex.Lock()
		defer server.mutex.Unlock()
		for i, c := range server.clients {
			if c == client {
				server.clients = append(server.clients[:i], server.clients[i+1:]...)
				break
			}
		}
		client.conn.Close()
	}()
	// Read client input.
	reader := bufio.NewReader(client.conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected\n", client.pseudo)
			// Broadcast disconnect message to remaining clients
			server.broadcastMessage(Client{}, fmt.Sprintf("%s has left the chat...\n", client.pseudo))
			break
		}
		// Erase the input line for the client.
		client.conn.Write([]byte("\033[F"))
		// Broadcasts the message to all clients indicating the sender.
		server.broadcastMessage(client, message)
	}
}
