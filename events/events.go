package events

import (
	"fmt"
	"strconv"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sir-farfan/hack-a-bot/model"
	"github.com/sir-farfan/hack-a-bot/service/sqlevent"
)

// CMDName - the name for this command
func CMDName() string {
	return "events"
}

// RegisterCommands - call in order to add new commands to the help section
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
		id := fmt.Sprintf("%d", evt.ID)
		button := tgapi.InlineKeyboardButton{Text: evt.Name,
			CallbackData: &id,
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
	db := sqlevent.New()

	return db.GetEvent(0)
}

// Create - will only put the user in "event creation mode"
func Create(recv tgapi.Update) (*tgapi.Chattable, error) {
	db := sqlevent.New()
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
	for k := range events {
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
	} else {
		msg.Text = "used up your events"
		db.UserCookieDelete(cookie)
	}

	var chat tgapi.Chattable
	chat = msg
	return &chat, nil
}

// Subscribe - Allows the user to subscribe to the selected event
func Subscribe(recv tgapi.Update) (*tgapi.Chattable, error) {
	db := sqlevent.New()
	defer db.DB.Close()

	msg := tgapi.NewMessage(recv.CallbackQuery.Message.Chat.ID, "")
	msg.ParseMode = "Markdown"

	eventID, _ := strconv.Atoi(recv.CallbackQuery.Data)
	events := db.GetEventByID(int64(eventID))
	if len(events) != 1 {
		msg.Text = "Error looking up event"
	} else {
		msg.Text = fmt.Sprintf("This is what we know about the event:\n%s", events[0].Description)
	}

	var chat tgapi.Chattable
	chat = msg
	return &chat, nil
}
