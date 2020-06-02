package events

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CMDName() string {
	return "events"
}

// Events - Pick an event from the list
func Events(recv tgapi.Update) (*tgapi.Chattable, error) {
	var choices []tgapi.InlineKeyboardButton

	events, _ := GetEvents()
	for _, evt := range events {
		log.Println(evt[1])
		button := tgapi.InlineKeyboardButton{Text: evt[0],
			CallbackData: &evt[1],
		}
		choices = append(choices, button)
	}

	row := tgapi.NewInlineKeyboardRow(choices...)
	keyboard := tgapi.NewInlineKeyboardMarkup(row)

	msg := tgapi.NewMessage(recv.Message.Chat.ID, "Found these events")
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = keyboard

	var chat tgapi.Chattable
	chat = msg

	return &chat, nil
}

// GetEvents - name, description
func GetEvents() ([][2]string, error) {
	return [][2]string{
		{"Silks", "Aerial silks"},
		{"straps", "Aerial straps"},
	}, nil
}
