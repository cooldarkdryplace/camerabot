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
	log.Println("Getting updates.")
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
	log.Printf("COmparing update ID #%v with last: %v", u.ID, lastUpdate)

	if u.ID > lastUpdate {
		log.Printf("Processing update: %v", u.ID)
		processed <- u.ID

		log.Printf("Message type: %s", u.Message.Entities[0].Type)

		if u.Message.Entities[0].Type == "bot_command" {
			sayHi()
		}
	}
}

func keepTrackOfUpdates(processed chan int) {
	log.Print("Keeping track of updates")

	p := <- processed

	if p > lastUpdate {
		log.Println("Updating last")
		lastUpdate = p
	}
}

func sayHi() {
	log.Print("Saying hi.")

	telegram.SendTextMessage(chatId, "Hi there.")
}

func sendPhoto() {

}
