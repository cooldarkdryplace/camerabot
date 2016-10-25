package handler

import (
	"os/exec"
	"log"

	"github.com/bilinguliar/camerabot/telegram"
	"github.com/bilinguliar/camerabot/connection"
)

const (
	picScript = "/opt/camerabot/updateFrame.sh"
	picCommand = "/pic"
)

type PictureHandler struct {
	cmd *exec.Cmd
	command string
	photoLocation string
}

func NewPictureHandler(photoLocation string) *PictureHandler {
	return &PictureHandler {
		cmd: exec.Command(picScript),
		command: picCommand,
		photoLocation: photoLocation,
	}
}

func (ph PictureHandler) Handle(client connection.Client, chatId int64) error {
	if err := ph.cmd.Run(); err != nil {
		log.Print("Failed generating new photo", err)
		return err
	}

	go telegram.SendPicture(client, chatId, ph.photoLocation)

	return nil
}

func (ph PictureHandler) GetCommand() string {
	return ph.command
}
