package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func (h *handlers) getStatistic(w http.ResponseWriter, r *http.Request) {
	dateFromStr := chi.URLParam(r, "date_from")
	dateToStr := chi.URLParam(r, "date_to")

	dateFrom, err := time.Parse("2006-01-02T15:04:05", dateFromStr)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusBadRequest, []byte("Invalid date_from format"))
		return
	}

	dateTo, err := time.Parse("2006-01-02T15:04:05", dateToStr)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusBadRequest, []byte("Invalid date_to format"))
		return
	}

	list, err := h.db.GetStatistic(r.Context(), dateFrom, dateTo)
	if err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	if err := json.NewEncoder(w).Encode(list); err != nil {
		h.logger.Error().Msg(err.Error())
		h.response.json(w, http.StatusInternalServerError, []byte("internal server error"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
