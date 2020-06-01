package multichoice

import (
	"errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CMDName() string {
	return "multi"
}

// MultiChoice prints the JSON encoded "message" field in the body
// of the request or "Hello, World!" if there isn't one.
func MultiChoice(tgbotapi.Update) (*tgbotapi.Message, error) {
	return nil, errors.New("unimplemented")
}
