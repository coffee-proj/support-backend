package ws

import (
	conv "github.com/coffee/support/internal/controller/http/v1/converter"
	"github.com/coffee/support/internal/controller/http/v1/dto"
	resp "github.com/coffee/support/internal/controller/response"
	types "github.com/coffee/support/internal/entity/type"
	"github.com/gorilla/websocket"
	"github.com/gosuit/lec"
	"github.com/gosuit/sl"
)

type Client struct {
	Conn   *websocket.Conn
	RoomId uint64
	UserId uint64
	Send   chan *dto.WsMessage
}

func (c *Client) Write() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		msg := <-c.Send

		c.Conn.WriteJSON(msg)
	}
}

func (c *Client) Read(uc ChatUseCase, h *Hub) {
	defer func() {
		h.Unregister <- c
		c.Conn.Close()
	}()

	for {
		var msg dto.WsMessage

		readErr := c.Conn.ReadJSON(&msg)
		if readErr != nil {
			sl.Default().Warn("Read error", readErr)
			// resp.WsErr(c.Conn, internalErr.WithErr(readErr))
			break
		}

		ctx := lec.New(sl.Default())

		switch msg.Type {

		case types.READ:
			err := uc.ReadChat(ctx, c.RoomId, c.UserId)
			if err != nil {
				resp.WsErr(c.Conn, err)
			} else {
				h.Broadcast <- &msg
			}
		case types.WRITE:
			msg.ChatId = c.RoomId
			msg.SenderId = c.UserId

			err := uc.WriteMessage(ctx, conv.WsToWrite(&msg))
			if err != nil {
				resp.WsErr(c.Conn, err)
			} else {
				h.Broadcast <- &msg
			}
		default:
			resp.WsErr(c.Conn, badTypeErr)

		}
	}
}
