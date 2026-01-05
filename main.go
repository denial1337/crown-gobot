package main

import (
    "log"
    
    "cron-bot/scheduler"
    "cron-bot/storage"
)

func main() {
    // Инициализируем хранилище
    store, err := storage.NewStorage("./tasks.db")
    if err != nil {
        log.Fatal("Failed to initialize storage:", err)
    }
    defer store.Close()

    // Создаем шедулер с обработчиком задач
    sch := scheduler.NewScheduler(store, taskHandler)
    
    // Запускаем шедулер
    if err := sch.Start(); err != nil {
        log.Fatal("Failed to start scheduler:", err)
    }
    defer sch.Stop()

    // Здесь позже будет инициализация Telegram бота
    log.Println("Cron service started. Waiting for Telegram bot integration...")
    
    // Блокируем главную горутину
    select {}
}

// Пример обработчика задач
func taskHandler(userID int64, taskID int64) error {
    log.Printf("Executing task %d for user %d", taskID, userID)
    // Здесь будет логика выполнения задачи
    // Например, отправка сообщения пользователю
    return nil
}
