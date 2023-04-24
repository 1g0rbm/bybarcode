package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	"bybarcode/internal/config"
	"bybarcode/internal/db"
)

func main() {
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger()

	if err := godotenv.Load(); err != nil {
		logger.Fatal().Msgf(".env file load error: %s", err)
		return
	}

	app := &cli.App{
		Name:  "Bybarcode CLI",
		Usage: "A simple CLI for bybarcode application",
		Commands: []*cli.Command{
			{
				Name:    "load",
				Aliases: []string{"l"},
				Usage:   "Load products data to DB from file.",
				Action:  helloCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}

func helloCommand(c *cli.Context) error {
	cfg := config.NewCliConfig()

	conn, err := db.NewConnect("pgx", cfg.DBDsn)
	defer func(conn *db.Connect) {
		err = conn.Close()
	}(&conn)
	if err != nil {
		return err
	}

	root, err := os.Getwd()
	if err != nil {
		return err
	}

	file, err := os.Open(root + cfg.BarcodeFilePath)
	if err != nil {
		return err
	}

	brands := map[string]int{}
	categories := map[string]int{}

	reader := csv.NewReader(file)
	reader.Comma = '\t'
	_, err = reader.Read()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		fmt.Println(line)
		brandName := line[6]
		if brandName == "" {
			brandName = "unknown brand"
		}

		brandId, ok := brands[brandName]
		if !ok {

			brandId, err = conn.CreateBrand(ctx, brandName)
			if err != nil {
				return err
			}
			brands[brandName] = brandId
		}

		categoryName := line[4]
		if categoryName == "" {
			categoryName = "unknown category"
		}

		categoryId, ok := categories[categoryName]
		if !ok {
			categoryId, err = conn.CreateCategory(ctx, categoryName)
			if err != nil {
				return err
			}
			categories[categoryName] = categoryId
		}

		fmt.Println(categoryId)
		fmt.Println(brandId)
		fmt.Println(categories)
		fmt.Println(brands)

		break
	}

	return err
}
