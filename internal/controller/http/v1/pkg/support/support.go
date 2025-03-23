package support

import (
	"fmt"
	"net/http"
	"strconv"

	conv "github.com/coffee/support/internal/controller/http/v1/converter"
	"github.com/coffee/support/internal/controller/http/v1/dto"
	"github.com/coffee/support/internal/controller/http/v1/pkg/support/ws"
	resp "github.com/coffee/support/internal/controller/response"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/gosuit/e"
	"github.com/gosuit/gins"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Support struct {
	support SupportUseCase
	chat    ChatUseCase
	wsHub   *ws.Hub
}

func New(support SupportUseCase, chat ChatUseCase) *Support {
	hub := ws.NewHub()
	go hub.Run()

	return &Support{
		support: support,
		chat:    chat,
		wsHub:   hub,
	}
}

// @Summary Connect to support
// @Description Upgrade connection to websockets
// @Tags Support
// @Accept json
// @Produce json
// @Security Bearer
// @Failure 401 {object} resp.JsonError "Authorization header wasn`t found, Token is not bearer"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /support/ws [get]
func (s *Support) Join(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		resp.AbortErrMsg(c, internalErr.WithErr(err))
		return
	}

	userId := c.GetUint64("userId")

	client := &ws.Client{
		Conn:   conn,
		RoomId: userId,
		UserId: userId,
		Send:   make(chan *dto.WsMessage),
	}

	s.wsHub.Register <- client

	go client.Write()
	go client.Read(s.chat, s.wsHub)
}

// @Summary Sup to support chat
// @Description Upgrade connection to websockets (available for users with role "SUPPORT")
// @Tags Support
// @Accept json
// @Produce json
// @Security Bearer
// @Param	id	path		int	true	"User ID"
// @Failure 401 {object} resp.JsonError "Authorization header wasn`t found, Token is not bearer"
// @Failure 403 {object} resp.JsonError "This resource is forbidden"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /support/s/chats/:id/ws [get]
func (s *Support) JoinSupport(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		resp.AbortErrMsg(c, internalErr.WithErr(err))
		return
	}

	supId := c.GetUint64("userId")

	userId, parseErr := strconv.ParseUint(c.Param("id"), 10, 64)
	if parseErr != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(parseErr))
		return
	}

	client := &ws.Client{
		Conn:   conn,
		RoomId: userId,
		UserId: supId,
		Send:   make(chan *dto.WsMessage),
	}

	s.wsHub.Register <- client

	go client.Write()
	go client.Read(s.chat, s.wsHub)
}

// @Summary History of chat
// @Description Returns list of messages of chat
// @Tags Support
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []dto.Message "Successful response with history"
// @Failure 401 {object} resp.JsonError "Authorization header wasn`t found, Token is not bearer"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /support/history [get]
func (s *Support) ChatHistory(c *gin.Context) {
	ctx := gins.GetCtx(c)
	userId := c.GetUint64("userId")

	chat, err := s.chat.GetChat(ctx, userId)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.MsgToChat(chat)

	c.JSON(okStatus, result)
}

// @Summary History of chat (sup)
// @Description Returns list of messages by userId (available for users with role "SUPPORT")
// @Tags Support
// @Accept json
// @Produce json
// @Security Bearer
// @Param	id	path		int	true	"User ID"
// @Success 200 {object} []dto.Message "Successful response with history"
// @Failure 401 {object} resp.JsonError "Authorization header wasn`t found, Token is not bearer"
// @Failure 403 {object} resp.JsonError "This resource is forbidden"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /support/s/chats/:id/history [get]
func (s *Support) SupChatHistory(c *gin.Context) {
	ctx := gins.GetCtx(c)

	userId, parseErr := strconv.ParseUint(c.Param("id"), 10, 64)
	if parseErr != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(parseErr))
		return
	}

	chat, err := s.chat.GetChat(ctx, userId)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	result := conv.MsgToChat(chat)

	c.JSON(okStatus, result)
}

