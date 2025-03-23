package support

import (
	"github.com/coffee/support/internal/entity"
	"github.com/gosuit/e"
	"github.com/gosuit/lec"
)

type SupportStorage interface {
	GetAllSupports(ctx lec.Context, offset, limit int64) ([]entity.Support, e.Error)
	AddSupport(ctx lec.Context, supportId uint64) e.Error
	ReplaceSupport(ctx lec.Context, targetId, replacementId uint64) e.Error
	RemoveSupport(ctx lec.Context, supportId uint64) e.Error
}

type SupportChooser interface {
	ChooseSupport(ctx lec.Context) (entity.Support, e.Error)
}
