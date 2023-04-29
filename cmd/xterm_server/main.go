package main

import (
	"fmt"
	"net/http"

	"webxterm/pkg/xterm"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

// main function  
func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Add("Content-type", "application/json")
		if _, err := w.Write([]byte("Hello World")); err != nil {
			printError(err)
		}
	})

	r.HandleFunc("/ws", xterm.HandleXtermConnection())

	listenOnAddress := fmt.Sprintf("%s:%d", "localhost", 3000)
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
