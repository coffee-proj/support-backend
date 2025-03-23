package auth

import e "github.com/gosuit/e"

type Auth struct {
	jwt *Jwt
}

func New(opts *JwtOptions) *Auth {
	return &Auth{
		jwt: NewJwt(opts),
	}
}

func (a *Auth) ValidateToken(token string) (*Claims, e.Error) {
	return a.jwt.ValidateToken(token)
}

func (a *Auth) GetAudience() []string {
	return a.jwt.audience
}
