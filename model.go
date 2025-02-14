package deepinfra

import (
	"github.com/fumiama/deepinfra/model"
)

type Model interface {
	model.Inputer
	model.Outputer
}
