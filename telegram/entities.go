package telegram

type Entity struct {
	Type   string `json:"type"`
	Offset int64  `json:"offset"`
	Length int64  `json:"length"`
}

type User struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	UserName  string `json:"username"`
}

type UpdatesResponse struct {
	Ok      bool     `json:"ok"`
	Updates []Update `json:"result"`
}

type Chat struct {
	ID    int64  `json:"id"`
	Title string `json:"title"`
	Type  string `json:"type"`
}

type Message struct {
	ID       int64    `json:"message_id"`
	Date     int      `json:"date"`
	Chat     Chat     `json:"chat"`
	Entities []Entity `json:"entities"`
	Text     string   `json:"text"`
	From     User     `json:"from"`
}

type Update struct {
	ID      int64   `json:"update_id"`
	Message Message `json:"message"`
}

type PhotoSize struct {
	ID       int64 `json:"file_id"`
	Width    int   `json:"width"`
	Height   int   `json:"height"`
	FileSize int64 `json:"file_size"`
}
