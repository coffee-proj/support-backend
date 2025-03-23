package storage

import (
	"github.com/coffee/support/internal/usecase/storage/chat"
	"github.com/coffee/support/internal/usecase/storage/support"
	"github.com/gosuit/e"
	"github.com/gosuit/lec"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Storage struct {
	mg *mongo.Database
	*chat.MessageStorage
	*support.SupportStorage
}

func New(mongo *mongo.Database) *Storage {
	return &Storage{
		mg:             mongo,
		MessageStorage: chat.New(mongo),
		SupportStorage: support.New(mongo),
	}
}

func (s *Storage) Close(ctx lec.Context) e.Error {
	err := s.mg.Client().Disconnect(ctx)
	if err != nil {
		return e.E(err).WithMessage("Failed to disconnect mongo")
	}

	return nil
}
