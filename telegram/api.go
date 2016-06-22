package telegram

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
)

const (
	baseURL = "https://api.telegram.org/bot"

	token = "181285124:AAEp5UShB5s7LyDMqJGWBDFR_DeBtBUlBXE"

	methodSendMessage = "sendMessage"
	methodSendPhoto = "sendPhoto"
	methodGetUpdates = "getUpdates"
	methodsendChatAction = "sendChatAction"
)

type Entity struct {
	Type   string `json:"type"`
	Offset int `json:"offset"`
	Length int `json:"length"`
}

type User struct {
	ID        int `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

type UpdatesResponse struct {
	Ok      bool `json:"ok"`
	Updates []Update `json:"result"`
}

type Chat struct {
	ID    int `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type Message struct {
	ID       int `json:"message_id"`
	Date     int `json:"date"`
	Chat     Chat `json:"chat"`
	Entities []Entity `json:"entities"`
	Text     string `json:"text"`
	From     User `json:"from"`
}

type Update struct {
	ID      int `json:"update_id"`
	Message Message `json:"message"`
}

type PhotoSize struct {
	ID       int `json:"file_id"`
	width    int `json:"width"`
	height   int `json:"height"`
	fileSize int `json:"file_size"`
}

func GetUpdates() []Update {
	apiResponse := new(UpdatesResponse)
	err := getJson(fmt.Sprintf("%s%s/%s", baseURL, token, methodGetUpdates), apiResponse)

	if err != nil {
		log.Panic("Error getting updates", err)
	}

	return apiResponse.Updates
}

func SendTextMessage(chat int, m string) {
	log.Printf("Sending test message: %s to chat: %v", m, chat)
	http.Get(fmt.Sprintf("%s%s/%s?chat_id=%v&text=%s", baseURL, token, methodSendMessage, chat, m))
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}