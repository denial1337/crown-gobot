package scheduler

import (
    "time"
)

type TaskHandler func(userID int64, taskID int64) error

type ScheduledTask struct {
    TaskID    int64
    UserID    int64
    Interval  string
    NextRun   time.Time
    Handler   TaskHandler
}
