package validator

import (
	"github.com/go-playground/validator/v10"
	"github.com/gosuit/e"
)

const (
	Password Arg = iota
)

var (
	lenErr = e.New("Bad string length", e.BadInput)
)

type Arg int

func setupArgs(validate *validator.Validate, args []Arg) error {
	for i := range len(args) {
		switch args[i] {
		default:
			return nil
		}
	}

	return nil
}
