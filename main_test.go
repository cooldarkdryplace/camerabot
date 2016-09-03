package main

import (
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
						ID:    3667788,
						Title: "AnotherChat",
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
	chatStatuses := make(map[int64]*ChatStatus)
	expected := map[int64]*ChatStatus{
		2556677: &ChatStatus{
			LastProcessed: 77665544332211,
			WillSend:      true,
		},
		3667788: &ChatStatus{
			LastProcessed: 111222333444,
			WillSend:      true,
		},
	}

	actual := setChatStatuses(chatStatuses, updates.Updates)

	if len(expected) != len(actual) {
		t.Error("Resulting maps differ in length")
	}

	for k, _ := range expected {
		actualStatus := actual[k]
		expectedStatus := expected[k]

		if actualStatus.LastProcessed != expectedStatus.LastProcessed {
			t.Errorf("Expected last processed message to be: %d but was %d", expectedStatus.LastProcessed, actualStatus.LastProcessed)
		}
	}
}
