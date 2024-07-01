package server

import (
	"fmt"
	"log"
	"time"
)

func (server *Server) broadcastMessage(sender Client, message string) {
	// Locks the mutex to synchronize access.
	server.mutex.Lock()
	// Ensures the mutex is unlocked after the function exits, preventing deadlock.
	defer server.mutex.Unlock()
	// If sender has a pseudo, broadcast message to all clients with the sender name.
	if sender.pseudo != "" {
		// Host log.
		log.Printf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), sender.pseudo, message)
		// Append the message to the historic.
		historic = append(historic, Logs{time.Now().Format("2006-01-02 15:04:05"), sender.pseudo, message})
		// Save the logs in a file.
		SaveLogs(fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), sender.pseudo, message))
		for _, client := range server.clients {
			// Send formatted message to each client.
			client.conn.Write([]byte(fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), sender.pseudo, message)))
		}
	} else {
		// Broadcast message to all clients (used for notifications).
		for _, client := range server.clients {
			client.conn.Write([]byte(message))
		}
	}
}
