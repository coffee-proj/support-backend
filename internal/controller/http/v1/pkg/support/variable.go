package support

import (
	resp "github.com/coffee/support/internal/controller/response"
	e "github.com/gosuit/e"
	"github.com/gosuit/httper"
)

const (
	okStatus = httper.StatusOK
)

var (
	addMsg = resp.NewMessage("Support is added.")
	rmMsg  = resp.NewMessage("Support is removed.")
)

var (
	internalErr = e.New("Something went wrong", e.Internal)
	badReqErr = e.New("Invalid input.", e.BadInput)
)
