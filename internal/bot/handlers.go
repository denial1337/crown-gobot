package bot

import (
    "context"
	"fmt"
    "log"
    tgBot "github.com/go-telegram/bot"
    "github.com/go-telegram/bot/models"
    // "os"
    // "bytes"
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
func (b *Bot) scheduleHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
    if update.Message == nil {
        return
    }

    chatID := update.Message.Chat.ID
    
    // Создаем inline-клавиатуру с двумя кнопками
    keyboard := [][]models.InlineKeyboardButton{
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
        ReplyMarkup: &models.InlineKeyboardMarkup{
            InlineKeyboard: keyboard,
        },
        ParseMode: models.ParseModeMarkdown,
    })
    
    if err != nil {
        log.Printf("Ошибка отправки сообщения: %v\n", err)
    }
}

var ShortWeekdays = [7]string{
	"Mon",
	"Tue",
	"Wed",
	"Thu",
	"Fri",
	"Sat",
	"Sun",
}


func (b *Bot) startHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
    keyboard := [][]models.InlineKeyboardButton{
        {
            {
                Text: "Создать таску",
                CallbackData: "create_task",
            },
            {
                Text: "Удалить таску",
                CallbackData: "delete_task",
            },
        },
    }
    //fileData, _ := os.ReadFile("./media/11.PNG")
    params := &tgBot.SendPhotoParams{
		ChatID:  update.Message.Chat.ID,
		//Photo:   &models.InputFileUpload{Filename: "media/11.PNG", Data: bytes.NewReader(fileData)},
        Photo: &models.InputFileString{Data: b.photo},
		Caption: "че делаем, пэпэ?",
        ReplyMarkup: &models.InlineKeyboardMarkup{InlineKeyboard: keyboard},
	}

    // Отправляем фото с inline клавиатурой
    sentPhoto, err := bot.SendPhoto(ctx, params)
    if err == nil {
        log.Println("FileID:", sentPhoto.Photo[0].FileID)
    }
}

func (b *Bot) createTaskHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
    keyboard := [][]models.InlineKeyboardButton{
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
        ReplyMarkup: &models.InlineKeyboardMarkup{
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

func (b *Bot) dailyCallbackHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
	log.Print("INSIDE dailyCallbackHandler")
	log.Print(update.CallbackQuery.ID)
    keyboard := [][]models.InlineKeyboardButton{}
    chatID := update.CallbackQuery.Message.Message.Chat.ID
	row := []models.InlineKeyboardButton{}
    for _, day := range ShortWeekdays {
        button := models.InlineKeyboardButton{
            Text:         day,
            CallbackData: fmt.Sprintf("schedule_weekday_%s", day),
        }
        row = append(row, button)
    }
	keyboard = append(keyboard, row)

   	bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
    })
    _, err := bot.SendMessage(ctx, &tgBot.SendMessageParams{
        ChatID: chatID,
        Text:   "Выберите день недели:",
        ReplyMarkup: &models.InlineKeyboardMarkup{
            InlineKeyboard: keyboard,
        },
    })
    
    if err != nil {
        log.Printf("Ошибка отправки сообщения: %v\n", err)
    }
}


func (b *Bot) weeklyCallbackHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
	log.Print("INSIDE weeklyCallbackHandler")
	log.Print(update.CallbackQuery.ID)
    keyboard := [][]models.InlineKeyboardButton{}
    chatID := update.CallbackQuery.Message.Message.Chat.ID
	row := []models.InlineKeyboardButton{}
    for _, day := range ShortWeekdays {
        button := models.InlineKeyboardButton{
            Text:         day,
            CallbackData: fmt.Sprintf("schedule_weekday_%s", day),
        }
        row = append(row, button)
    }
	keyboard = append(keyboard, row)

   	bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
    })
    _, err := bot.SendMessage(ctx, &tgBot.SendMessageParams{
        ChatID: chatID,
        Text:   "Выберите день недели:",
        ReplyMarkup: &models.InlineKeyboardMarkup{
            InlineKeyboard: keyboard,
        },
    })
    
    if err != nil {
        log.Printf("Ошибка отправки сообщения: %v\n", err)
    }
}

func (b *Bot) monthlyCallbackHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
	log.Print("INSIDE monthlyHandler")
	log.Print(update.CallbackQuery.ID)
    keyboard := [][]models.InlineKeyboardButton{}
    chatID := update.CallbackQuery.Message.Message.Chat.ID

    var row []models.InlineKeyboardButton
    for i := 1; i <= 31; i++ {
        button := models.InlineKeyboardButton{
            Text:         fmt.Sprintf("%d", i),
            CallbackData: fmt.Sprintf("schedule_day_%d", i),
        }
        row = append(row, button)
        
        if i % 7 == 0 || i == 31 {
            keyboard = append(keyboard, row)
            row = []models.InlineKeyboardButton{}
        }
    }
   	bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
    })
    _, err := bot.SendMessage(ctx, &tgBot.SendMessageParams{
        ChatID: chatID,
        Text:   "Выберите день месяца:",
        ReplyMarkup: &models.InlineKeyboardMarkup{
            InlineKeyboard: keyboard,
        },
    })
    
    if err != nil {
        log.Printf("Ошибка отправки сообщения: %v\n", err)
    }
}

func (b *Bot) scheduleCallbackHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
    if update.CallbackQuery == nil {
        return
    }
	log.Print("INSIDE scheduleCallbackHandler")
	log.Print(update.CallbackQuery.ID)
    callbackData := update.CallbackQuery.Data
	
	chatID := update.CallbackQuery.Message.Message.Chat.ID
    log.Println(callbackData, chatID)
	

    switch callbackData {
    case "schedule_monthly":
        log.Println("MONTHLY")
		b.monthlyCallbackHandler(ctx, bot, update)
	case "schedule_weekly":
		log.Println("WEEKLY")
		b.weeklyCallbackHandler(ctx, bot, update)
    // Добавьте другие case для обработки других callback'ов
    }
    
    // Ответим на callback, чтобы убрать "часики" у кнопки
    bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
        CallbackQueryID: update.CallbackQuery.ID,
    })
}

func (b *Bot) callbackHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
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

func (b *Bot) askWithForceReply(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
    chatID := update.Message.Chat.ID
    
    _, err := bot.SendMessage(ctx, &tgBot.SendMessageParams{
        ChatID: chatID,
        Text:   "Введите ваше имя:",
        ReplyMarkup: &models.ForceReply{
            ForceReply: true,
            InputFieldPlaceholder: "Имя пользователя",
            Selective: true, // Только для этого пользователя
        },
    })
    
    if err != nil {
        log.Printf("Ошибка: %v", err)
    }
}

func (b *Bot) messageHandler(ctx context.Context, bot *tgBot.Bot, update *models.Update) {
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

