package controller

import (
	"net/http"

	v1 "github.com/coffee/support/internal/controller/http/v1"
	"github.com/coffee/support/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/gosuit/lec"
)

type Config struct {
	V1   v1.Config `yaml:"v1"`
	Mode string    `yaml:"mode"`
}

type Controller struct {
	v1  *v1.Router
	cfg *Config
}

func New(uc *usecase.UseCase, cfg *Config) *Controller {
	return &Controller{
		v1:  v1.New(uc, &cfg.V1),
		cfg: cfg,
	}
}

func (c *Controller) InitRoutes(ctx lec.Context) *gin.Engine {
	setGinMode(c.cfg.Mode)

	router := gin.New()

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, "pong")
	})

	api := router.Group("/api")
	{
		c.v1.InitRoutes(ctx, api)
	}

	return router
}

func setGinMode(mode string) {
	switch mode {

	case "RELEASE":
		gin.SetMode(gin.ReleaseMode)

	case "TEST":
		gin.SetMode(gin.TestMode)

	case "DEBUG":
		gin.SetMode(gin.DebugMode)

	default:
		gin.SetMode(gin.DebugMode)

	}
}
