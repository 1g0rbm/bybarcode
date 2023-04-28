package main

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"

	"bybarcode/internal/api"
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

	apiApp := api.NewAppApi(conn, cfg, logger)
	if err = apiApp.Run(); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}
