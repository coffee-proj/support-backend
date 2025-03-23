package dto

import (
	"time"

	types "github.com/coffee/support/internal/entity/type"
)

type Message struct {
	MessageId string    `json:"message_id"`
	UserId    uint64    `json:"user_id"`
	SupportId uint64    `json:"support_id"`
	SenderId  uint64    `json:"sender_id" `
	Content   string    `json:"content"`
	IsRead    bool      `json:"is_read"`
	Timestamp time.Time `json:"timestamp"`
}

type WsMessage struct {
	ChatId   uint64        `json:"-"`
	SenderId uint64        `json:"-" `
	Content  string        `json:"content"`
	Type     types.MsgType `json:"type"`
}

type Support struct {
	SupportId uint64 `json:"supportId"`
	CountChat uint64 `json:"countChat"`
}
