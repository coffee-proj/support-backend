package auth

import (
	"github.com/golang-jwt/jwt/v5"
	e "github.com/gosuit/e"
)

type Jwt struct {
	audience []string
	issuer   string
	key      string
}

func NewJwt(options *JwtOptions) *Jwt {
	return &Jwt{
		audience: options.Audience,
		issuer:   options.Issuer,
		key:      options.Key,
	}
}

func (j *Jwt) ValidateToken(jwtString string) (*Claims, e.Error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(j.key), nil
	}

	token, err := jwt.ParseWithClaims(jwtString, &Claims{}, keyFunc)
	if err != nil {
		return nil, unauthErr.WithErr(err)
	}

	if !token.Valid {
		return nil, unauthErr.WithErr(err)
	}

	return token.Claims.(*Claims), nil
}
