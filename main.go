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
	chatUpdatesMap map[int64]*ChatStatus
	Client         connection.Client
)

func init() {
	chatUpdatesMap = make(map[int64]*ChatStatus)

	Client = &connection.HttpClient{
		Impl: &http.Client{},
	}

	go sayHi()
}

func main() {
	for {
		updates := getUpdates()
		chatUpdatesMap = setChatStatuses(chatUpdatesMap, updates)
		sendPictures(chatUpdatesMap)

		time.Sleep(time.Second * 10)
		log.Print("Main sleeping...")
	}
}

func getUpdates() []telegram.Update {
	log.Println("Getting updates.")

	return telegram.GetUpdates(Client)
}

func sendPictures(chatsMap map[int64]*ChatStatus) {
	for chatID, status := range chatsMap {
		if status.WillSend {
			sendPhoto(chatID)
		}
	}
}

type ChatStatus struct {
	LastProcessed int64
	WillSend      bool
}

func setChatStatuses(chatUpdatesMap map[int64]*ChatStatus, updates []telegram.Update) map[int64]*ChatStatus {
	for _, u := range updates {
		if isUpdateContainsPicRequest(u) {
			status, present := chatUpdatesMap[u.Message.Chat.ID]

			if present {
				if status.LastProcessed < u.ID {
					status.LastProcessed = u.ID
					status.WillSend = true
				} else {
					status.WillSend = false
				}
			} else {
				chatUpdatesMap[u.Message.Chat.ID] = &ChatStatus{
					LastProcessed: u.ID,
					WillSend:      true,
				}
			}
		}
	}

	return chatUpdatesMap
}

func isUpdateContainsPicRequest(u telegram.Update) bool {
	return u.Message.Entities[0].Type == "bot_command" && strings.Contains(u.Message.Text, "/pic")
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
