package deepinfra

import (
	"bytes"
	"io"
)

type Model interface {
	Body() *bytes.Buffer
	Parse(io.Reader) error
	Output() string
}
