package main

import (
    "context"
    "log"
    "os"
    "os/signal"
    "time"

    "github.com/denial1337/crown-gobot/internal/bot"
    "github.com/denial1337/crown-gobot/internal/db"
    "github.com/denial1337/crown-gobot/internal/config"
    "github.com/denial1337/crown-gobot/internal/models"
)

// ChatJoinRequest https://core.telegram.org/bots/api#chatjoinrequest

func main() {
    // Загрузка конфигурации
    cfg := config.Load()
    
    if cfg.TelegramToken == "" {
        log.Fatal("TELEGRAM_TOKEN environment variable is required")
    }

    // Обработка сигналов для graceful shutdown
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
    defer cancel()

    conn, err := db.InitDb(ctx)
    if err != nil {
        log.Print("Error while creating db", err)
    }
    defer conn.Close(ctx)

    ms := models.NewMemoryStorage()

    // Создание бота
    b, err := bot.New(cfg.TelegramToken, conn, ms)
    if err != nil {
        log.Fatalf("Failed to create bot: %v", err)
    }

    log.Println("Starting bot...")
    
    // Запуск бота
    go b.Start(ctx)
    log.Println("bot started ...")
    for {
            select {
            case <-ctx.Done():
                // Получен сигнал завершения
                log.Println("Received shutdown signal, stopping...")
                time.Sleep(1 * time.Second)
                log.Println("Bot stopped")
                return // Выход из функции main
                
            default:
                log.Println("In def...")
                // Получаем задачи (это быстро, т.к. это просто срез)
                tasks, err := conn.GetCurrentTasks(ctx, time.Now())
                log.Println(len(tasks))
                if err != nil {
                    // Не используем Fatalf, чтобы не падать при ошибке
                    // Просто логируем и продолжаем
                    log.Printf("Failed while getting tasks: %v", err)
                } else {
                    // Обрабатываем задачи
                    for _, task := range tasks {
                        // Отправляем сообщение (раскомментируй когда нужно)
                        //b.SendMessage(ctx, task.ChatID, task.Message)
                        log.Print(task.ChatID)
                        // Логируем для отладки
                        //log.Printf("Task: ChatID=%d, Message=%s", 
                            //task.ChatID, task.Message)
                    }
                }
                
                // Ждем 20 секунд ИЛИ сигнал завершения
                select {
                case <-time.After(20 * time.Second):
                    // Просто продолжаем цикл
                case <-ctx.Done():
                    // Выходим если получили сигнал завершения
                    return
                }
            }
        }
            
    log.Println("Bot stopped")
}
