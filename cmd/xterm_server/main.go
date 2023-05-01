package main

import (
	"fmt"
	"net/http"

	"webxterm/internal/httphandler"
	"webxterm/internal/websockethandler"
	"webxterm/pkg/xterm"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// main function  
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/ws/{name}", xterm.HandleXtermConnection())
	r.HandleFunc("/ip", httphandler.GetIP())
	r.HandleFunc("/ping", websockethandler.PingServer("google.com"))

	listenOnAddress := fmt.Sprintf("%s:%d", "0.0.0.0", 3000)
	log.Info().Msg("Listening on " + listenOnAddress)

	server := http.Server{
		Addr:    listenOnAddress,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		printError(err)
	}
}

func printError(err error) {
	log.Error().Err(err).Msg("")
}
