package main

import (
	"errors"
	"fmt"
)

var (
	ERR_NO_REDIRECTIONS      = errors.New("No active redirections.")
	ERR_NO_ADDRESS_SPECIFIED = errors.New("No remote address specified.")
)

func NewError(format string, vals ...interface{}) error {
	return errors.New(fmt.Sprintf(format, vals...))
}
