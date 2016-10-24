package main

import (
	"log"
	"net/http"
	"time"

	"github.com/bilinguliar/camerabot/connection"
	"github.com/bilinguliar/camerabot/telegram"
	"github.com/bilinguliar/camerabot/handler"
)

const (
	mainChatId  int64 = -1001077692103
	sourcePhoto       = "/tmp/frame.png"
)

var (
	chatUpdatesMap map[int64]*ChatStatus
	Client         connection.Client
	handlers [1]handler.Handler
)

func init() {
	chatUpdatesMap = make(map[int64]*ChatStatus)

	Client = &connection.HttpClient{
		Impl: &http.Client{},
	}

	handlers = [...]handler.Handler{
		handler.NewPictureHandler(sourcePhoto),
	}

	sayHi()
}

func main() {
	for {
		updates := getUpdates()
		chatUpdatesMap = setChatStatuses(chatUpdatesMap, updates)

		handleUpdates(chatUpdatesMap)

		time.Sleep(time.Second * 10)
	}
}

func getUpdates() []telegram.Update {
	return telegram.GetUpdates(Client)
}

func handleUpdates(chatsMap map[int64]*ChatStatus) {
	for chatID, status := range chatsMap {
		if !status.WillSend {
			continue
		}

		for _, h := range handlers {
			if status.Command == h.GetCommand() {
				h.Handle(Client, chatID)
			}
		}
	}
}

type ChatStatus struct {
	LastProcessed int64
	Command string
	WillSend      bool
}

func setChatStatuses(chatUpdatesMap map[int64]*ChatStatus, updates []telegram.Update) map[int64]*ChatStatus {
	for _, u := range updates {
		if !isUpdateContainsCommand(u) {
			continue
		}

		status, present := chatUpdatesMap[u.Message.Chat.ID]

		if present {
			if status.LastProcessed < u.ID {
				status.LastProcessed = u.ID
				status.WillSend = true
				status.Command = u.Message.Text
			} else {
				status.WillSend = false
			}
		} else {
			chatUpdatesMap[u.Message.Chat.ID] = &ChatStatus{
				LastProcessed: u.ID,
				WillSend:      true,
				Command: u.Message.Text,
			}
		}
	}

	return chatUpdatesMap
}

func isUpdateContainsCommand(u telegram.Update) bool {
	return u.Message.Entities[0].Type == "bot_command"
}

func sayHi() {
	log.Print("Saying hi.")

	telegram.SendTextMessage(Client, mainChatId, "Hi there.")
}

