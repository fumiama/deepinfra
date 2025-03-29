package model

import (
	"bytes"
	"io"
	"net/http"
)

type Inputer interface {
	Body() *bytes.Buffer
	Parse(io.Reader) error
}

type Outputer interface {
	Output() string
	OutputRaw() string
}

type Requester interface {
	API(api, key string) string       // API decorator
	Header(key string, h http.Header) // Header decorator
}

type MessageBuilder[T any] interface {
	System(prompt string) T
	User(prompt string) T
	Assistant(prompt string) T
}

type Protocol interface {
	Inputer
	Outputer
	Requester
	MessageBuilder[Protocol]
}
