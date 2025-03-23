package entity

import (
	"time"
)

type Message struct {
	MessageId string    `json:"messageId" bson:"_id"`
	UserId    uint64    `json:"userId" bson:"userId"`
	SupportId uint64    `json:"supportId" bson:"supportId"`
	SenderId  uint64    `json:"senderId" bson:"senderId"`
	Content   string    `json:"content" bson:"content"`
	IsRead    bool      `json:"isRead" bson:"isRead"`
	Timestamp time.Time `json:"timestamp" bson:"timestamp"`
}

type Support struct {
	SupportId uint64 `json:"supportId" bson:"supportId"`
	CountChat uint64 `json:"countChat" bson:"countChat"`
}

type SupportToChat struct {
	SupportId uint64 `json:"supportId" bson:"supportId"`
	UserId    uint64 `json:"userId" bson:"userId"`
}
