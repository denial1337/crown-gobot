
package bot

import (
    "context"
    "log"

    "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
)

type Bot struct {
    tgBot *bot.Bot
}

func New(token string) (*Bot, error) {
    opts := []bot.Option{
        bot.WithDefaultHandler(defaultHandler),
    }

    tgBot, err := bot.New(token, opts...)
    if err != nil {
        return nil, err
    }

    return &Bot{
        tgBot: tgBot,
    }, nil
}

func (b *Bot) Start(ctx context.Context) {
    b.registerHandlers()
    b.tgBot.Start(ctx)
}

func (b *Bot) registerHandlers() {
    b.tgBot.RegisterHandler(bot.HandlerTypeMessageText, "/start", bot.MatchTypeExact, startHandler)
}

func startHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
    _, err := b.SendMessage(ctx, &bot.SendMessageParams{
        ChatID: update.Message.Chat.ID,
        Text:   "Hello",
    })
    if err != nil {
        log.Printf("Failed to send message: %v", err)
    }
}

func defaultHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
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
