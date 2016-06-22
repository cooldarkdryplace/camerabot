package main

import (
	"log"

	"github.com/bilinguliar/camerabot/telegram"
	"time"
)

func main() {
	log.Print("Hi")

	updatesChan := make(chan telegram.Update)
	go processUpdates(updatesChan)
	go monitorUpdates(updatesChan)

	for {
		time.Sleep(time.Second * 10)
		log.Print("Main sleeping...")
	}
}

func monitorUpdates(c chan telegram.Update) {
	for {
		for _, u:= range telegram.GetUpdates() {
			c <- u
		}

		log.Print("Updates mon sleeping...")
		time.Sleep(time.Second * 5)
	}
}

func processUpdates(c chan telegram.Update) {

}

func sayHi() {

}

func sendPhoto() {

}
