package handler

import (
	"log"
	"os/exec"

	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	zoomedPhoto = "zoomedFrame.jpg"
	zoomScript  = "/opt/camerabot/updateZoomFrame.sh"
	zoomCommand = "/zoom"
)

type ZoomHandler struct {
	command       string
	script        string
	photoLocation string
}

func NewZoomHandler(cacheDir string) *ZoomHandler {
	return &ZoomHandler{
		command:       zoomCommand,
		script:        zoomScript,
		photoLocation: cacheDir + "/" + zoomedPhoto,
	}
}

func (zh ZoomHandler) Handle(chatId int64) error {
	if err := exec.Command(zh.script).Run(); err != nil {
		log.Print("Failed generating new zoomed photo: ", err)
		return err
	}

	go telegram.SendPicture(chatId, zh.photoLocation)

	return nil
}

func (zh ZoomHandler) Command() string {
	return zh.command
}
