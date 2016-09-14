package bing

import (
	"github.com/kataras/iris"
)

type myjson struct {
	Name string `json:"name"`
}

//Handler form "/"(index) route
func Index(c *iris.Context) {
	c.Render("bing/index.html", map[string]interface{}{
		"title":   "BING INDEX",
		"content": "BING CONTENT",
	})
}

func Piccy(c *iris.Context) {
	c.Render("bing/index.html", map[string]interface{}{
		"title":   "BING PICTURE",
		"content": "BING PICTURE CONTENT",
	})
}

func Json(c *iris.Context) {
	c.JSON(iris.StatusOK, myjson{Name: "iris"})
}