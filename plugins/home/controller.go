package home

import (
	//"fmt"
	"github.com/kataras/iris"
)

func Index(ctx *iris.Context) {
	ctx.Render("plugins/home/templates/index.html", map[string]interface{}{
		"page_title": "Iris", 
		"page_head": "", 
		"page_content": "CONTENT",
		"page_script": "",
		"parentLayout": "templates/theme/layouts/layout.semantic.html",
		"theme": "theme",
	})
}