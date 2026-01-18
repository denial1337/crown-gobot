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
    //"strconv"
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


func (b *Bot) startHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
    keyboard := models.GetMainMenuState("")
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


// func (b *Bot) callbackHandler(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
//     if update.CallbackQuery == nil {
//         return
//     }

//     callbackData := update.CallbackQuery.Data
//     log.Println(callbackData)

//     // Ответим на callback, чтобы убрать "часики" у кнопки
//     bot.AnswerCallbackQuery(ctx, &tgBot.AnswerCallbackQueryParams{
//         CallbackQueryID: update.CallbackQuery.ID,
//     })
// }

// func (b *Bot) askWithForceReply(ctx context.Context, bot *tgBot.Bot, update *tgModels.Update) {
//     chatID := update.Message.Chat.ID
    
//     _, err := bot.SendMessage(ctx, &tgBot.SendMessageParams{
//         ChatID: chatID,
//         Text:   "Введите ваше имя:",
//         ReplyMarkup: &tgModels.ForceReply{
//             ForceReply: true,
//             InputFieldPlaceholder: "Имя пользователя",
//             Selective: true, // Только для этого пользователя
//         },
//     })
    
//     if err != nil {
//         log.Printf("Ошибка: %v", err)
//     }
// }

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
    var currentKeyboard func(string)[][]tgModels.InlineKeyboardButton
    var nextKeyboard func(string)[][]tgModels.InlineKeyboardButton
    var keyboard [][]tgModels.InlineKeyboardButton
    var nextTimeFrame string
    //text := update.CallbackQuery.Message.Message.Caption
    
    callback := update.CallbackQuery.Data
    
    switch true {
    case callback == "create_task":
        keyboard = models.WithSpecialButtons("", models.GetMinutesKeyboard(consts.TIMEFRAME_MINUTES))
    case strings.Contains(callback, consts.TIMEFRAME_DOW):
        currentKeyboard = models.GetDOWKeyboard
        nextKeyboard = models.GetMainMenuState
    case strings.Contains(callback, consts.TIMEFRAME_MONTHES):
        currentKeyboard = models.GetMonthKeyboard
        nextKeyboard = models.GetDOWKeyboard
        nextTimeFrame = consts.TIMEFRAME_DOW
    case strings.Contains(callback, consts.TIMEFRAME_DOM):
        currentKeyboard = models.GetDOMKeyboard
        nextKeyboard = models.GetMonthKeyboard
        nextTimeFrame = consts.TIMEFRAME_MONTHES
    case strings.Contains(callback, consts.TIMEFRAME_HOURS):
        currentKeyboard = models.GetHoursKeyboard
        nextKeyboard = models.GetDOMKeyboard
        nextTimeFrame = consts.TIMEFRAME_DOM
    case strings.Contains(callback, consts.TIMEFRAME_MINUTES):
        currentKeyboard = models.GetMinutesKeyboard
        nextKeyboard = models.GetHoursKeyboard
        nextTimeFrame = consts.TIMEFRAME_HOURS
    }

    if keyboard == nil {
        keyboard = getNextKeyboard(&callback, nextTimeFrame, currentKeyboard, nextKeyboard)
    }
    log.Print("INSIDE dailyCallbackHandler")
	log.Print(update.CallbackQuery.Data)
    
    _, err := bot.EditMessageCaption(ctx, &tgBot.EditMessageCaptionParams{
        ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
        MessageID:   update.CallbackQuery.Message.Message.ID,
        Caption:     callback,
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


func getNextKeyboard(callBackData *string, 
    nextTimeFrame string,
    currentKeyboard func(string)[][]tgModels.InlineKeyboardButton,
    nextKeyboard func(string)[][]tgModels.InlineKeyboardButton,
    ) [][]tgModels.InlineKeyboardButton {
    parts := strings.Split(*callBackData, consts.VALUES_SEPARATOR)
    for _, v := range parts {
        log.Println(v)
    }
    var lastPart string
    if len(parts) > 1 {
        lastPart = parts[len(parts) - 2]
    } else { return models.WithSpecialButtons(*callBackData, currentKeyboard(*callBackData))}
    switch lastPart {
    case consts.TIMEFRAME_SEPARATOR:
        if nextTimeFrame != "" {
            *callBackData = *callBackData + nextTimeFrame
        }
        return models.WithSpecialButtons(*callBackData, nextKeyboard(*callBackData))
    case consts.DELETE:
        *callBackData = deleteLastPart(*callBackData)
        return models.WithSpecialButtons(*callBackData, currentKeyboard(*callBackData))
    default:
        return models.WithSpecialButtons(*callBackData, currentKeyboard(*callBackData))
    }
    
}

func deleteLastPart(s string) string {
    parts := strings.Split(s, consts.VALUES_SEPARATOR)
    parts = parts[:len(parts) - 3]
    
    return strings.Join(parts, consts.VALUES_SEPARATOR) + "|"
}

// func containsWeekday(s string) bool {
//     for _, day := range consts.ShortWeekdays {
//         if day == s {
//             return true
//         } 
//     }
//     return false
// }