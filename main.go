package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bilinguliar/camerabot/connection"
	"github.com/bilinguliar/camerabot/handler"
	"github.com/bilinguliar/camerabot/telegram"
)

type ChatStatus struct {
	LastProcessed int64
	Command       string
	WillSend      bool
}

const (
	mainChatId      int64 = -1001077692103
	sourcePhoto           = "/tmp/frame.png"
	fallbackTimeout       = 20
)

var (
	chatUpdatesMap map[int64]*ChatStatus
	client         connection.Client
	handlers       map[string]handler.Handler
	lastUpdateID   int64
)

func init() {
	chatUpdatesMap = make(map[int64]*ChatStatus)

	client = &connection.HttpClient{
		Impl: &http.Client{},
	}

	handlers = make(map[string]handler.Handler)

	picHandler := handler.NewPictureHandler(sourcePhoto)

	handlers[picHandler.GetCommand()] = picHandler
}

func main() {
	sayHi()

	for {
		updates, err := getUpdates()
		if err != nil {
			telegram.SendTextMessage(client, mainChatId, fmt.Sprintf("Failed getting updates: %v", err))
			time.Sleep(fallbackTimeout * time.Second)
		}

		chatUpdatesMap = setChatStatuses(chatUpdatesMap, updates)

		log.Print("Polling...")
		handleUpdates(chatUpdatesMap)
	}
}

func getUpdates() ([]telegram.Update, error) {
	return telegram.GetUpdates(client, lastUpdateID+1)
}

func handleUpdates(chatsMap map[int64]*ChatStatus) {
	for chatID, status := range chatsMap {
		if !status.WillSend {
			continue
		}

		h, exists := handlers[status.Command]
		if exists {
			h.Handle(client, chatID)
			return
		}

		log.Printf("Unknown command: %q ignored", status.Command)
	}
}

func setChatStatuses(chatUpdatesMap map[int64]*ChatStatus, updates []telegram.Update) map[int64]*ChatStatus {
	for _, u := range updates {
		trackLastUpdateID(u.ID)

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
				Command:       u.Message.Text,
			}
		}
	}

	return chatUpdatesMap
}

func isUpdateContainsCommand(u telegram.Update) bool {
	if len(u.Message.Entities) == 0 {
		return false
	}

	return u.Message.Entities[0].Type == "bot_command"
}

func trackLastUpdateID(ID int64) {
	log.Printf("Last update ID: %d, incoming update ID: %d", lastUpdateID, ID)
	if lastUpdateID < ID {
		lastUpdateID = ID
	}
}

func sayHi() {
	log.Print("Saying hi.")

	telegram.SendTextMessage(client, mainChatId, "Hi there.")
}
