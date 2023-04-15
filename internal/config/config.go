package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultBarcodeFilePath = "../../data/tmp/products_data_all.csv"
	defaultDBDsn           = ""
)

type ApiConfig struct {
	DBDsn string
}

type CliConfig struct {
	DBDsn           string
	BarcodeFilePath string
}

func NewApiConfig() *ApiConfig {
	return &ApiConfig{
		DBDsn: getEnvString("DATABASE_DSN", defaultDBDsn),
	}
}

func NewCliConfig() *CliConfig {
	return &CliConfig{
		DBDsn:           getEnvString("POSTGRESQL_URL", defaultDBDsn),
		BarcodeFilePath: getEnvString("BARCODE_FILE_PATH", defaultBarcodeFilePath),
	}
}

func getEnvString(name string, defaultValue string) string {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	return value
}

func getEnvDuration(name string, defaultValue time.Duration) time.Duration {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	normalizedValue, err := time.ParseDuration(value)
	if err != nil {
		return defaultValue
	}

	return normalizedValue
}

func getEnvBool(name string, defaultValue bool) bool {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	b, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}
	return b
}

func getEnvInt(name string, defaultValue int) int {
	value, ok := os.LookupEnv(name)
	if !ok {
		return defaultValue
	}

	i, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return defaultValue
	}

	return int(i)
}
