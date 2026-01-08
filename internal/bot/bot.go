package bot

import (
    "context"
    "log"

    "github.com/go-telegram/bot"
    tg "github.com/go-telegram/bot/models"
    "github.com/denial1337/crown-gobot/internal/db"
    //"github.com/denial1337/crown-gobot/internal/models"
)

type Bot struct {
    tgBot *bot.Bot
    handler *Handler
}

type Handler struct
{
    conn *db.DbConnection
}

func (h *Handler) testHandler(ctx context.Context, b *bot.Bot, update *tg.Update) {
    task := MessageToTask(update)
    if task == nil { return }

    if err := h.conn.InsertTask(task.String(), ctx); err != nil {
        log.Print("Error while iserting", task.String(), err)
    } else {
        b.SendMessage(ctx, &bot.SendMessageParams {
            ChatID: update.Message.Chat.ID,
            Text: "task added succesfuly",
        })
    }
}

func New(token string, conn *db.DbConnection) (*Bot, error) {
    opts := []bot.Option{
        bot.WithDefaultHandler(defaultHandler),
    }

    tgBot, err := bot.New(token, opts...)
    if err != nil {
        return nil, err
    }

    handler := &Handler{
        conn: conn,
    }

    return &Bot{
        tgBot: tgBot,
        handler: handler,
    }, nil
}

func (b *Bot) Start(ctx context.Context) {
    b.registerHandlers()
    b.tgBot.Start(ctx)
}

func (b *Bot) registerHandlers() {
    b.tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/test", bot.MatchTypePrefix, b.handler.testHandler)
    b.tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
    b.tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/task", bot.MatchTypePrefix, taskHandler)
}

func (b *Bot) SendMessage(ctx context.Context, chat_id int64, message string) error {
    _, err := b.tgBot.SendMessage(ctx, &bot.SendMessageParams{
        ChatID: chat_id,
        Text: message,
    })
    if err != nil {
        log.Printf("Failed to send message: %v", err)
        return err
    }
    return nil
}


func startHandler(ctx context.Context, b *bot.Bot, update *tg.Update) {
    _, err := b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID: update.Message.Chat.ID,
        Text:   "Hello",
    })
    if err != nil {
        log.Printf("Failed to send message: %v", err)
    } else {
        log.Println("sayd Hello to", update.Message.From.Username)
    }
}

func taskHandler(ctx context.Context, b *bot.Bot, update *tg.Update) {
    task := MessageToTask(update)
    if task != nil {
        /*if err := b.conn.InsertTask(task.String()); err != nil {
            log.Print("Error while iserting", task.String())
        } else {
            b.SendMessage(ctx, &bot.SendMessageParams {
                ChatID: update.Message.Chat.ID,
                Text: "task added succesfuly",
            })
        }*/
        _, err := b.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text: task.String(),
        })
        if err != nil {
        log.Print("Error in taskHandler", err)
        }
    }
}


func defaultHandler(ctx context.Context, b *bot.Bot, update *tg.Update) {
    // Обработка сообщений без команд
    if update.Message != nil {
        _, err := b.SendMessage(ctx, &bot.SendMessageParams{
            ChatID: update.Message.Chat.ID,
            Text:   "Use /start command",
        })
        if err != nil {
            log.Printf("Failed to send message: %v", err)
        }
    }
}
