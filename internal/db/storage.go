package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/denial1337/crown-gobot/internal/models"
	"context"
	"time"
	"errors"
	"log"
)

const connString = "postgres://gouser:gopassword@localhost:5432/crown_gobot"



type DbConnection struct {
	conn *pgx.Conn
}

func InitDb(ctx context.Context) (*DbConnection, error) {
    conn, err := pgx.Connect(ctx, connString)
    if err != nil {
      	return nil, err
	}

	_, err = conn.Exec(ctx, CreateTaskTableQuery())
	if err != nil {
		return nil, err
	}
	return &DbConnection{conn: conn}, nil
}

func (conn *DbConnection) Close(ctx context.Context) error {
	if conn != nil {
		return conn.conn.Close(ctx)
	}
	return nil
}

func (conn *DbConnection) InsertTask(values string, ctx context.Context) error {
	if conn != nil {
		_, err := conn.conn.Exec(ctx, InsertTaskQuery(values))
		if err != nil { return err}
	}
	return nil
}

func (conn *DbConnection) GetCurrentTasks(ctx context.Context, t time.Time) ([]models.TaskMessage, error) {
	log.Print("we are in GetCurrentTasks")
	if conn == nil {
		return nil,  errors.New("Отсутствует подключение к базе данных")
	}
	
	query := GetCurrentTasksQuery(t)
	log.Print("we get query", query)
	rows, err := conn.conn.Query(ctx, query)
	log.Print("we get rows", rows != nil)
	if err != nil {
		return nil, err
	}

	var tasks []models.TaskMessage

	for rows.Next() {
		log.Print("we iter")
		var task models.TaskMessage
		if err := rows.Scan(&task.ChatID, &task.Message); err != nil {
			return nil, errors.New("Ошибка при попытке чтения данных")
		}
		tasks = append(tasks, task)
	}
	log.Print("we return")
	return tasks, nil
}


