package main

import (
	"errors"
)

var (
	InvalidInputErr  = errors.New("usage > main.go <URL> [-H key:value]")
	InvalidMethodErr = errors.New("invalid method, please use an appropriate HTTP method.")

	bodyNotAllowederr = errors.New("body not allowed with method.")
)
