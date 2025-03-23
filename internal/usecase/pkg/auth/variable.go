package auth

import e "github.com/gosuit/e"

var (
	unauthErr = e.New("Token is invalid", e.Unauthorize)
)
