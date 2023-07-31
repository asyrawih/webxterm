package main

import (
	"fmt"
	"net/http"

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

	tm.AddServer(server1)
	tm.AddServer(server2)

	go tm.Serve()

	select {}
}

func NewhttpServer() {
	r := mux.NewRouter()
	r.HandleFunc("/spawn", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintf(w, "hello")
	})
}
