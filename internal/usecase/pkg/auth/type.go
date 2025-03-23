package auth

import (
	types "github.com/coffee/support/internal/entity/type"
	"github.com/golang-jwt/jwt/v5"
)

type JwtOptions struct {
	Audience []string `yaml:"audience" env:"JWT_AUDIENCE"`
	Issuer   string   `yaml:"issuer"   env:"JWT_ISSUER"`
	Key      string   `env:"JWT_SECRET"`
}

type Claims struct {
	Id    uint64       `json:"id"`
	Roles []types.Role `json:"roles"`
	jwt.RegisteredClaims
}
