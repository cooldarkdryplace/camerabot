package main

import (
	"log"
	"time"

	"github.com/bilinguliar/camerabot/telegram"
	"strings"
)

const (
	chatId = -136923106
)

var lastUpdate int

func main() {
	for {
		processUpdates(getUpdates())
		time.Sleep(time.Second * 10)
		log.Print("Main sleeping...")
	}
}

func getUpdates() []telegram.Update {
	log.Println("Getting updates.")
	return telegram.GetUpdates()
}

func processUpdates(updates []telegram.Update) {
	for _, u := range updates {
		log.Printf("COmparing update ID #%v with last: %v", u.ID, lastUpdate)

		if !shouldBeProcessed(u) {
			continue
		}

		log.Printf("Message type: %s", u.Message.Entities[0].Type)
		if u.Message.Entities[0].Type == "bot_command" {

			if strings.Contains(u.Message.Text, "/pic") {
				go sayHi()
			}
		}

		keepTrackOfUpdates(u.ID)
	}
}

func shouldBeProcessed(u telegram.Update) bool {
	if u.ID <= lastUpdate || len(u.Message.Entities) == 0 {
		return false
	}

	return true
}

func keepTrackOfUpdates(id int) {
	if id > lastUpdate {
		log.Println("Updating last")
		lastUpdate = id
	}
}

func sayHi() {
	log.Print("Saying hi.")

	telegram.SendTextMessage(chatId, "Hi there.")
}

func sendPhoto() {

}
