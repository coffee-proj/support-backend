package middleware

type Middleware struct {
	auth AuthUseCase
}

func New(auth AuthUseCase) *Middleware {
	return &Middleware{
		auth,
	}
}
