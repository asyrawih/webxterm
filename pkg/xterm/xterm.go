package xterm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/creack/pty"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

/*
Returns a function that handles an xterm.js connection by creating a PTY and
forwarding its input and output through the websocket connection. The function
takes an http.ResponseWriter and an http.Request as arguments.

The returned function upgrades the connection to a websocket connection, sets
up a ping handler, and starts three goroutines that handle the connection. The
first goroutine sends a ping message to the client every 10 seconds. The second
goroutine reads from the PTY and writes to the websocket connection. The third
goroutine reads from the websocket connection and writes to the PTY.
The function returns no values.
*/
func HandleXtermConnection() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader := getConnectionUpgrader()
		connection, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Err(err).Msg("")
		}

		command := exec.Command("bash")
		command.Env = os.Environ()
		tty, err := pty.Start(command)
		if err != nil {
			message := fmt.Sprintf("failed to start tty: %s", err)
			if err := connection.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Err(err).Msg("")
			}
			return
		}

		defer func() {
			if err := command.Process.Kill(); err != nil {
				log.Err(err).Msg("Error Kill The Process")
			}

			if _, err := command.Process.Wait(); err != nil {
				log.Err(err).Msg("")
			}

			if err := tty.Close(); err != nil {
				log.Err(err).Msg("Error Close The TTY")
			}

			if err := connection.Close(); err != nil {
				log.Err(err).Msg("Error Close The Connection")
			}
		}()

		connection.SetPingHandler(func(_ string) error {
			return nil
		})

		var connectionClosed bool
		var waiter sync.WaitGroup
		waiter.Add(1)

		go func() {
			for {
				if err := connection.WriteMessage(websocket.PingMessage, []byte("ping")); err != nil {
					log.Err(err).Msg("Error Ping the Connection")
					return
				}
				time.Sleep(10 * time.Second)
				log.Info().Msg("Recv Ping")
			}
		}()

		// Send the buffer from tty into websocket connection
		go func() {
			errorCounter := 0
			for {
				if errorCounter >= 10 {
					waiter.Done()
					break
				}
				buffer := make([]byte, 4028)
				n, err := tty.Read(buffer)
				if err != nil {
					if errConn := connection.WriteMessage(websocket.TextMessage, []byte(err.Error())); errConn != nil {
						log.Err(err).Msg("Error Write the TTY To xtermjs")
					}
					log.Err(err).Msg("Error Read the TTY")
					return
				}

				if err := connection.WriteMessage(websocket.BinaryMessage, buffer[:n]); err != nil {
					errorCounter++
					log.Err(err).Msg("")
					continue
				}
				errorCounter = 0
			}
		}()

		go func() {
			for {
				messageType, data, err := connection.ReadMessage()
				log.Info().Msg("Repeated")
				if err != nil {
					if !connectionClosed {
						log.Warn().Msgf("failed to get next reader: %s", err)
						return
					}
					log.Err(err).Msg("")
				}

				dataLength := len(data)

				dataBuffer := bytes.Trim(data, "\x00")

				dataType, ok := WebsocketMessageType[messageType]
				if !ok {
					dataType = "unknown"
				}

				log.Info().Msgf(
					"received %s (type: %v) message of size %v byte(s) from xterm.js with key sequence: %v",
					dataType,
					messageType,
					dataLength,
					dataBuffer,
				)

				if dataLength == -1 { // invalid
					log.Warn().Msg("failed to get the correct number of bytes read, ignoring message")
					continue
				}

				if messageType == websocket.BinaryMessage {
					if dataBuffer[0] == 1 {
						ttySize := &TTYSize{}
						resizeMessage := bytes.Trim(dataBuffer[1:], " \n\r\t\x00\x01")
						if err = json.Unmarshal(resizeMessage, ttySize); err != nil {
							log.Warn().Msgf(
								"failed to unmarshal received resize message '%s': %s",
								string(resizeMessage),
								err,
							)
							continue
						}
						log.Info().Msgf("resizing tty to use %v rows and %v columns...", ttySize.Rows, ttySize.Cols)
						if err = pty.Setsize(tty, &pty.Winsize{
							Rows: ttySize.Rows,
							Cols: ttySize.Cols,
						}); err != nil {
							log.Warn().Msgf("failed to resize tty, error: %s", err)
						}
						continue
					}
				}

				// write to tty
				bytesWritten, err := tty.Write(dataBuffer)
				if err != nil {
					log.Warn().Msg(fmt.Sprintf("failed to write %v bytes to tty: %s", len(dataBuffer), err))
					continue
				}
				log.Info().Msgf("%v bytes written to tty...", bytesWritten)
			}
		}()

		waiter.Wait()
		log.Info().Msg("clossing the connection")
		connectionClosed = true
	}
}

/*
getConnectionUpgrader returns a websocket.Upgrader with a CheckOrigin function that always returns true and a HandshakeTimeout set to 0.

Returns:

	websocket.Upgrader: An object that can upgrade an HTTP connection to a WebSocket connection.
*/
func getConnectionUpgrader() websocket.Upgrader {
	return websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		HandshakeTimeout: 0,
	}
}
