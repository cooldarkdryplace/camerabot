package handler

import (
	"log"
	"os/exec"

	"github.com/cooldarkdryplace/camerabot/connection"
	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	zoomScript  = "/opt/camerabot/updateZoomFrame.sh"
	zoomCommand = "/zoom"
)

type ZoomHandler struct {
	command       string
	script        string
	photoLocation string
}

func NewZoomHandler(photoLocation string) *ZoomHandler {
	return &ZoomHandler{
		command:       zoomCommand,
		script:        zoomScript,
		photoLocation: photoLocation,
	}
}

func (zh ZoomHandler) Handle(client connection.Client, chatId int64) error {
	if err := exec.Command(zh.script).Run(); err != nil {
		log.Print("Failed generating new zoomed photo: ", err)
		return err
	}

	go telegram.SendPicture(client, chatId, zh.photoLocation)

	return nil
}

func (zh ZoomHandler) GetCommand() string {
	return zh.command
}
