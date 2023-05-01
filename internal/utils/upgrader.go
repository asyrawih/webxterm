package utils

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// GetConnectionUpgrader function  
func GetConnectionUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout: 0,
	}
}
