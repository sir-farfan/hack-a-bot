package events

import (
	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sir-farfan/hack-a-bot/model"
	"github.com/sir-farfan/hack-a-bot/service/sql_event"
)

func CMDName() string {
	return "events"
}

func RegisterCommands(processors map[string]model.Processor) error {
	processors["events"] = Events
	return nil
}

// Events - Pick an event from the list
func Events(recv tgapi.Update) (*tgapi.Chattable, error) {
	var choices []tgapi.InlineKeyboardButton

	events := GetEvents()
	for _, evt := range events {
		button := tgapi.InlineKeyboardButton{Text: evt.Name,
			CallbackData: &evt.Description,
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
func GetEvents() []model.Event {
	db := sql_event.New()

	return db.GetEvent("")
}
