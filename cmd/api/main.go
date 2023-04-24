package main

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"bybarcode/internal/config"
	"bybarcode/internal/db"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msgf(".env file load error: %s", err)
		return
	}

	cfg := config.NewApiConfig()

	conn, err := db.NewConnect("pgx", cfg.DBDsn)
	defer func(conn *db.Connect) {
		err = conn.Close()
	}(&conn)
	if err != nil {
		logger.Fatal().Msg(err.Error())
	}

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/api/v1/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, err := w.Write([]byte(`{"status": "ok"}`))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})

	server := &http.Server{
		Addr:    cfg.Address,
		Handler: r,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}
