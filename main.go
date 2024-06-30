package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

type Server struct {
	host    string
	port    string
	clients []Client
	mutex   sync.Mutex
}

type Client struct {
	conn   net.Conn
	pseudo string
}

type Config struct {
	Host string
	Port string
}

type Logs struct {
	Time    string
	Pseudo  string
	Message string
}

func New(config Config) Server {
	return Server{
		host: config.Host,
		port: config.Port,
	}
}

var historic []Logs

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

func (server *Server) broadcastMessage(sender Client, message string) {
	// Locks the mutex to synchronize access.
	server.mutex.Lock()
	// Ensures the mutex is unlocked after the function exits, preventing deadlock.
	defer server.mutex.Unlock()
	// If sender has a pseudo, broadcast message to all clients with the sender name.
	if sender.pseudo != "" {
		for _, client := range server.clients {
			// Append the message to the historic.
			historic = append(historic, Logs{time.Now().Format("2006-01-02 15:04:05"), sender.pseudo, message})
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

func main() {
	port := GetPort()
	server := New(Config{
		Host: "localhost",
		Port: port,
	})
	server.Run()
}
