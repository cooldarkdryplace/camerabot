package handler

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/VictoriaMetrics/metrics"

	"github.com/cooldarkdryplace/camerabot"
	"github.com/cooldarkdryplace/camerabot/telegram"
)

const (
	sourcePhoto = "frame.jpg"
	picScript   = "/opt/camerabot/updateFrame.sh"
	picCommand  = "/pic"
)

var picsSentTotal = metrics.NewCounter("pictures_sent_total")

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
		log.Printf("Failed generating new photo: %s", err)
		return err
	}

	err := telegram.SendPicture(chatID, filepath.Join(camerabot.CacheDir, sourcePhoto))
	if err != nil {
		return fmt.Errorf("failed to send pictures: %w", err)
	}
	picsSentTotal.Inc()

	return nil
}
