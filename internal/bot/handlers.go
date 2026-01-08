package bot

import (
	"context"
	"log"

	tgBot "github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

// Обработчик для тестовой команды
func (b *Bot) testHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
    if update.Message == nil {
        return
    }

    // Если нужно конвертировать сообщение в задачу
    task := MessageToTask(update)
    if task != nil {
        if err := b.conn.InsertTask(task.String(), ctx); err != nil {
            log.Printf("Error while inserting %s: %v", task.String(), err)
        } else {
            bot.SendMessage(ctx, &tgBot.SendMessageParams{
                ChatID: update.Message.Chat.ID,
                Text:   "Task added successfully",
            })
        }
        return
    }
}

// Обработчик для команды /com с кнопками
func (b *Bot) comHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
    if update.Message == nil {
        return
    }

    chatID := update.Message.Chat.ID
    
    // Создаем inline-клавиатуру с двумя кнопками
    keyboard := [][]models.InlineKeyboardButton{
        {
            {
                Text:         "📝 Первая кнопка",
                CallbackData: "com_option_1",
            },
            {
                Text:         "🔢 Вторая кнопка", 
                CallbackData: "com_option_2",
            },
        },
    }
    
    // Отправляем сообщение с кнопками
    _, err := bot.SendMessage(ctx, &tgBot.SendMessageParams{
        ChatID: chatID,
        Text:   "Выберите вариант:",
        ReplyMarkup: &models.InlineKeyboardMarkup{
            InlineKeyboard: keyboard,
        },
        ParseMode: models.ParseModeMarkdown,
    })
    
    if err != nil {
        log.Printf("Ошибка отправки сообщения: %v\n", err)
    }
}

// Обработчик callback-запросов от inline-кнопок
func (b *Bot) callbackHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
    if update.CallbackQuery == nil {
        return
    }

    // Минимальная обработка - логируем "test"
    log.Println("test")

    // Ответим на callback, чтобы убрать "часики" у кнопки
    bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
    })
}