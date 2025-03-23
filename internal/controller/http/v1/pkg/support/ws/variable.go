package ws

import e "github.com/gosuit/e"

var (
	internalErr = e.New("Something went wrong", e.Internal)
	badTypeErr  = e.New("Invalid message type.", e.BadInput)
)
