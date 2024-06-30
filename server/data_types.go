package server

import (
	"net"
	"sync"
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
