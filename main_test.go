package main

import (
	"reflect"
	"testing"

	"github.com/bilinguliar/camerabot/telegram"
)

var updates telegram.UpdatesResponse

func init() {
	updates = telegram.UpdatesResponse{
		Ok: true,
		Updates: []telegram.Update{
			telegram.Update{
				ID: 11223344556677,
				Message: telegram.Message{
					ID:   22334455667788,
					Date: 123123123,
					Chat: telegram.Chat{
						ID:    2556677,
						Title: "Kilnspotting",
						Type:  "superchat",
					},
					Entities: []telegram.Entity{
						telegram.Entity{
							Type:   "bot_command",
							Offset: 0,
							Length: 4,
						},
					},
					Text: "/pic",
					From: telegram.User{
						ID:        888,
						FirstName: "Eighter",
						UserName:  "eighteighteight",
					},
				},
			},
			telegram.Update{
				ID: 77665544332211,
				Message: telegram.Message{
					ID:   88776655443322,
					Date: 123123124,
					Chat: telegram.Chat{
						ID:    2556677,
						Title: "Kilnspotting",
						Type:  "superchat",
					},
					Entities: []telegram.Entity{
						telegram.Entity{
							Type:   "bot_command",
							Offset: 0,
							Length: 4,
						},
					},
					Text: "/pic",
					From: telegram.User{
						ID:        999,
						FirstName: "Lollercoaster",
						UserName:  "bilinguliar",
					},
				},
			},
			telegram.Update{
				ID: 111222333444,
				Message: telegram.Message{
					ID:   333444555666777,
					Date: 123123130,
					Chat: telegram.Chat{
						ID:    2556677,
						Title: "Kilnspotting",
						Type:  "superchat",
					},
					Entities: []telegram.Entity{
						telegram.Entity{
							Type:   "bot_command",
							Offset: 0,
							Length: 4,
						},
					},
					Text: "/pic",
					From: telegram.User{
						ID:        777,
						FirstName: "RandomCitizen",
						UserName:  "bigbutt",
					},
				},
			},
		},
	}

}

func TestMustProcessSinglePictureRequestPerChatIfThereAreABunchOfThem(t *testing.T) {
	expected := map[int64]struct{}{2556677: struct{}{}}
	actual := getChatsToSendPictureTo(updates.Updates)

	if !reflect.DeepEqual(actual, expected) {
		t.Errorf("Expected: %v but was %v", actual, expected)
	}
}
