package scheduler

import (
    "fmt"
    "log"
    "time"

    "cron-bot/storage"

    "github.com/robfig/cron/v3"
)

type Scheduler struct {
    storage    *storage.Storage
    cronParser cron.Parser
    tasks      map[int64]*ScheduledTask // taskID -> task
    cron       *cron.Cron
    handler    TaskHandler
}

func NewScheduler(storage *storage.Storage, handler TaskHandler) *Scheduler {
    parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
    
    return &Scheduler{
        storage:    storage,
        cronParser: parser,
        tasks:      make(map[int64]*ScheduledTask),
        cron:       cron.New(cron.WithParser(parser)),
        handler:    handler,
    }
}

func (s *Scheduler) Start() error {
    // Загружаем существующие задачи из базы
    tasks, err := s.storage.GetDueTasks()
    if err != nil {
        return fmt.Errorf("failed to load tasks: %w", err)
    }

    // Запускаем крон планировщик
    s.cron.Start()
    
    // Запускаем шедулер для проверки задач
    go s.runScheduler()
    
    log.Println("Scheduler started")
    return nil
}

func (s *Scheduler) runScheduler() {
    ticker := time.NewTicker(1 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        s.checkDueTasks()
    }
}

func (s *Scheduler) checkDueTasks() {
    tasks, err := s.storage.GetDueTasks()
    if err != nil {
        log.Printf("Error getting due tasks: %v", err)
        return
    }

    for _, task := range tasks {
        // Выполняем задачу
        go s.executeTask(task)
        
        // Пересчитываем следующее время выполнения
        s.rescheduleTask(task)
    }
}

func (s *Scheduler) executeTask(task storage.Task) {
    if s.handler != nil {
        if err := s.handler(task.UserID, task.TaskID); err != nil {
            log.Printf("Error executing task %d for user %d: %v", task.TaskID, task.UserID, err)
        }
    }
}

func (s *Scheduler) rescheduleTask(task storage.Task) {
    schedule, err := s.cronParser.Parse(task.Interval)
    if err != nil {
        log.Printf("Error parsing cron expression for task %d: %v", task.TaskID, err)
        return
    }

    nextRun := schedule.Next(time.Now())
    if err := s.storage.UpdateNextRun(task.TaskID, nextRun); err != nil {
        log.Printf("Error updating next run for task %d: %v", task.TaskID, err)
    }
}

func (s *Scheduler) AddTask(userID int64, interval string) (int64, error) {
    // Парсим cron выражение для валидации
    _, err := s.cronParser.Parse(interval)
    if err != nil {
        return 0, fmt.Errorf("invalid cron expression: %w", err)
    }

    // Рассчитываем время следующего запуска
    schedule, _ := s.cronParser.Parse(interval)
    nextRun := schedule.Next(time.Now())

    // Сохраняем в базу
    taskID, err := s.storage.AddTask(userID, interval, nextRun)
    if err != nil {
        return 0, fmt.Errorf("failed to save task: %w", err)
    }

    return taskID, nil
}

func (s *Scheduler) Stop() {
    s.cron.Stop()
    log.Println("Scheduler stopped")
}
