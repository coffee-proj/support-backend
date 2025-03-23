package app

import (
	"fmt"

	"github.com/coffee/support/internal/usecase"
	"github.com/coffee/support/internal/usecase/storage"

	v1 "github.com/coffee/support/internal/controller"
	"github.com/gosuit/e"
	"github.com/gosuit/httper"
	"github.com/gosuit/lec"
	"github.com/gosuit/mongo"
	"github.com/gosuit/sl"
)

type App struct {
	controller *v1.Controller
	usecase    *usecase.UseCase
	storage    *storage.Storage
	ctx        lec.Context
	server     *httper.Server
}

func New() *App {
	cfg, err := getAppConfig()
	if err != nil {
		panic(fmt.Errorf("can`t get application config. Error: %s", err.Error()))
	}

	logger := sl.New(&cfg.Logger)
	ctx := lec.New(logger)

	mg, err := mongo.New(ctx, &cfg.Mongo)
	if err != nil {
		panic(err)
	}
	logger.Info("Connect to mongo")

	app := &App{}

	app.ctx = ctx

	app.storage = storage.New(mg)

	app.usecase = usecase.New(app.storage, &cfg.Usecase)

	app.controller = v1.New(app.usecase, &cfg.Controller)

	app.server = httper.NewServer(&cfg.Server, app.controller.InitRoutes(ctx))

	return app
}

func (a *App) Run() {
	c := a.ctx

	log := c.Logger()

	a.server.Start()

	if err := a.shutdown(); err != nil {
		log.Error("Failed to shutdown app", sl.ErrAttr(err))
		return
	}

	log.Info("Application stopped successfully")
}

func (a *App) shutdown() e.Error {
	log := a.ctx.Logger()

	a.server.Shutdown(log.ToSlog())

	log.Info("Server stopped")

	err := a.storage.Close(a.ctx)

	if err != nil {
		log.Error("Failed to close storage", err.SlErr())
		return err
	}

	log.Info("Storage closed")

	return nil
}
