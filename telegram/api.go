package telegram

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/bilinguliar/camerabot/connection"
)

const (
	baseURL = "https://api.telegram.org/bot"

	methodSendMessage    = "sendMessage"
	methodSendPhoto      = "sendPhoto"
	methodGetUpdates     = "getUpdates"
	methodsendChatAction = "sendChatAction"
)

var token string

func init() {
	token = os.Getenv("TOKEN")
}

func GetUpdates(client connection.Client) ([]Update, error) {
	apiResponse := &UpdatesResponse{}

	err := getJson(client, fmt.Sprintf("%s%s/%s", baseURL, token, methodGetUpdates), apiResponse)

	if err != nil {
		log.Print("Error getting updates", err)
		return nil, err
	}

	return apiResponse.Updates, nil
}

func SendTextMessage(client connection.Client, chat int64, m string) {
	log.Printf("Sending test message: %s to chat: %v", m, chat)
	client.Get(fmt.Sprintf("%s%s/%s?chat_id=%v&text=%s", baseURL, token, methodSendMessage, chat, m))
}

func SendPicture(client connection.Client, chat int64, filename string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	picture, err := os.Open(filename)
	if err != nil {
		log.Panic("error opening file")
	}

	defer picture.Close()

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

	res, err := client.Do(req)
	if err != nil {
		log.Panic(err)
	}

	if res.StatusCode != http.StatusOK {
		log.Panic(fmt.Errorf("bad status: %s", res.Status))
	}
}

func getJson(client connection.Client, url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}
