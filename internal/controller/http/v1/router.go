package v1

import (
	"github.com/coffee/support/internal/controller/http/v1/middleware"
	"github.com/coffee/support/internal/controller/http/v1/pkg/support"
	"github.com/coffee/support/internal/usecase"
	"github.com/coffee/support/pkg/swagger"
	"github.com/gin-gonic/gin"
	"github.com/gosuit/lec"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Config struct {
	Swagger swagger.SwaggerSpec `yaml:"swagger"`
}

type Router struct {
	support SupportHandler
	mid     Middleware
}

func New(uc *usecase.UseCase, cfg *Config) *Router {
	swagger.SetSwaggerConfig(cfg.Swagger)

	return &Router{
		support: support.New(uc.Support, uc.Chat),
		mid:     middleware.New(uc.Auth),
	}
}

func (r *Router) InitRoutes(ctx lec.Context, h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/v1")
	{
		router.Use(r.mid.InitLogger(ctx))

		r.initSupportRoutes(router)
		r.initSwaggerRoute(router)
	}

	return router
}

func (r *Router) initSupportRoutes(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("/support")
	{
		router.GET("/history", r.mid.CheckAccess(), r.support.ChatHistory)
		router.GET("/ws", r.mid.CheckAccess(), r.support.Join)

		supports := router.Group("/s/chats")
		{
			supports.Use(r.mid.CheckAccess("SUPPORT"))

			supports.GET("/", r.support.SupportChats)
			supports.GET("/:id/history", r.support.SupChatHistory)
			supports.GET("/:id/ws", r.support.JoinSupport)
		}

		admin := router.Group("/admin")
		{
			admin.Use(r.mid.CheckAccess("SUPER_ADMIN"))

			admin.POST("/sup/add", r.mid.CheckAccess("SUPER_ADMIN"), r.support.AddSupport)
			admin.DELETE("/sup/remove", r.mid.CheckAccess("SUPER_ADMIN"), r.support.RemoveSupport)
			admin.GET("/sup/all", r.support.GetAllSupports)
		}
	}

	return router
}

func (r *Router) initSwaggerRoute(h *gin.RouterGroup) *gin.RouterGroup {
	router := h.Group("swagger")
	{
		router.GET("/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	return router
}
