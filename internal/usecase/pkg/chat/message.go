package chat

import (
	"time"

	"github.com/coffee/support/internal/entity"
	"github.com/gosuit/e"
	"github.com/gosuit/lec"
	"github.com/gosuit/sl"
)

type Chat struct {
	mw MessageSaver
}


type ChatUseCase interface {
	WriteMessage(c lec.Context, msg *entity.Message) e.Error
	GetChat(c lec.Context, userId uint64) ([]*entity.Message, e.Error)
	GetSupportChats(c lec.Context, id uint64) ([]uint64, e.Error)
	ReadChat(c lec.Context, userId uint64, readerId uint64) e.Error
}

func New(mw MessageSaver) *Chat {
	return &Chat{
		mw: mw,
	}
}
func (c *Chat) WriteMessage(ctx lec.Context, msg *entity.Message) e.Error {
	msg.Timestamp = time.Now()

	supportId, err := c.mw.GetSupportIdFromUser(ctx, msg.UserId)
	if err != nil {
		return err
	}

	msg.SupportId = supportId

	ctx.Logger().Debug("Support id", sl.Uint64Attr("id", supportId))

	if msg.SupportId == 0 {
		newSupport, err := c.mw.ChooseSupport(ctx)
		if err != nil {
			return err
		}

		msg.SupportId = newSupport.SupportId

		ctx.Logger().Debug("Support choose", sl.Uint64Attr("id", newSupport.SupportId))

		err = c.mw.SetSupportToChat(ctx, entity.SupportToChat{
			SupportId: msg.SupportId,
			UserId:    msg.UserId,
		})
		if err != nil {
			return err
		}
	}

	return c.mw.WriteMessage(ctx, *msg)
}

func (c *Chat) GetChat(ctx lec.Context, userId uint64) ([]*entity.Message, e.Error) {
	return c.mw.GetChat(ctx, userId)
}

func (c *Chat) GetSupportChats(ctx lec.Context, supportId uint64) ([]uint64, e.Error) {
	return c.mw.GetSupportChats(ctx, supportId)
}

func (c *Chat) ReadChat(ctx lec.Context, userId, readerId uint64) e.Error {
	return c.mw.ReadChat(ctx, userId, readerId)
}
