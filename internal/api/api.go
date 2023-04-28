package api

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"

	"bybarcode/internal/config"
	"bybarcode/internal/db"
)

type AppApi struct {
	db     db.Connect
	cfg    *config.ApiConfig
	logger zerolog.Logger
	router *chi.Mux
	server *http.Server
}

func NewAppApi(db db.Connect, cfg *config.ApiConfig, logger zerolog.Logger) *AppApi {
	r := chi.NewRouter()

	api := &AppApi{
		db:     db,
		cfg:    cfg,
		logger: logger,
		router: r,
		server: &http.Server{
			Addr:    cfg.Address,
			Handler: r,
		},
	}

	api.router.Use(middleware.RequestID)
	api.router.Use(middleware.RealIP)
	api.router.Use(middleware.Logger)
	api.router.Use(middleware.Recoverer)

	api.router.Get("/api/v1/ping", api.pingHandler)

	return api
}

func (aa *AppApi) Run() error {
	return aa.server.ListenAndServe()
}

func (aa *AppApi) pingHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	_, err := w.Write([]byte(`{"status": "ok"}`))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}
