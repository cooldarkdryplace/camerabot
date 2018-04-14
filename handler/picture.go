package handler

import (
	"log"
	"os/exec"

	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	sourcePhoto = "frame.jpg"
	picScript   = "/opt/camerabot/updateFrame.sh"
	picCommand  = "/pic"
)

type PictureHandler struct {
	command       string
	script        string
	photoLocation string
}

func NewPictureHandler(cacheDir string) *PictureHandler {
	return &PictureHandler{
		command:       picCommand,
		script:        picScript,
		photoLocation: cacheDir + "/" + sourcePhoto,
	}
}

func (ph PictureHandler) Handle(chatID int64) error {
	if err := exec.Command(ph.script).Run(); err != nil {
		log.Print("Failed generating new photo: ", err)
		return err
	}

	go telegram.SendPicture(chatID, ph.photoLocation)

	return nil
}

func (ph PictureHandler) Command() string {
	return ph.command
}