// @Summary Support chats
// @Description Returns list of chats support must to answer (available for users with role "SUPPORT")
// @Tags Support
// @Accept json
// @Produce json
// @Security Bearer
// @Success 200 {object} []int "Successful response with chats"
// @Failure 401 {object} resp.JsonError "Authorization header wasn`t found, Token is not bearer"
// @Failure 403 {object} resp.JsonError "This resource is forbidden"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /support/s/chats/ [get]
func (s *Support) SupportChats(c *gin.Context) {
	ctx := gins.GetCtx(c)
	userId := c.GetUint64("userId")

	chats, err := s.chat.GetSupportChats(ctx, userId)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(okStatus, chats)
}

// @Summary Add support
// @Description Register new SUPPORT user (available for users with role "SUPER_ADMIN")
// @Tags Support
// @Accept json
// @Produce json
// @Security Bearer
// @Param	id	query		int	true	"User ID"
// @Failure 401 {object} resp.JsonError "Authorization header wasn`t found, Token is not bearer"
// @Failure 403 {object} resp.JsonError "This resource is forbidden"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /support/admin/sup/add [post]
func (s *Support) AddSupport(c *gin.Context) {
	ctx := gins.GetCtx(c)

	supId, parseErr := strconv.ParseUint(c.Query("id"), 10, 64)
	if parseErr != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(parseErr))
		return
	}

	err := s.support.AddSupport(ctx, supId)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(okStatus, addMsg)
}

// @Summary Remove support
// @Description Unregister SUPPORT user (available for users with role "SUPER_ADMIN")
// @Tags Support
// @Accept json
// @Produce json
// @Security Bearer
// @Param	id	query		int	true	"User ID"
// @Failure 401 {object} resp.JsonError "Authorization header wasn`t found, Token is not bearer"
// @Failure 403 {object} resp.JsonError "This resource is forbidden"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /support/admin/sup/remove [delete]
func (s *Support) RemoveSupport(c *gin.Context) {
	ctx := gins.GetCtx(c)

	supId, parseErr := strconv.ParseUint(c.Query("id"), 10, 64)
	if parseErr != nil {
		resp.AbortErrMsg(c, badReqErr.WithErr(parseErr))
		return
	}

	err := s.support.RemoveSupport(ctx, supId)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	c.JSON(okStatus, rmMsg)
}

// @Summary Get all supports
// @Description Unregister SUPPORT user (available for users with role "SUPER_ADMIN")
// @Tags Support
// @Accept json
// @Produce json
// @Security Bearer
// @Param limit query uint false "Limit"
// @Param offset query uint false "Offset"s
// @Success 200  {object} []dto.Support "List of supportss"
// @Failure 400 {object} resp.JsonError "Invalid request params"
// @Failure 401 {object} resp.JsonError "Authorization header wasn`t found, Token is not bearer"
// @Failure 403 {object} resp.JsonError "This resource is forbidden"
// @Failure 500 {object} resp.JsonError "Something going wrong..."
// @Router /support/admin/sup/all [get]
func (s *Support) GetAllSupports(c *gin.Context) {
	ctx := gins.GetCtx(c)

	limit, offset, err := getLimitAndOffsetParams(c)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	supports, err := s.support.GetAllSupports(ctx, offset, limit)
	if err != nil {
		resp.AbortErrMsg(c, err)
		return
	}

	res := conv.ListSupportsToDto(supports)

	c.JSON(http.StatusOK, res)
}

func getLimitAndOffsetParams(c *gin.Context) (limit int64, offset int64, er e.Error) {
	var errF = func(err error) e.Error {
		return e.E(err).WithCode(e.BadInput).WithMessage("Invalid request params")
	}

	var l, ofst int
	var err error

	limitPar := c.Query("limit")
	if len(limitPar) == 0 {
		limit = 10
	} else {
		l, err = strconv.Atoi(limitPar)
		if err != nil {
			return 0, 0, errF(err)
		}
	}

	offsetPar := c.Query("offset")
	if len(offsetPar) == 0 {
		offset = 0
	} else {
		ofst, err = strconv.Atoi(offsetPar)
		if err != nil {
			return 0, 0, errF(err)
		}
	}

	if (ofst < 0 || ofst > 10000) || (limit < 0 || limit > 10000) {
		return 0, 0, errF(fmt.Errorf("invalid variable of limit or offset"))
	}

	return int64(l), int64(ofst), nil
}
