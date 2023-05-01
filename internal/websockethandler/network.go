package websockethandler

import (
	"net/http"
	"os"
	"os/exec"
	"sync"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"

	"webxterm/internal/utils"
)

func PingServer(address string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := utils.GetConnectionUpgrader()
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Err(err).Msg("Failed to upgrade connection")
		}
		waiter := sync.WaitGroup{}
		waiter.Add(1)

		go SendPing(conn, &waiter)

		waiter.Wait()

	}
}

// SendPing function  
// Send ping to tty
func SendPing(conn *websocket.Conn, wg *sync.WaitGroup) {
	defer wg.Done()
	c := exec.Command("ping", "google.com")
	c.Env = os.Environ()
	tty, err := pty.Start(c)
	if err != nil {
		log.Err(err).Msg("Error Start TTY")
	}

	for i := 0; i < 10; i++ {
		buffer := make([]byte, 1024)
		n, err := tty.Read(buffer)
		if err != nil {
			log.Err(err).Msg("Error Read Buffer from TTY")
			break
		}

		if err := conn.WriteMessage(websocket.TextMessage, buffer[:n]); err != nil {
			log.Err(err).Msg("Error Write to websocket connection")
			break
		}
	}

	defer func() {
		log.Info().Msg("Closing TTY")
		if err := c.Process.Kill(); err != nil {
			log.Err(err).Msg("Error Kill The Process")
		}

		if err := c.Wait(); err != nil {
			log.Err(err).Msg("Error Kill The Process")
		}
		if err := conn.Close(); err != nil {
			log.Err(err).Msg("Error Kill The Process")
		}
	}()

}
