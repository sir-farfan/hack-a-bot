package model

import tgapi "github.com/go-telegram-bot-api/telegram-bot-api"

type CommandProcessor struct {
	CMDName string
	Process Processor
}

type Processor func(recv tgapi.Update) (*tgapi.Chattable, error)
