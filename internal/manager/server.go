package manager

import (
	"fmt"
	"log"
	"net/http"

	"webxterm/pkg/xterm"

	"github.com/gorilla/mux"
)

// TTYServer struct  
type TTYServer struct {
	Name string
	Options
}

// This Options Will Contains Information About Server
type Options struct {
	// Websocket Port
	Port string
	// Websocket Host
	Host string
}

type HandlerOpts struct {
	// Arguments is a list of strings to pass as arguments to the specified COmmand
	Arguments []string
	// Command is the path to the binary we should create a TTY for
	Command string
}

type TTYServerManager struct {
	TTYServers []*TTYServer
}

// NewTTYServerManagar function  
func NewTTYServerManagar() *TTYServerManager {
	return &TTYServerManager{}
}

// AddServer method  
func (m *TTYServerManager) AddServer(server *TTYServer) {
	m.TTYServers = append(m.TTYServers, server)
}

// GetServer method  
func (m *TTYServerManager) GetServer() []*TTYServer {
	return m.TTYServers
}

// Serve method  
// Server The Websocket Server
func (m *TTYServerManager) Serve() {
	for _, server := range m.TTYServers {
		go func(server *TTYServer) {
			r := mux.NewRouter()

			r.HandleFunc("/"+server.Name, xterm.HandleXtermConnection())

			listenOnAddress := fmt.Sprintf("%s:%s", server.Host, server.Port)
			log.Print("ListenAndServe on " + listenOnAddress)

			listener := http.Server{
				Addr:    listenOnAddress,
				Handler: r,
			}

			if err := listener.ListenAndServe(); err != nil {
				log.Fatal(err.Error())
			}
		}(server)
	}
}

// NewOptions function  
func NewOptions(host, port string) Options {
	return Options{
		Port: port,
		Host: host,
	}
}

// NewServer function  
func NewServer(Name string, opt Options) *TTYServer {
	return &TTYServer{
		Name:    Name,
		Options: opt,
	}
}
