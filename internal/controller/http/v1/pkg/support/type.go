package support

import (
	"github.com/coffee/support/internal/entity"
	e "github.com/gosuit/e"
	"github.com/gosuit/lec"
)

type ChatUseCase interface {
	WriteMessage(c lec.Context, msg *entity.Message) e.Error
	GetChat(c lec.Context, userId uint64) ([]*entity.Message, e.Error)
	GetSupportChats(c lec.Context, id uint64) ([]uint64, e.Error)
	ReadChat(c lec.Context, userId uint64, readerId uint64) e.Error
}

type SupportUseCase interface {
	GetAllSupports(ctx lec.Context, offset, limit int64) ([]entity.Support, e.Error)
	AddSupport(c lec.Context, id uint64) e.Error
	RemoveSupport(c lec.Context, id uint64) e.Error
}
