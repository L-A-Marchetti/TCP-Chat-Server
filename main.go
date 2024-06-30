package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type Server struct {
	host    string
	port    string
	clients []*Client
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

func New(config *Config) *Server {
	return &Server{
		host: config.Host,
		port: config.Port,
	}
}

var historic []Logs

func (server *Server) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", server.host, server.port))
	log.Printf("Listening on port %s", server.port)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}

		client := &Client{
			conn: conn,
		}
		// Welcome message.
		welcome := PrintWelcome()
		client.conn.Write([]byte(welcome))
		// Ask for pseudo and add client to server's list
		client.conn.Write([]byte("[ENTER YOUR PSEUDO]: "))
		reader := bufio.NewReader(client.conn)
		pseudo, _ := reader.ReadString('\n')
		client.pseudo = strings.TrimRight(pseudo, "\n")

		server.mutex.Lock()
		server.clients = append(server.clients, client)
		server.mutex.Unlock()

		// Handle client messages in a goroutine
		go server.handleClient(client)
	}
}

func (server *Server) handleClient(client *Client) {
	// Display historic of the conversations.
	for _, log := range historic {
		client.conn.Write([]byte("[" + log.Time + "][" + log.Pseudo + "]: " + log.Message))
	}
	log.Printf("Client %s connected\n", client.pseudo)
	defer func() {
		// Remove client from server's list when client disconnects
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

	reader := bufio.NewReader(client.conn)

	// Notify other clients about new connection
	server.broadcastMessage(nil, fmt.Sprintf("%s has joind the chat...\n", client.pseudo))

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Client %s disconnected\n", client.pseudo)
			// Broadcast disconnect message to remaining clients
			server.broadcastMessage(nil, fmt.Sprintf("%s has left the chat...\n", client.pseudo))
			break
		}
		client.conn.Write([]byte("\033[F")) // Erase the input line.
		server.broadcastMessage(client, message)
	}
}

func (server *Server) broadcastMessage(sender *Client, message string) {
	server.mutex.Lock()
	defer server.mutex.Unlock()

	if sender != nil {
		for _, client := range server.clients {
			historic = append(historic, Logs{time.Now().Format("2006-01-02 15:04:05"), sender.pseudo, message})
			client.conn.Write([]byte(fmt.Sprintf("[%s][%s]: %s", time.Now().Format("2006-01-02 15:04:05"), sender.pseudo, message)))
		}
	} else {
		// Broadcast message to all clients (used for notifications)
		for _, client := range server.clients {
			client.conn.Write([]byte(message))
		}
	}
}

func main() {
	port := ""
	if len(os.Args) > 2 {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return
	}
	if len(os.Args) == 1 {
		port = "8989"
	} else {
		port = os.Args[1]
	}
	server := New(&Config{
		Host: "localhost",
		Port: port,
	})
	server.Run()
}
