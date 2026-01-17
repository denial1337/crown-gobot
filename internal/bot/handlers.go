package bot

import (
    "context"
	//"fmt"
    "log"
    tgBot "github.com/go-telegram/bot"
    tgModels "github.com/go-telegram/bot/models"
    "github.com/denial1337/crown-gobot/internal/models"
   	"github.com/denial1337/crown-gobot/internal/consts"
    "strings"
    "strconv"
    // "os"
    // "bytes"
)

// Обработчик для тестовой команды
func (b *Bot) testHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
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
func (b *Bot) scheduleHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
    if update.Message == nil {
        return
    }

    chatID := update.Message.Chat.ID
    
    // Создаем inline-клавиатуру с двумя кнопками
    keyboard := [][]tgModels.InlineKeyboardButton{
        {
			{
                Text:         "Hourly",
                CallbackData: "schedule_daily",
            },
            {
                Text:         "Daily",
                CallbackData: "schedule_daily",
            },
            {
                Text:         "Weekly", 
                CallbackData: "schedule_weekly",
            },
            {
                Text:         "Monthly",
                CallbackData: "schedule_monthly",
            },
            {
                Text:         "Exact days", 
                CallbackData: "schedule_exact_days",
            },
        },
    }
    
    // Отправляем сообщение с кнопками
    _, err := bot.SendMessage(ctx, &tgBot.SendMessageParams{
        ChatID: chatID,
        Text:   "Выберите вариант:",
        ReplyMarkup: &tgModels.InlineKeyboardMarkup{
            InlineKeyboard: keyboard,
        },
        ParseMode: tgModels.ParseModeMarkdown,
    })
    
    if err != nil {
        log.Printf("Ошибка отправки сообщения: %v\n", err)
    }
}


func (b *Bot) startHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
    keyboard := models.GetMainMenuState()
    //fileData, _ := os.ReadFile("./media/11.PNG")
    params := &tgBot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		//Photo:   &tgModels.InputFileUpload{Filename: "media/11.PNG", Data: bytes.NewReader(fileData)},
        Photo: &tgModels.InputFileString{Data: b.photo},
		Caption: "че делаем, пэпэ?",
        ReplyMarkup: &tgModels.InlineKeyboardMarkup{InlineKeyboard: keyboard},
	}

    // Отправляем фото с inline клавиатурой
    sentPhoto, err := bot.SendPhoto(ctx, params)
    if err == nil {
        log.Println("FileID:", sentPhoto.Photo[0].FileID)
    }
}

func (b *Bot) createTaskHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
    keyboard := [][]tgModels.InlineKeyboardButton{
        {
            {
                Text: "1b",
                CallbackData: "bb1",
            },
            {
                Text: "2b",
                CallbackData: "bb2",
            },
        },
    }
    _, err := bot.EditMessageReplyMarkup(ctx, &tgBot.EditMessageReplyMarkupParams{
        ChatID:    update.CallbackQuery.Message.Message.Chat.ID,
        MessageID: update.CallbackQuery.Message.Message.ID,
        ReplyMarkup: &tgModels.InlineKeyboardMarkup{
            InlineKeyboard: keyboard,
        },
    })
    
    if err != nil {
        log.Printf("Ошибка изменения клавиатуры: %v", err)
    }
    
    // Обязательно отвечаем на callback query (убираем "часики")
    bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
    })
}


func (b *Bot) callbackHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
    if update.CallbackQuery == nil {
        return
    }

    callbackData := update.CallbackQuery.Data
    log.Println(callbackData)

    // Ответим на callback, чтобы убрать "часики" у кнопки
    bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
    })
}

func (b *Bot) askWithForceReply(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
    chatID := update.Message.Chat.ID
    
    _, err := bot.SendMessage(ctx, &tgBot.SendMessageParams{
        ChatID: chatID,
        Text:   "Введите ваше имя:",
        ReplyMarkup: &tgModels.ForceReply{
            ForceReply: true,
            InputFieldPlaceholder: "Имя пользователя",
            Selective: true, // Только для этого пользователя
        },
    })
    
    if err != nil {
        log.Printf("Ошибка: %v", err)
    }
}

func (b *Bot) messageHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
    if update.Message == nil {
        return
    }
    
    // Проверяем, является ли сообщение ответом на ForceReply
    if update.Message.ReplyToMessage != nil {
        // Проверяем, было ли это сообщение с ForceReply
        if update.Message.ReplyToMessage.Text == "Введите ваше имя:" {
            // Обрабатываем ввод имени
            name := update.Message.Text
            bot.SendMessage(ctx, &tgBot.SendMessageParams{
                ChatID: update.Message.Chat.ID,
                Text:   "Привет, " + name + "!",
            })
        }
    }
}

func (b *Bot) callBackHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
    var keyboard [][]tgModels.InlineKeyboardButton
    text := update.CallbackQuery.Message.Message.Caption
    
    callback := update.CallbackQuery.Data
    
    switch true {
    case callback == "create_task":
        keyboard = models.WithSpecialButtons("", models.GetMinutesKeyboard(consts.TIMEFRAME_MINUTES))
        text = callback + "\nМинуты"
    case strings.Contains(callback, consts.TIMEFRAME_MINUTES):
        keyboard = models.WithSpecialButtons(callback, models.GetHoursKeyboard(callback))
        text = callback + "\nЧасы"
    case lastPart == consts.DELETE:
        callback = deleteLastPart(callback)
        text = callback
     case func() bool {
        if _, err := strconv.Atoi(lastPart); err == nil { return true }
        return false
    }():
        callback = callback + consts.VALUES_SEPARATOR
        text = callback
    case containsWeekday(lastPart):
        callback = callback + consts.VALUES_SEPARATOR
        text = callback
    }
    
    log.Print("INSIDE dailyCallbackHandler")
	log.Print(update.CallbackQuery.Data)
    
    _, err := bot.EditMessageCaption(ctx, &tgBot.EditMessageCaptionParams{
        ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
        MessageID:   update.CallbackQuery.Message.Message.ID,
        Caption:     text,
        ReplyMarkup: &tgModels.InlineKeyboardMarkup{
            InlineKeyboard: keyboard,
        },
    })
    if err != nil {
        log.Printf("Error editing message: %v", err)
    }    
    
    bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
    })
}


func lastPartHandle(callBackData string) string {
    parts := strings.Split(callback, consts.VALUES_SEPARATOR)
    var lastPart string
    if len(parts) > 0 {
        lastPart = parts[len(parts) - 1]
    } else { return callBackData }

    switch lastPart {
    case consts.TIMEFRAME_SEPARATOR:
        
    }
}

func deleteLastPart(s string) string {
    parts := strings.Split(s, consts.TIMEFRAME_SEPARATOR)
    parts = parts[:len(parts) - 2]
    
    return strings.Join(parts, consts.TIMEFRAME_SEPARATOR)
}

func containsWeekday(s string) bool {
    for _, day := range consts.ShortWeekdays {
        if day == s {
            return true
        }
    }
    return false
}