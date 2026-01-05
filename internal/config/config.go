package config

import (
    "os"
)

type Config struct {
    TelegramToken string
}

func Load() *Config {
    return &Config{
        TelegramToken: os.Getenv("TELEGRAM_TOKEN"),
    }
}
