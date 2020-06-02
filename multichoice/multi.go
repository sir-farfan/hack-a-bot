package multichoice

import (
	"errors"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func CMDName() string {
	return "multi"
}

// MultiChoice - returns a multichoice event for the user to pick
func MultiChoice(recv tgapi.Update) (*tgapi.Chattable, error) {
	return nil, errors.New("unimplemented")
}
