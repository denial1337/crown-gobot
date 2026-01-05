package main

import (
    "context"
    "log"
    "os"
    "os/signal"

    "github.com/denial1337/crown-gobot/internal/bot"
    "github.com/denial1337/crown-gobot/internal/config"
)

func main() {
    // Загрузка конфигурации
    cfg := config.Load()
    
    if cfg.TelegramToken == "" {
        log.Fatal("TELEGRAM_TOKEN environment variable is required")
    }

    // Создание бота
    b, err := bot.New(cfg.TelegramToken)
    if err != nil {
        log.Fatalf("Failed to create bot: %v", err)
    }

    // Обработка сигналов для graceful shutdown
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    log.Println("Starting bot...")
    
    // Запуск бота
    b.Start(ctx)
    
    log.Println("Bot stopped")
}
