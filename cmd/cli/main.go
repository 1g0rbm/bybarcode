package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/urfave/cli/v2"

	"bybarcode/internal/config"
	"bybarcode/internal/db"
	"bybarcode/internal/products"
)

const (
	lineChBufferSize = 2000
	numWorkers       = 16
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
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "workers",
						Value: numWorkers,
						Usage: "Number of worker goroutines to use",
					},
				},
				Action: loadCommand,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logger.Fatal().Msg(err.Error())
	}
}

func loadCommand(c *cli.Context) error {
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
	defer func(file *os.File) {
		err = file.Close()
	}(file)

	lineCh := make(chan []string, lineChBufferSize)
	errCh := make(chan error)

	go func() {
		reader := csv.NewReader(file)
		reader.Comma = '\t'
		defer close(lineCh)

		_, err = reader.Read()
		if err != nil {
			errCh <- err
		}

		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				fmt.Printf("****\nPARSE ERR: %s\n****\n", err)
				continue
			}

			lineCh <- line
		}
		return
	}()

	var wg sync.WaitGroup
	for i := 0; i < 20; i += 1 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case line := <-lineCh:
					if err = handleLine(c.Context, &conn, line); err != nil {
						errCh <- err
					}
				case <-c.Context.Done():
					return
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(errCh)
	}()

	for {
		select {
		case err = <-errCh:
			return err
		}
	}
}

func handleLine(
	ctx context.Context,
	conn *db.Connect,
	line []string,
) error {
	brandName := line[6]
	if brandName == "" {
		brandName = "unknown brand"
	}

	brandId, err := conn.CreateBrand(ctx, brandName)
	if err != nil {
		return err
	}

	categoryName := line[4]
	if categoryName == "" {
		categoryName = "unknown category"
	}

	categoryId, err := conn.CreateCategory(ctx, categoryName)
	if err != nil {
		return err
	}

	id, err := strconv.ParseInt(line[0], 10, 64)
	if err != nil {
		return err
	}

	barcode, err := strconv.ParseInt(line[1], 10, 64)
	if err != nil {
		return err
	}

	p := products.Product{
		ID:         id,
		Upcean:     barcode,
		Name:       line[2],
		CategoryId: int64(categoryId),
		BrandId:    int64(brandId),
	}

	_, err = conn.CreateProduct(ctx, p)
	if err != nil {
		return err
	}

	fmt.Printf("Loaded product %v\n", p)

	return nil
}
