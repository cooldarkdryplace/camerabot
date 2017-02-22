package handler

import (
	"log"
	"os/exec"

	"github.com/bilinguliar/camerabot/connection"
	"github.com/bilinguliar/camerabot/telegram"
)

const (
	picScript  = "/opt/camerabot/updateFrame.sh"
	picCommand = "/pic"
)

type PictureHandler struct {
	command       string
	script        string
	photoLocation string
}

func NewPictureHandler(photoLocation string) *PictureHandler {
	return &PictureHandler{
		command:       picCommand,
		script:        picScript,
		photoLocation: photoLocation,
	}
}

func (ph PictureHandler) Handle(client connection.Client, chatId int64) error {
	if err := exec.Command(ph.script).Run(); err != nil {
		log.Print("Failed generating new photo: ", err)
		return err
	}

	go telegram.SendPicture(client, chatId, ph.photoLocation)

	return nil
}

func (ph PictureHandler) GetCommand() string {
	return ph.command
}
