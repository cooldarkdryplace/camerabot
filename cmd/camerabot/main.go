package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cooldarkdryplace/camerabot/connection"
	"github.com/cooldarkdryplace/camerabot/handler"
	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	mainChatID      int64 = -1001077692103
	sourcePhoto           = "/tmp/frame.jpg"
	zoomedPhoto           = "/tmp/zoomedFrame.jpg"
	fallbackTimeout       = 20
)

var (
	handlers     map[string]Handler
	lastUpdateID int64
)

// Handler processes command sent to bot.
type Handler interface {
	Handle(client connection.Client, chatID int64) error
	GetCommand() string
}

func init() {
	handlers = make(map[string]Handler)

	picHandler := handler.NewPictureHandler(sourcePhoto)
	zoomHandler := handler.NewZoomHandler(zoomedPhoto)

	handlers[picHandler.GetCommand()] = picHandler
	handlers[zoomHandler.GetCommand()] = zoomHandler
}

func main() {
	client := &connection.HttpClient{
		Impl: &http.Client{},
	}

	sayHi(client)

	for {
		updates, err := getUpdates(client)
		if err != nil {
			telegram.SendTextMessage(client, mainChatID, fmt.Sprintf("Failed getting updates: %v", err))
			time.Sleep(fallbackTimeout * time.Second)
		}

		log.Print("Polling...")
		handleUpdates(updates, client)
	}
}

func getUpdates(client connection.Client) ([]telegram.Update, error) {
	return telegram.GetUpdates(client, lastUpdateID+1)
}

func handleUpdates(updates []telegram.Update, client connection.Client) {
	for _, u := range updates {
		trackLastUpdateID(u.ID)

		command := getCommand(u)
		chatID := u.Message.Chat.ID

		if command == "" {
			continue
		}

		if h, exists := handlers[command]; exists {
			h.Handle(client, chatID)
			continue
		}

		log.Printf("Unknown command: %q ignored", command)
	}
}

func getCommand(u telegram.Update) string {
	if len(u.Message.Entities) == 0 {
		return ""
	}

	if u.Message.Entities[0].Type == "bot_command" {
		return u.Message.Text
	}

	return ""
}

func trackLastUpdateID(ID int64) {
	log.Printf("Last update ID: %d, incoming update ID: %d", lastUpdateID, ID)
	if lastUpdateID < ID {
		lastUpdateID = ID
	}
}

func sayHi(client connection.Client) {
	log.Print("Saying hi.")

	telegram.SendTextMessage(client, mainChatID, "Hi there.")
}
