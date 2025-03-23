package ws

import (
	"github.com/coffee/support/internal/entity"
	"github.com/gosuit/e"
	"github.com/gosuit/lec"
)

type ChatUseCase interface {
	WriteMessage(c lec.Context, msg *entity.Message) e.Error
	ReadChat(c lec.Context, userId uint64, readerId uint64) e.Error
}
