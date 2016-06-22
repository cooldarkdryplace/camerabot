package main

import (
	"log"
	"time"

	"github.com/bilinguliar/camerabot/telegram"
)

const (
	chatId = -136923106
)

var lastUpdate int

func main() {
	log.Print("Starting...")
	processedChan := make (chan int)
	updatesChan := make(chan telegram.Update)
	go startConsumer(updatesChan, processedChan)
	go keepTrackOfUpdates(processedChan)

	log.Print("Started")

	for {
		getUpdates(updatesChan)
		time.Sleep(time.Second * 10)
		log.Print("Main sleeping...")
	}
}

func getUpdates(out chan telegram.Update) {
	for _, u := range telegram.GetUpdates() {
		out <- u
	}
}

func startConsumer(updates chan telegram.Update, processed chan int) {
	log.Print("Initing consumer...")

	for {
		go processUpdate(<-updates, processed)
	}
}

func processUpdate(u telegram.Update, processed chan int) {
	log.Printf("Processing update #%s", u.ID)

	if u.ID > lastUpdate {

		processed <- u.ID
		if u.Message.Entities[0].Type == "bot_command" {
			sayHi()
		}
	}
}

func keepTrackOfUpdates(processed chan int) {
	log.Print("Keeping track of updates")

	p := <- processed

	if p > lastUpdate {
		lastUpdate = p
	}
}

func sayHi() {
	log.Print("Saying hi.")

	telegram.SendTextMessage(chatId, "Hi there.")
}

func sendPhoto() {

}
