package converter

import (
	"github.com/coffee/support/internal/controller/http/v1/dto"
	"github.com/coffee/support/internal/entity"
)

func MsgToChat(msgs []*entity.Message) []*dto.Message {
	res := make([]*dto.Message, 0)

	for _, msg := range msgs {
		toAdd := &dto.Message{
			MessageId: msg.MessageId,
			UserId:    msg.UserId,
			SupportId: msg.SupportId,
			SenderId:  msg.SenderId,
			Content:   msg.Content,
			IsRead:    msg.IsRead,
			Timestamp: msg.Timestamp,
		}

		res = append(res, toAdd)
	}

	return res
}

func WsToWrite(msg *dto.WsMessage) *entity.Message {
	return &entity.Message{
		SenderId: msg.SenderId,
		Content:  msg.Content,
		UserId:   msg.ChatId,
	}
}

func ListSupportsToDto(sups []entity.Support) []dto.Support {
	res := make([]dto.Support, 0, len(sups))
	for i := range sups {
		res = append(res, SupportToDto(sups[i]))
	}

	return res
}

func SupportToDto(sup entity.Support) dto.Support {
	return dto.Support{
		SupportId: sup.SupportId,
		CountChat: sup.CountChat,
	}
}
