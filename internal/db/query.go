package db

import (
	"fmt"
	"time"
)

func CreateTable() string {
	return `CREATE TABLE IF NOT EXISTS test_db (
            telegram_id BIGINT PRIMARY KEY,
            username TEXT,
            first_name TEXT)`
}


func CreateTaskTableQuery() string {
	return `CREATE TABLE IF NOT EXISTS tasks (
            id BIGSERIAL PRIMARY KEY,
            username TEXT,
            chat_id BIGINT,
			message TEXT,
			schedule TEXT,
			last_run TIMESTAMP,
			next_run TIMESTAMP)`
}


func InsertTaskQuery(values string) string {
	return fmt.Sprintf("INSERT INTO tasks (username, chat_id, message,schedule,last_run, next_run) VALUES (%s)", values)
}


func GetCurrentTasksQuery(t time.Time) string {
	//formattedTime := t.Format("2006-01-02T15:04:05Z07:00")
	//_ := formattedTime
	return "SELECT chat_id, message FROM tasks"
	//return fmt.Sprintf("SELECT chat_id, message FROM tasks WHERE next_run = '%s'", formattedTime)
}