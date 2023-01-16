package models

type ChatMessage struct {
	Channel uint32 `json:"channel"`
	Msg     string `json:"msg"`
}
