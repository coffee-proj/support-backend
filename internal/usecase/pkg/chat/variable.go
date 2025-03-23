package chat

import (
	"github.com/coffee/support/internal/entity"
	"github.com/gosuit/e"
	"github.com/gosuit/lec"
)

type MessageSaver interface {
	WriteMessage(lec.Context, entity.Message) e.Error
	GetChat(lec.Context, uint64) ([]*entity.Message, e.Error)
	GetSupportChats(lec lec.Context, supportId uint64) ([]uint64, e.Error)
	GetSupportIdFromUser(lec lec.Context, userId uint64) (uint64, e.Error)
	SetSupportToChat(lec lec.Context, sc entity.SupportToChat) e.Error
	ChooseSupport(lec lec.Context) (entity.Support, e.Error)
	ReadChat(lec lec.Context, userId, readerId uint64) e.Error
}
