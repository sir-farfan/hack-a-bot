package main

import (
	"log"
	"os"
	"reflect"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sir-farfan/hack-a-bot/events"
	"github.com/sir-farfan/hack-a-bot/model"
	"github.com/sir-farfan/hack-a-bot/multichoice"
	"github.com/sir-farfan/hack-a-bot/service/sql_event"
)

var Processors map[string]model.Processor

func Help(recv tgapi.Update) (*tgapi.Chattable, error) {
	help := tgapi.NewMessage(recv.Message.Chat.ID, "")
	help.Text = "Call an action typing (or clicking) any of:"

	cmds := reflect.ValueOf(Processors).MapKeys()
	log.Println(cmds)
	for _, cmd := range cmds {
		log.Println(cmd.String())
		help.Text += "\n/" + cmd.String()
	}

	var chat tgapi.Chattable
	chat = help

	return &chat, nil
}

func main() {
	Processors = make(map[string]model.Processor)
	Processors["help"] = Help
	Processors[multichoice.CMDName()] = multichoice.MultiChoice
	events.RegisterCommands(Processors)

	db := sql_event.New()

	tgToken := os.Getenv("BOT_TOKEN")

	bot, err := tgapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	u := tgapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		var exec model.Processor
		cookie := db.UserCookieGet(update.Message.Chat.ID)
		if cookie.ID == update.Message.Chat.ID {
			exec = Processors[cookie.Cookie]
		} else if update.Message.IsCommand() {
			exec = Processors[update.Message.Command()]
		}

		if exec != nil {
			response, err := exec(update)
			if err == nil {
				_, err = bot.Send(*response)
				if err != nil {
					log.Printf("ERROR answering the user: %v\n", err)
				}
			}
		} else {
			response, _ := Help(update)
			bot.Send(*response)
		}

		// log.Printf("DEBUG: [%s] %s\n", update.Message.From.UserName, update.Message.Text)
	}
}
