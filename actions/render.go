package actions

import (
	"crypto/md5"
	"encoding/hex"
	"strconv"
	"strings"

	"github.com/riton/toutdoux/public"
	"github.com/riton/toutdoux/templates"

	"github.com/gobuffalo/buffalo/render"
)

var r *render.Engine

func init() {
	r = render.New(render.Options{
		// HTML layout to be used for all HTML requests:
		HTMLLayout: "application.plush.html",

		// fs.FS containing templates
		TemplatesFS: templates.FS(),

		// fs.FS containing assets
		AssetsFS: public.FS(),

		// Add template helpers here:
		Helpers: render.Helpers{
			// for non-bootstrap form helpers uncomment the lines
			// below and import "github.com/gobuffalo/helpers/forms"
			// forms.FormKey:     forms.Form,
			// forms.FormForKey:  forms.FormFor,
			"gravatarURL": gravatarURL,
		},
	})
}

func gravatarURL(email string, size int) string {
	cksum := md5.Sum([]byte(strings.ToLower(email)))
	return "https://www.gravatar.com/avatar/" + hex.EncodeToString(cksum[:]) + "s=" + strconv.Itoa(size)
}
