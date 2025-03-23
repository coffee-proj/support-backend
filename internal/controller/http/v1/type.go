package v1

import (
	types "github.com/coffee/support/internal/entity/type"
	"github.com/gin-gonic/gin"
	"github.com/gosuit/lec"
)

type SupportHandler interface {
	Join(c *gin.Context)
	JoinSupport(c *gin.Context)
	ChatHistory(c *gin.Context)
	SupChatHistory(c *gin.Context)
	SupportChats(c *gin.Context)
	AddSupport(c *gin.Context)
	RemoveSupport(c *gin.Context)
	GetAllSupports(c *gin.Context)
}

type Middleware interface {
	CheckAccess(roles ...types.Role) gin.HandlerFunc
	InitLogger(ctx lec.Context) gin.HandlerFunc
}
