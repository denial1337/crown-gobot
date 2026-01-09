package bot

import (
    "context"
    "log"
    //"fmt"
    tgbotapi "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
    "github.com/denial1337/crown-gobot/internal/db"
    m "github.com/denial1337/crown-gobot/internal/models"
    // "os"
    // "bytes"
)

type Bot struct {
    tgBot *tgbotapi.Bot
    conn  *db.DbConnection
    stor *m.MemoryStorage
    photo string
}

func New(token string, conn *db.DbConnection, stor *m.MemoryStorage) (*Bot, error) {
    opts := []tgbotapi.Option{
        tgbotapi.WithDefaultHandler(defaultHandler),
    }

    tgBot, err := tgbotapi.New(token, opts...)
    if err != nil {
        return nil, err
    }

    return &Bot{
        tgBot: tgBot,
        conn:  conn,
        stor: stor,
        photo: "AgACAgIAAxkDAAN6aWDoD_bJ1xWRf7drohnkhwpknsMAAjYOaxs9DwlLhXhLEDXgTJ0BAAMCAANzAAM4BA",
    }, nil
}

func (b *Bot) Start(ctx context.Context) {
    b.registerHandlers()
    b.tgBot.Start(ctx)
}

func (b *Bot) registerHandlers() {
    b.tgBot.RegisterHandler(tgbotapi.HandlerTypeMessageText, "/test", tgbotapi.MatchTypePrefix, b.testHandler)
    b.tgBot.RegisterHandler(tgbotapi.HandlerTypeMessageText, "/com", tgbotapi.MatchTypePrefix, b.scheduleHandler)
    b.tgBot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData, "com_option_", tgbotapi.MatchTypePrefix, b.callbackHandler)
    b.tgBot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData, "schedule_", tgbotapi.MatchTypePrefix, b.scheduleCallbackHandler)
     b.tgBot.RegisterHandler(tgbotapi.HandlerTypeCallbackQueryData, "create_task", tgbotapi.MatchTypePrefix, b.createTaskHandler)
    b.tgBot.RegisterHandler(tgbotapi.HandlerTypeMessageText, "/task", tgbotapi.MatchTypePrefix, taskHandler)
    b.tgBot.RegisterHandler(tgbotapi.HandlerTypeMessageText, "/start", tgbotapi.MatchTypePrefix, b.startHandler)
    b.tgBot.RegisterHandler(tgbotapi.HandlerTypeMessageText, "", tgbotapi.MatchTypePrefix, b.messageHandler)
}

func (b *Bot) SendMessage(ctx context.Context, chatID int64, message string) error {
    _, err := b.tgBot.SendMessage(ctx, &tgbotapi.SendMessageParams{
        ChatID: chatID,
        Text:   message,
    })
    if err != nil {
        log.Printf("Failed to send message: %v", err)
        return err
    }
    return nil
}

func taskHandler(ctx context.Context, bot *tgbotapi.Bot, update *models.Update) {
    task := MessageToTask(update)
    if task != nil {
        /*if err := b.conn.InsertTask(task.String()); err != nil {
            log.Print("Error while inserting", task.String())
        } else {
            b.SendMessage(ctx, &bot.SendMessageParams {
                ChatID: update.Message.Chat.ID,
                Text: "task added successfully",
            })
        }*/
        _, err := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text:   task.String(),
        })
        if err != nil {
            log.Print("Error in taskHandler", err)
        }
    }
}

// func (b *Bot) startHandler(ctx context.Context, bot *tgbotapi.Bot, update *models.Update) {
//     fileData, errReadFile := os.ReadFile("./media/11.PNG")
// 	if errReadFile != nil {
// 		log.Printf("error read file, %v\n", errReadFile)
// 		return
// 	}

// 	params := &tgbotapi.SendPhotoParams{
// 		ChatID:  update.Message.Chat.ID,
// 		Photo:   &models.InputFileUpload{Filename: "media/11.PNG", Data: bytes.NewReader(fileData)},
// 		Caption: "New uploaded Facebook logo",
// 	}

	
//     sentPhoto, err := bot.SendPhoto(ctx, params)
//     if err == nil {
//         // sentPhoto.Photo[0].FileID - это ваш file_id для будущего использования
//         log.Println("FileID:", sentPhoto.Photo[0].FileID)
//     }
// }
func defaultHandler(ctx context.Context, bot *tgbotapi.Bot, update *models.Update) {
    if update.Message != nil {
        _, err := bot.SendMessage(ctx, &tgbotapi.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text:   "Use /start command",
        })
        if err != nil {
            log.Printf("Failed to send message: %v", err)
        }
    }
}