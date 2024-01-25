package internal

import (
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

func GetWANIP() string {
	res, err := http.Get("https://ipv4.icanhazip.com")

	if err != nil {
		log.Fatal().Err(err).Msg("")
	}

	resBody, _ := io.ReadAll(res.Body)

	return strings.TrimSpace(string(resBody))
}

func WanIPCheck(checkInterval int) {
	for {
		time.Sleep(time.Duration(checkInterval) * time.Second)
		log.Info().Str("WAN_IP", GetWANIP()).Msg("WAN IP check")
	}
}
