package api

import (
	"net/http"

	"github.com/rs/zerolog"
)

type response struct {
	logger zerolog.Logger
}

func newResponse(l zerolog.Logger) *response {
	return &response{
		logger: l,
	}
}

func (r response) json(w http.ResponseWriter, status int, body []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(body); err != nil {
		r.logger.Fatal().Msg(err.Error())
	}

	r.logger.Debug().Msgf("Send response with headers %s and body %s", w.Header(), string(body))
}
