package main

import (
	"webxterm/internal/manager"
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
	tm.Serve()
	select {}
}
