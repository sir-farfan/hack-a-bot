package main

import (
	"fmt"
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sir-farfan/hack-a-bot/multichoice"
)

type CommandProcessor struct {
	CMDName string
	Process Processor
}

type Processor func(recv tgbotapi.Update) (*tgbotapi.Message, error)

func main() {
	processors := make(map[string]Processor)

	processors[multichoice.CMDName()] = multichoice.MultiChoice

	fmt.Println(processors)
	log.Printf("hello world 14\n")
	tgToken := os.Getenv("BOT_TOKEN")
	log.Printf("bot key: %s\n", tgToken)

	bot, err := tgbotapi.NewBotAPI(tgToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		if update.Message.IsCommand() {
			exec := processors[update.Message.Command()]
			if exec != nil {
				response, _ := exec(update)
				fmt.Println(response)
			}
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		msg.ReplyToMessageID = update.Message.MessageID

		bot.Send(msg)
	}
}
