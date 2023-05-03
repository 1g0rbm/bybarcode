package api

import (
	"net/http"
)

func (h *handlers) pingHandler(w http.ResponseWriter, r *http.Request) {
	_, err := w.Write([]byte(`{"status": "ok"}`))
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte(err.Error()))
	}
}
