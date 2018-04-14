package camerabot

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/cooldarkdryplace/camerabot/telegram"
)

const fallbackTimeout = 20 * time.Second

var (
	MainChatID int64
	CacheDir   string

	mu           sync.Mutex
	lastUpdateID int64
)

var Handlers = make(map[string]Handler)

// Handler processes command sent to bot.
type Handler interface {
	Command() string
	Handle(chatID int64) error
	Help() string
}

func command(u telegram.Update) string {
	if len(u.Message.Entities) == 0 {
		return ""
	}

	if u.Message.Entities[0].Type == "bot_command" {
		return u.Message.Text
	}

	return ""
}

func trackLastUpdateID(ID int64) {
	mu.Lock()
	log.Printf("Last update ID: %d, incoming update ID: %d", lastUpdateID, ID)
	if lastUpdateID < ID {
		lastUpdateID = ID
	}
	mu.Unlock()
}

func handleUpdates(updates []telegram.Update) {
	for _, u := range updates {
		trackLastUpdateID(u.ID)

		cmd := command(u)
		chatID := u.Message.Chat.ID

		if cmd == "" {
			continue
		}

		if h, exists := Handlers[cmd]; exists {
			h.Handle(chatID)
			continue
		}

		log.Printf("Unknown command: %q in chat: %d ignored", cmd, u.Message.Chat.ID)
	}
}

func ListenAndServe() {
	telegram.SendTextMessage(MainChatID, "Hi there.")

	for {
		updates, err := telegram.GetUpdates(lastUpdateID + 1)
		if err != nil {
			telegram.SendTextMessage(MainChatID, fmt.Sprintf("Failed getting updates: %v", err))
			time.Sleep(fallbackTimeout)
		}

		log.Print("Polling...")
		handleUpdates(updates)
	}
}
