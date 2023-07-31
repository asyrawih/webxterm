package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"webxterm/internal/manager"

	"github.com/gorilla/mux"
)

// main function  
func main() {
	tm := manager.NewTTYServerManagar()

	server1 := manager.NewServer("ws", manager.Options{
		Host: "0.0.0.0",
		Port: "3000",
	})

	server2 := manager.NewServer("sw", manager.Options{
		Host: "0.0.0.0",
		Port: "3001",
	})

	server3 := manager.NewServer("sw", manager.Options{
		Host: "0.0.0.0",
		Port: "3002",
	})

	tm.AddServer(server1)
	tm.AddServer(server2)
	tm.AddServer(server3)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	go tm.Serve()

	<-signalChan
}

func NewhttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/spawn", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "hello")
	})
}
