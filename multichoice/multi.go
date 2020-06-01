package multichoice

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CMDName() string {
	return "multi"
}

// MultiChoice - returns a multichoice event for the user to pick
func MultiChoice(recv tgbotapi.Update) (*tgbotapi.Message, error) {
	return nil, errors.New("unimplemented")
}
