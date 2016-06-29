package main

import (
	"github.com/iris-contrib/middleware/logger"
	"github.com/iris-contrib/middleware/recovery"
	"github.com/kataras/iris"
	"os"
	"time"
)

func main() {
	iris.Use(recovery.New(os.Stderr))
	iris.Use(logger.New(iris.Logger))

	iris.Config.Render.Template.Directory = "./public"

	iris.Get("/", func(ctx *iris.Context) {
		ctx.ServeFile("./public/index.html", false)
	})

	iris.Get("/ppf", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, map[string]string{"sentence": time.Now().Format(time.RFC3339)})
	})

	iris.Static("/public", "./public", 1)
	iris.Listen(":8080")
}
