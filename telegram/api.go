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

var (
	token  string
	client = &http.Client{}
)

func init() {
	token = os.Getenv("TOKEN")
	if token == "" {
		log.Fatal("API token was not set as ENV variable with name: 'TOKEN'")
	}
}

// GetUpdates with IDs greater than provided value.
func GetUpdates(lastMsgID int64) ([]Update, error) {
	apiResponse := &UpdatesResponse{}

	err := getJSON(
		fmt.Sprintf("%s%s/%s?timeout=%d&offset=%d", baseURL, token, methodGetUpdates, timeout, lastMsgID),
		apiResponse,
	)
	if err != nil {
		log.Println("Error getting updates: ", err)
		return nil, err
	}

	return apiResponse.Updates, nil
}

// SendTextMessage to the chat with provided ID.
func SendTextMessage(chat int64, m string) error {
	log.Printf("Sending text message: %q to chat: %v", m, chat)
	resp, err := client.Get(fmt.Sprintf("%s%s/%s?chat_id=%v&text=%s", baseURL, token, methodSendMessage, chat, m))
	if err != nil {
		return fmt.Errorf("failed to send message: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Message was not send, status: %s", resp.Status)
		return fmt.Errorf("failed to send message, status: %s", resp.Status)
	}

	return nil
}

// SendPicture to the chat.
func SendPicture(chat int64, filename string) error {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	picture, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}

	defer picture.Close()

	fileWriter, err := bodyWriter.CreateFormFile("photo", "img.png")
	if err != nil {
		return fmt.Errorf("error writing to buffer: %w", err)
	}

	if _, err := io.Copy(fileWriter, picture); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	fieldWriter, err := bodyWriter.CreateFormField("chat_id")
	if err != nil {
		return fmt.Errorf("error writing chat_id to buffer: %w", err)
	}

	if err = binary.Write(fieldWriter, binary.LittleEndian, chat); err != nil {
		return fmt.Errorf("failed to write data as binary to form field: %w", err)
	}
	bodyWriter.Close()

	req, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%s%s/%s?chat_id=%v", baseURL, token, methodSendPhoto, chat),
		bodyBuf,
	)
	if err != nil {
		return fmt.Errorf("failed to create POST request with picture: %w", err)
	}

	req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error during POST to Telegram API: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP status for API call was not OK: %s\n", res.Status)
	}

	return nil
}

func getJSON(url string, target interface{}) error {
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
