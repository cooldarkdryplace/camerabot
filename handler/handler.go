package handler

import "github.com/bilinguliar/camerabot/connection"

type Handler interface {
	Handle(client connection.Client, chatID int64) error
	GetCommand() string
}
