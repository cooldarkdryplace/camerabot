package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/cooldarkdryplace/camerabot/handler"
	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	fallbackTimeout = 20 * time.Second
	defaultCacheDir = "/tmp"
)

var (
	handlers map[string]Handler

	mu           sync.Mutex
	lastUpdateID int64

	mainChatID int64
	cacheDir   string
)

// Handler processes command sent to bot.
type Handler interface {
	Handle(chatID int64) error
	Command() string
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

		if h, exists := handlers[cmd]; exists {
			h.Handle(chatID)
			continue
		}

		log.Printf("Unknown command: %q in chat: %d ignored", cmd, u.Message.Chat.ID)
	}
}

func main() {
	if v := os.Getenv("MAIN_CHAT_ID"); v != "" {
		var err error
		if mainChatID, err = strconv.ParseInt(v, 10, 64); err != nil {
			log.Fatalf("Main chat ID is not a valid integer: %s", err)
		}
	} else {
		log.Fatal("MAIN_CHAT_ID env var not set")
	}

	if v := os.Getenv("CACHE_DIR"); v != "" {
		cacheDir = v
	} else {
		log.Printf("Using default cache directory: %s", defaultCacheDir)
		cacheDir = defaultCacheDir
	}

	picHandler := handler.NewPictureHandler(cacheDir)
	zoomHandler := handler.NewZoomHandler(cacheDir)

	handlers = map[string]Handler{
		picHandler.Command():  picHandler,
		zoomHandler.Command(): zoomHandler,
	}

	telegram.SendTextMessage(mainChatID, "Hi there.")

	for {
		updates, err := telegram.GetUpdates(lastUpdateID + 1)
		if err != nil {
			telegram.SendTextMessage(mainChatID, fmt.Sprintf("Failed getting updates: %v", err))
			time.Sleep(fallbackTimeout)
		}

		log.Print("Polling...")
		handleUpdates(updates)
	}
}
