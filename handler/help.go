package handler

import (
	"strings"

	"github.com/cooldarkdryplace/camerabot"
	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	helpCommand = "/help"
	endOfLine   = "%0A"
)

func init() {
	camerabot.Handlers[helpCommand] = &HelpHandler{}
}

type HelpHandler struct{}

func (hh *HelpHandler) Command() string {
	return helpCommand
}

func (hh *HelpHandler) Help() string {
	return "Show this help"
}

func (hh *HelpHandler) Handle(chatID int64) error {
	var b strings.Builder

	for k, v := range camerabot.Handlers {
		b.WriteString(k)
		b.WriteString(" - ")
		b.WriteString(v.Help())
		b.WriteString(endOfLine)
	}

	return telegram.SendTextMessage(chatID, b.String())
}
