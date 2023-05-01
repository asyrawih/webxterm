package httphandler

import (
	"net"
	"net/http"

	"github.com/rs/zerolog/log"
)

func GetIP() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, _ *http.Request) {
		c, err := net.Dial("udp", "8.8.8.8:80")
		if err != nil {
			log.Err(err).Msg("")
		}
		defer c.Close()

		ipLocal := c.LocalAddr().(*net.UDPAddr)

		_, err = w.Write([]byte(ipLocal.IP.String()))
		if err != nil {
			log.Err(err).Msg("")
		}

	}
}
