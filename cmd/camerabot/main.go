package main

import (
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/cooldarkdryplace/camerabot"
)

const defaultCacheDir = "/tmp"

func main() {
	if v := os.Getenv("MAIN_CHAT_ID"); v != "" {
		var err error
		if camerabot.MainChatID, err = strconv.ParseInt(v, 10, 64); err != nil {
			log.Fatalf("Main chat ID is not a valid integer: %s", err)
		}
	} else {
		log.Fatal("MAIN_CHAT_ID env var not set")
	}

	if v := os.Getenv("CACHE_DIR"); v != "" {
		camerabot.CacheDir = v
	} else {
		log.Printf("Using default cache directory: %s", defaultCacheDir)
		camerabot.CacheDir = defaultCacheDir
	}

	go camerabot.ListenAndServe()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-signalChan
	log.Println("Interrupt recieved. Graceful shutdown.")
}
