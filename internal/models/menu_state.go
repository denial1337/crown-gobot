package models

import	(
	tgModels "github.com/go-telegram/bot/models"
	"strconv"
	"github.com/denial1337/crown-gobot/internal/consts"
)


type MenuState struct {
	photo string
	keyboard [][]tgModels.InlineKeyboardButton
}



func WithSpecialButtons(prev string, keyboard [][]tgModels.InlineKeyboardButton) [][]tgModels.InlineKeyboardButton {
	row := []tgModels.InlineKeyboardButton{
		{
			Text: consts.FREQ_EVERY,
			CallbackData: prev + consts.FREQ_EVERY + consts.VALUES_SEPARATOR,
		},
		{
			Text: consts.FREQ_INTERVAL,
			CallbackData: prev + consts.FREQ_INTERVAL + consts.VALUES_SEPARATOR,
		},
		{
			Text: consts.FREQ_MANY,
			CallbackData: prev + consts.FREQ_MANY + consts.VALUES_SEPARATOR,
		},
		{
			Text: consts.DELETE,
			CallbackData: prev + consts.DELETE + consts.VALUES_SEPARATOR,
		},
		{
			Text: "OK",
			CallbackData: prev + consts.TIMEFRAME_SEPARATOR + consts.VALUES_SEPARATOR,
		},
	}

	keyboard = append(keyboard, row)

	return keyboard
}

func GetMainMenuState() [][]tgModels.InlineKeyboardButton {
	keyboard := [][]tgModels.InlineKeyboardButton{
		{
			{
				Text: "Мои таски",
				CallbackData: "show_tasks",
			},
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

	return keyboard
}

// добавить декоратор который имеет / , - и del
// и мб окно ввода которое будет отображать текущее значение
// после нажатия на del 

func GetMinutesKeyboard(prev string) [][]tgModels.InlineKeyboardButton {
	row := []tgModels.InlineKeyboardButton{}
    keyboard := [][]tgModels.InlineKeyboardButton{}
    for i := range 12 {
        button := tgModels.InlineKeyboardButton{
            Text:         strconv.Itoa(i*5),
            CallbackData: prev + strconv.Itoa(i*5) + consts.VALUES_SEPARATOR,
        }
        row = append(row, button)
		if (i + 1) % 3 == 0 {
			keyboard = append(keyboard, row)
			row = []tgModels.InlineKeyboardButton{}
		}
    }
	keyboard = append(keyboard, row)

	return keyboard
}


func GetHoursKeyboard(prev string) [][]tgModels.InlineKeyboardButton {
	row := []tgModels.InlineKeyboardButton{}
    keyboard := [][]tgModels.InlineKeyboardButton{}
    for i := range 23 {
        button := tgModels.InlineKeyboardButton{
            Text:         strconv.Itoa(i),
            CallbackData: prev + strconv.Itoa(i) + consts.VALUES_SEPARATOR,
        }
        row = append(row, button)
    }
	keyboard = append(keyboard, row)

	return keyboard
}


func GetDOWKeyboard(prev string) [][]tgModels.InlineKeyboardButton {
	keyboard := [][]tgModels.InlineKeyboardButton{}
	row := []tgModels.InlineKeyboardButton{}
    for _, dow := range consts.ShortWeekdays {
        button := tgModels.InlineKeyboardButton{
            Text:         dow,
            CallbackData: prev + dow + consts.VALUES_SEPARATOR,
        }
        row = append(row, button)
    }
	keyboard = append(keyboard, row)

	return keyboard
}


func GetDOMKeyboard(prev string) [][]tgModels.InlineKeyboardButton {
	keyboard := [][]tgModels.InlineKeyboardButton{}
    var row []tgModels.InlineKeyboardButton
    for i := 1; i <= 31; i++ {
        button := tgModels.InlineKeyboardButton{
            Text:         strconv.Itoa(i),
            CallbackData: prev + strconv.Itoa(i) + consts.VALUES_SEPARATOR,
        }
        row = append(row, button)
        
        if i % 7 == 0 || i == 31 {
            keyboard = append(keyboard, row)
            row = []tgModels.InlineKeyboardButton{}
        }
    }

	return keyboard
}

func GetMonthKeyboard(prev string) [][]tgModels.InlineKeyboardButton {
	row := []tgModels.InlineKeyboardButton{}
    keyboard := [][]tgModels.InlineKeyboardButton{}
    for i := range 12 {
        button := tgModels.InlineKeyboardButton{
            Text:         strconv.Itoa(i + 1),
            CallbackData: prev + strconv.Itoa(i + 1) + consts.VALUES_SEPARATOR,
        }
        row = append(row, button)
		if (i + 1) % 3 == 0 {
			keyboard = append(keyboard, row)
			row = []tgModels.InlineKeyboardButton{}
		}
    }
	keyboard = append(keyboard, row)

	return keyboard
}
