package main

import (
	db2 "bybarcode/internal/db"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/urfave/cli"

	"bybarcode/internal/config"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msgf(".env file load error: %s", err)
		return
	}

	cfg := config.NewCliConfig()

	db, err := db2.NewConnect("pgx", cfg.DBDsn)
	if err != nil {
		logger.Fatal().Msgf("db connection error: %s", err)
		return
	}

	app := cli.NewApp()
	app.Name = "Bybarcode CLI"
	app.Usage = "A simple CLI for bybarcode application"
	app.Commands = []cli.Command{
		{
			Name:    "load",
			Aliases: []string{"l"},
			Usage:   "Load products data to DB from file.",
			Action:  helloCommand,
		},
	}

	if err = app.Run(os.Args); err != nil {
		logger.Fatal().Msgf("CLI error: %s", err)
	}
}

func helloCommand(c *cli.Context) error {
	fmt.Println("LOAD!")
	return nil
}
