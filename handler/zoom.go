package handler

import (
	"log"
	"os/exec"

	"github.com/cooldarkdryplace/camerabot"
	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	zoomedPhoto = "zoomedFrame.jpg"
	zoomScript  = "/opt/camerabot/updateZoomFrame.sh"
	zoomCommand = "/zoom"
)

func init() {
	camerabot.Handlers[zoomCommand] = &ZoomHandler{}
}

type ZoomHandler struct{}

func (zh *ZoomHandler) Command() string {
	return zoomCommand
}

func (zh *ZoomHandler) Help() string {
	return "Make zoomed photo and send it to the chat."
}

func (zh *ZoomHandler) Handle(chatId int64) error {
	if err := exec.Command(zoomScript).Run(); err != nil {
		log.Print("Failed generating new zoomed photo: ", err)
		return err
	}

	go telegram.SendPicture(chatId, camerabot.CacheDir+"/"+zoomedPhoto)

	return nil
}
