package model

import (
	"bytes"
	"io"
)

type Inputer interface {
	Body() *bytes.Buffer
	Parse(io.Reader) error
}

type Outputer interface {
	Output() string
	OutputRaw() string
}

type MessageBuilder[T any] interface {
	System(prompt string) T
	User(prompt string) T
	Assistant(prompt string) T
}

type Protocol[T any] interface {
	Inputer
	Outputer
	MessageBuilder[T]
}
