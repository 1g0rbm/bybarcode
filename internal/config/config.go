package config

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultBarcodeFilePath = "/data/products_data_all.csv"
	defaultAddress         = "127.0.0.1:8080"
)

type ApiConfig struct {
	DBDsn    string
	Address  string
	Gateway  string
	BotToken string
	TgApiUrl string
}

type CliConfig struct {
	DBDsn           string
	BarcodeFilePath string
}

type BotConfig struct {
	Token       string
	DBDsn       string
	TgWebAppUrl string
}

func NewBotConfig() *BotConfig {
	return &BotConfig{
		Token:       getEnvString("TELEGRAM_BOT_TOKEN", ""),
		DBDsn:       getEnvString("POSTGRESQL_URL", ""),
		TgWebAppUrl: getEnvString("TELEGRAM_BOT_WEB_APP_URL", ""),
	}
}

func NewApiConfig() *ApiConfig {
	return &ApiConfig{
		DBDsn:    getEnvString("POSTGRESQL_URL", ""),
		Address:  getEnvString("APO_ADDRESS", defaultAddress),
		Gateway:  getEnvString("API_GATEWAY", ""),
		BotToken: getEnvString("TELEGRAM_BOT_TOKEN", ""),
		TgApiUrl: getEnvString("TELEGRAM_API_URL", ""),
	}
}

func NewCliConfig() *CliConfig {
	return &CliConfig{
		DBDsn:           getEnvString("POSTGRESQL_URL", ""),
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
