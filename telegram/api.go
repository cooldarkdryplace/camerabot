package telegram

import (
	"net/http"
	"fmt"
	"encoding/json"
	"log"
	"io"
	"bytes"
	"mime/multipart"
	"os"
	"encoding/binary"
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

func SendTextMessage(chat int32, m string) {
	log.Printf("Sending test message: %s to chat: %v", m, chat)
	http.Get(fmt.Sprintf("%s%s/%s?chat_id=%v&text=%s", baseURL, token, methodSendMessage, chat, m))
}

func SendPicture(chat int32, filename string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	// open file handle
	picture, err := os.Open(filename)
	if err != nil {
		log.Panic("error opening file")
	}

	defer picture.Close()

	// this step is very important
	fileWriter, err := bodyWriter.CreateFormFile("photo", "img.png")
	if err != nil {
		log.Panic("error writing to buffer")
	}

	_, err = io.Copy(fileWriter, picture)
	if err != nil {
		log.Panic("error reading file", err)
	}

	fieldWriter, err := bodyWriter.CreateFormField("chat_id")
	if err != nil {
		log.Panic("error writing chat_id to buffer")
	}

	err = binary.Write(fieldWriter, binary.LittleEndian, chat)
	if err != nil {
		fmt.Println("binary.Write failed:", err)
	}

	bodyWriter.Close()

	req, err := http.NewRequest(http.MethodPost,
		fmt.Sprintf("%s%s/%s?chat_id=%v", baseURL, token, methodSendPhoto, chat), bodyBuf)
	if err != nil {
		log.Panic(err)
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	// Submit the request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		log.Panic(fmt.Errorf("bad status: %s", res.Status))
	}
}

func getJson(url string, target interface{}) error {
	r, err := http.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}