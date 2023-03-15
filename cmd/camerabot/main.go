package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"

	"github.com/VictoriaMetrics/metrics"

	"github.com/cooldarkdryplace/camerabot"
	_ "github.com/cooldarkdryplace/camerabot/handler"
)

const (
	defaultCacheDir = "/tmp"
	defaultPort     = "8080"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	mainChatID, err := strconv.ParseInt(os.Getenv("MAIN_CHAT_ID"), 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse MAIN_CHAT_ID: %s", err)
	}
	camerabot.MainChatID = mainChatID

	camerabot.CacheDir = os.Getenv("CACHE_DIR")
	if camerabot.CacheDir == "" {
		log.Printf("Using default cache directory: %s", defaultCacheDir)
		camerabot.CacheDir = defaultCacheDir
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Printf("Using default port: %s", defaultPort)
		port = defaultPort
	}

	go camerabot.ListenAndServe()

	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, false)
	})

	errChan := make(chan error)
	go func() {
		errChan <- http.ListenAndServe(":"+port, nil)
	}()

	select {
	case <-ctx.Done():
		log.Println("Interrupt received. Graceful shutdown.")
	case err := <-errChan:
		log.Printf("Camerabot failed: %s", err)
	}
}
