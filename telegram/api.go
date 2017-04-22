package telegram

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
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

	// Wait timeout for longpolling
	timeout = 60
)

var token string

func init() {
	token = os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("API token was not set as ENV variable with name: 'TOKEN'")
	}
}

func GetUpdates(client connection.Client, lastMsgID int64) ([]Update, error) {
	apiResponse := &UpdatesResponse{}

	err := getJson(
		client,
		fmt.Sprintf("%s%s/%s?timeout=%d&offset=%d", baseURL, token, methodGetUpdates, timeout, lastMsgID),
		apiResponse,
	)
	if err != nil {
		log.Println("Error getting updates: ", err)
		return nil, err
	}

	return apiResponse.Updates, nil
}

func SendTextMessage(client connection.Client, chat int64, m string) {
	log.Printf("Sending text message: %q to chat: %v", m, chat)
	client.Get(fmt.Sprintf("%s%s/%s?chat_id=%v&text=%s", baseURL, token, methodSendMessage, chat, m))
}

func SendPicture(client connection.Client, chat int64, filename string) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	picture, err := os.Open(filename)
	if err != nil {
		log.Println("Error opening file: ", err)
		return
	}

	defer picture.Close()

	fileWriter, err := bodyWriter.CreateFormFile("photo", "img.png")
	if err != nil {
		log.Println("Error writing to buffer: ", err)
		return
	}

	_, err = io.Copy(fileWriter, picture)
	if err != nil {
		log.Println("Error copying file: ", err)
		return
	}

	fieldWriter, err := bodyWriter.CreateFormField("chat_id")
	if err != nil {
		log.Println("Error writing chat_id to buffer: ", err)
		return
	}

	err = binary.Write(fieldWriter, binary.LittleEndian, chat)
	if err != nil {
		fmt.Println("Failed to write data as binary to form field: ", err)
		return
	}

	bodyWriter.Close()

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s/%s?chat_id=%v", baseURL, token, methodSendPhoto, chat),
		bodyBuf,
	)
	if err != nil {
		log.Println("Failed to create POST request with picture: ", err)
		return
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		log.Println("Error during POST to Telegram API: ", err)
		return
	}

	if res.StatusCode != http.StatusOK {
		log.Printf("HTTP status for API call was not OK: %s\n", res.Status)
	}
}

func getJson(client connection.Client, url string, target interface{}) error {
	r, err := client.Get(url)
	if err != nil {
		log.Printf("Tried to get conversation updates, error occurred: %q\n", err)
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		log.Printf("HTTP Get failed. Status code: %d\n", r.StatusCode)
		return errors.New("get was not successful")
	}

	if err = json.NewDecoder(r.Body).Decode(target); err != nil {
		log.Println("Failed to unmarshal updates from JSON:", err)
		return err
	}

	return nil
}
