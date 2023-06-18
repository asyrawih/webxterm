package manager

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	options := NewOptions("127.0.0.1", "3000")
	server := NewServer("server-1", options)

	options2 := NewOptions("127.0.0.1", "2000")
	server2 := NewServer("server-2", options2)

	tm := NewTTYServerManagar()
	tm.AddServer(server)
	tm.AddServer(server2)

	t.Log(tm.TTYServers)
}
