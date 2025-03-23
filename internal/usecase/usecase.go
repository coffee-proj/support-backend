package usecase

import (
	"github.com/coffee/support/internal/usecase/pkg/auth"
	"github.com/coffee/support/internal/usecase/pkg/chat"
	"github.com/coffee/support/internal/usecase/pkg/support"
	"github.com/coffee/support/internal/usecase/storage"
)

type UseCase struct {
	*auth.Auth
	*chat.Chat
	*support.Support
}

type Config struct {
	Jwt auth.JwtOptions `yaml:"jwt"`
}

func New(storage *storage.Storage, cfg *Config) *UseCase {
	return &UseCase{
		Chat:    chat.New(storage.MessageStorage),
		Support: support.New(storage.SupportStorage, storage.MessageStorage),
		Auth:    auth.New(&cfg.Jwt),
	}
}
