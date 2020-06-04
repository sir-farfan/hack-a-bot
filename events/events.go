package events

import (
	"log"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sir-farfan/hack-a-bot/model"
	"github.com/sir-farfan/hack-a-bot/service/sql_event"
)

func CMDName() string {
	return "events"
}

func RegisterCommands(processors map[string]model.Processor) error {
	processors["events"] = Events
	processors["eventcreate"] = Create
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

	return db.GetEvent(0)
}

// Create - will only put the user in "event creation mode"
func Create(recv tgapi.Update) (*tgapi.Chattable, error) {
	log.Printf("DEBUG: [%s] %s\n", recv.Message.From.UserName, recv.Message.Text)

	db := sql_event.New()
	cookie := db.UserCookieGet(recv.Message.Chat.ID)
	if cookie.Cookie != "eventcreate" {
		db.UserCookieCreate(model.User{ID: recv.Message.Chat.ID, Cookie: "eventcreate"})
	}

	events := db.GetEvent(recv.Message.Chat.ID)
	if len(events) == 0 {
		db.CreateEvent(model.Event{Owner: recv.Message.Chat.ID})
		events = db.GetEvent(recv.Message.Chat.ID)
	}
	event := events[0]
	for k, _ := range events {
		if events[k].Name == "" || events[k].Description == "" {
			event = events[k]
		}
	}

	msg := tgapi.NewMessage(recv.Message.Chat.ID, "")
	text := recv.Message.Text
	if text == "/eventcreate" {
		text = ""
	}
	if event.Name == "" {
		if text == "" {
			msg.Text = "What is the name of the event?"
		} else {
			event.Name = text
			db.UpdateEvent(event)
			msg.Text = "Done, now the description"
		}
	} else if event.Description == "" {
		if text == "" {
			msg.Text = "Give a description for your event"
		} else {
			event.Name = ""
			event.Description = text
			db.UpdateEvent(event)
			msg.Text = "We're done"
			db.UserCookieDelete(cookie)
		}
	}

	var chat tgapi.Chattable
	chat = msg
	return &chat, nil
}
