package grifts

import (
	"github.com/riton/toutdoux/actions"

	"github.com/gobuffalo/buffalo"
)

func init() {
	buffalo.Grifts(actions.App())
}
