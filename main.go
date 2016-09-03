package main

import (
	"log"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/bilinguliar/camerabot/connection"
	"github.com/bilinguliar/camerabot/telegram"
)

const (
	mainChatId  int64 = -1001077692103
	sourcePhoto       = "/tmp/frame.png"
)

var (
	lastUpdate int64
	Client     connection.Client
)

func init() {
	Client = &connection.HttpClient{
		Impl: &http.Client{},
	}
	go sayHi()
}

func main() {
	for {
		processUpdates(getUpdates())
		time.Sleep(time.Second * 10)
		log.Print("Main sleeping...")
	}
}

func getUpdates() []telegram.Update {
	log.Println("Getting updates.")

	return telegram.GetUpdates(Client)
}

func processUpdates(updates []telegram.Update) {
	for _, u := range updates {
		if !shouldBeProcessed(u) {
			continue
		}

		log.Printf("Message type: %s", u.Message.Entities[0].Type)
		if u.Message.Entities[0].Type == "bot_command" {

			if strings.Contains(u.Message.Text, "/pic") {
				log.Println("Picture requested!")
				go sendPhoto(u.Message.Chat.ID)
			}
		}

		keepTrackOfUpdates(u.ID)
	}
}

func getChatsToSendPictureTo(updates []telegram.Update) map[int64]struct{} {
	chats := make(map[int64]struct{}, len(updates))

	for _, u := range updates {
		if isUpdateContainsPicRequest(u) {
			chats[u.Message.Chat.ID] = struct{}{}
		}
	}

	return chats
}

func isUpdateContainsPicRequest(u telegram.Update) bool {
	return u.Message.Entities[0].Type == "bot_command" && strings.Contains(u.Message.Text, "/pic")
}

func shouldBeProcessed(u telegram.Update) bool {
	if u.ID <= lastUpdate || len(u.Message.Entities) == 0 {
		return false
	}

	return true
}

func keepTrackOfUpdates(id int64) {
	if id > lastUpdate {
		log.Println("Updating last")
		lastUpdate = id
	}
}

func sayHi() {
	log.Print("Saying hi.")

	telegram.SendTextMessage(Client, mainChatId, "Hi there.")
}

func sendPhoto(chatId int64) {
	err := exec.Command("/opt/camerabot/updateFrame.sh").Run()
	if err != nil {
		log.Print(err)
		return
	}

	telegram.SendPicture(Client, chatId, sourcePhoto)
}
