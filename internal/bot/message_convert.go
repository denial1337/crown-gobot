package bot

import (
	"github.com/denial1337/crown-gobot/internal/models"
	tg "github.com/go-telegram/bot/models"
	"time"
)


func MessageToTask(upd *tg.Update) *models.Task {
	msg := upd.Message.Text
	if len(msg) <= 6 { return nil }
	taskText := msg[6:]
	task := models.Task{
		Username: upd.Message.From.Username,
		ChatID: upd.Message.Chat.ID,
		Message: taskText,
		Schedule: "test",
		LastRun: time.Now(),
		NextRun: time.Now(),
		}
	return &task
}

