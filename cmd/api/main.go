package main

import (
	"bybarcode/internal/bot"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"net/http"
	"os"

	"bybarcode/internal/auth"
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

	r.Post("/api/v1/auth", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var uid auth.UserId
		if err := uid.Decode(r.Body); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				panic(err)
			}
			return
		}

		acc, err := conn.FindAccountById(r.Context(), uid.Value)
		if errors.As(err, &pgx.ErrNoRows) {
			w.WriteHeader(http.StatusBadRequest)
			_, err = w.Write([]byte(err.Error()))
			if err != nil {
				panic(err)
			}
			return
		}

		tgSender := bot.NewSender(cfg.BotToken, cfg.TgApiUrl)
		if err = tgSender.SendMessage(acc.ID, "/start"); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte("internal server error"))
			if err != nil {
				panic(err)
			}
			return
		}

		//session := auth.Session{
		//	ID:           uuid.New(),
		//	Token:        uuid.New(),
		//	RefreshToken: uuid.New(),
		//	AccountID:    uid.Value,
		//	ExpireAt:     time.Now().Add(24 * time.Hour),
		//	CreatedAt:    time.Now(),
		//	UpdatedAt:    time.Now(),
		//}
		//
		//if err := conn.CreateSession(r.Context(), session); err != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//	_, err = w.Write([]byte(err.Error()))
		//	if err != nil {
		//		panic(err)
		//	}
		//	return
		//}
		//
		//b, err := session.Encode()
		//if err != nil {
		//	w.WriteHeader(http.StatusInternalServerError)
		//	_, err = w.Write([]byte(err.Error()))
		//	if err != nil {
		//		panic(err)
		//	}
		//	return
		//}

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(""))
		if err != nil {
			panic(err)
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
