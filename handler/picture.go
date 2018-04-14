package handler

import (
	"log"
	"os/exec"

	"github.com/cooldarkdryplace/camerabot"
	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	sourcePhoto = "frame.jpg"
	picScript   = "/opt/camerabot/updateFrame.sh"
	picCommand  = "/pic"
)

func init() {
	camerabot.Handlers[picCommand] = &PictureHandler{}
}

type PictureHandler struct{}

func (ph *PictureHandler) Command() string {
	return picCommand
}

func (ph *PictureHandler) Help() string {
	return "Make a full size photo and send it to the chat."
}

func (ph *PictureHandler) Handle(chatID int64) error {
	if err := exec.Command(picScript).Run(); err != nil {
		log.Print("Failed generating new photo: ", err)
		return err
	}

	go telegram.SendPicture(chatID, camerabot.CacheDir+"/"+sourcePhoto)

	return nil
}
