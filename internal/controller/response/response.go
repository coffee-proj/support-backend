package resp

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gosuit/e"
	"github.com/gosuit/gins"
	"github.com/gosuit/sl"
)

type Message struct {
	Message string `json:"message"`
}

func NewMessage(msg string) *Message {
	return &Message{
		Message: msg,
	}
}

// JsonError use only for doc and represent e.JsonError
type JsonError struct {
	Error string `json:"error"`
}

func AbortErrMsg(c *gin.Context, err e.Error) {
	ctx := gins.GetCtx(c)
	log := ctx.Logger()

	if err.GetCode() == e.Internal {
		log.Error("Something going wrong...", err.SlErr())
	} else {
		log.Info("Invalid input data", err.SlErr())
	}

	c.AbortWithStatusJSON(
		err.ToHttpCode(),
		err.ToJson(),
	)
}

func WsErr(conn *websocket.Conn, err e.Error) {
	log := sl.Default()

	if err.GetCode() == e.Internal {
		log.Error("Something going wrong...", err.SlErr())
	} else {
		log.Info("Invalid input data", err.SlErr())
	}

	conn.WriteJSON(err.ToJson())
}
