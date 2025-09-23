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
	Output() Contents
	OutputRaw() Contents
}

type Requester interface {
	API(api, key string) string       // API decorator
	Header(key string, h http.Header) // Header decorator
}

type MessageBuilder[T any] interface {
	System(prompt string) T
	User(prompt ...Content) T
	Assistant(prompt ...Content) T
}

type Protocol interface {
	Inputer
	Outputer
	Requester
	MessageBuilder[Protocol]
}
