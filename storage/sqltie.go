package storage

import (
    "database/sql"
    "fmt"
    "log"
    "time"

    _ "github.com/mattn/go-sqlite3"
)

type Task struct {
    TaskID    int64
    UserID    int64
    CreatedAt time.Time
    NextRun   time.Time
    Interval  string // cron expression like "* * * * *"
}

type Storage struct {
    db *sql.DB
}

func NewStorage(dbPath string) (*Storage, error) {
    db, err := sql.Open("sqlite3", dbPath)
    if err != nil {
        return nil, fmt.Errorf("failed to open database: %w", err)
    }

    if err := db.Ping(); err != nil {
        return nil, fmt.Errorf("failed to ping database: %w", err)
    }

    storage := &Storage{db: db}
    if err := storage.createTables(); err != nil {
        return nil, fmt.Errorf("failed to create tables: %w", err)
    }

    return storage, nil
}

func (s *Storage) createTables() error {
    query := `
    CREATE TABLE IF NOT EXISTS tasks (
        task_id INTEGER PRIMARY KEY AUTOINCREMENT,
        user_id INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
        next_run TIMESTAMP NOT NULL,
        interval TEXT NOT NULL,
        UNIQUE(user_id, task_id)
    );

    CREATE INDEX IF NOT EXISTS idx_next_run ON tasks(next_run);
    CREATE INDEX IF NOT EXISTS idx_user_id ON tasks(user_id);
    `

    _, err := s.db.Exec(query)
    return err
}

func (s *Storage) AddTask(userID int64, interval string, nextRun time.Time) (int64, error) {
    query := `INSERT INTO tasks (user_id, interval, next_run) VALUES (?, ?, ?)`
    
    result, err := s.db.Exec(query, userID, interval, nextRun)
    if err != nil {
        return 0, fmt.Errorf("failed to insert task: %w", err)
    }
    
    taskID, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("failed to get last insert ID: %w", err)
    }
    
    return taskID, nil
}

func (s *Storage) DeleteTask(userID, taskID int64) error {
    query := `DELETE FROM tasks WHERE user_id = ? AND task_id = ?`
    
    result, err := s.db.Exec(query, userID, taskID)
    if err != nil {
        return fmt.Errorf("failed to delete task: %w", err)
    }
    
    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("failed to get rows affected: %w", err)
    }
    
    if rowsAffected == 0 {
        return fmt.Errorf("task not found")
    }
    
    return nil
}

func (s *Storage) GetUserTasks(userID int64) ([]Task, error) {
    query := `SELECT task_id, user_id, created_at, next_run, interval 
              FROM tasks WHERE user_id = ? ORDER BY next_run`
    
    rows, err := s.db.Query(query, userID)
    if err != nil {
        return nil, fmt.Errorf("failed to query tasks: %w", err)
    }
    defer rows.Close()
    
    var tasks []Task
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.TaskID, &task.UserID, &task.CreatedAt, &task.NextRun, &task.Interval)
        if err != nil {
            return nil, fmt.Errorf("failed to scan task: %w", err)
        }
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}

func (s *Storage) GetDueTasks() ([]Task, error) {
    query := `SELECT task_id, user_id, created_at, next_run, interval 
              FROM tasks WHERE next_run <= ? ORDER BY next_run`
    
    rows, err := s.db.Query(query, time.Now())
    if err != nil {
        return nil, fmt.Errorf("failed to query due tasks: %w", err)
    }
    defer rows.Close()
    
    var tasks []Task
    for rows.Next() {
        var task Task
        err := rows.Scan(&task.TaskID, &task.UserID, &task.CreatedAt, &task.NextRun, &task.Interval)
        if err != nil {
            return nil, fmt.Errorf("failed to scan task: %w", err)
        }
        tasks = append(tasks, task)
    }
    
    return tasks, nil
}

func (s *Storage) UpdateNextRun(taskID int64, nextRun time.Time) error {
    query := `UPDATE tasks SET next_run = ? WHERE task_id = ?`
    
    _, err := s.db.Exec(query, nextRun, taskID)
    if err != nil {
        return fmt.Errorf("failed to update next run: %w", err)
    }
    
    return nil
}

func (s *Storage) Close() error {
    return s.db.Close()
}
