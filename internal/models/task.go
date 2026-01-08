package models

import (
    "time"
    "fmt"
)

type Task struct {
	Username string
	ChatID int64
	Message string
	Schedule string
	LastRun time.Time
	NextRun time.Time
}


type TaskMessage struct {
	ChatID int64
	Message string
}

func (t *Task) String() string {
    return fmt.Sprintf("'%s', %d, '%s', '%s', '%s', '%s'",
        t.Username,
        t.ChatID,
        t.Message,
        t.Schedule,
        t.LastRun.Format("2006-01-02T15:04:05Z07:00"),
        t.NextRun.Format("2006-01-02T15:04:05Z07:00"))
}