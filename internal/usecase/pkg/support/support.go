package support

import (
	"github.com/coffee/support/internal/entity"
	"github.com/gosuit/e"
	"github.com/gosuit/lec"
)

type Support struct {
	ss SupportStorage
	sc SupportChooser
}

func New(ss SupportStorage, sc SupportChooser) *Support {
	return &Support{
		ss: ss,
		sc: sc,
	}
}

func (s *Support) AddSupport(ctx lec.Context, supportId uint64) e.Error {
	return s.ss.AddSupport(ctx, supportId)
}

func (s *Support) RemoveSupport(ctx lec.Context, supportId uint64) e.Error {
	err := s.ss.RemoveSupport(ctx, supportId)
	if err != nil {
		return err
	}

	newSupport, err := s.sc.ChooseSupport(ctx)
	if err != nil {
		return err
	}

	return s.ss.ReplaceSupport(ctx, supportId, newSupport.SupportId)
}

func (s *Support) GetAllSupports(ctx lec.Context, offset, limit int64) ([]entity.Support, e.Error) {
	return s.ss.GetAllSupports(ctx, offset, limit)
}
