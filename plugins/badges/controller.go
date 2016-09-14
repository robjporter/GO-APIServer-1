package badges

import (
	"io"
	//"fmt"
	"github.com/kataras/iris"
	"github.com/roporter/go-libs/go-badge"
)

func Draw(ctx *iris.Context) {
	subject := ctx.Param("subject")
	status := ctx.Param("status")
	color := ctx.Param("color")
	
	if color == "" {
		color = "green"
	}
	
	ctx.SetHeader("Content-Type","image/svg+xml")
	ctx.SetHeader("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.SetStatusCode(iris.StatusOK)
	var out io.Writer
	out = ctx.RequestCtx.Response.BodyWriter()

	badge.Render(subject, status, badge.ColorString(color), out)
}