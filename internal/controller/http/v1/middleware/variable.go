package middleware

import "github.com/gosuit/e"

const (
	bearerType = "Bearer"
)

var (
	foundErr     = e.New("Authorization header wasn`t found", e.Unauthorize)
	bearerErr    = e.New("Token is not bearer", e.Unauthorize)
	forbiddenErr = e.New("This resource is forbidden", e.Forbidden)
	// unauthErr    = e.New("Invalid token.", e.Unauthorize) // is unused
)
