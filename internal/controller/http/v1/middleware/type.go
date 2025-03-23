package middleware

import (
	"github.com/coffee/support/internal/usecase/pkg/auth"
	"github.com/gosuit/e"
)

type AuthUseCase interface {
	ValidateToken(token string) (*auth.Claims, e.Error)
	GetAudience() []string
}
