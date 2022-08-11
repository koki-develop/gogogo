package router

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/koki-develop/gogogo/backend/pkg/controllers"
)

func New() *gin.Engine {
	r := gin.Default()

	r.Use(func(ctx *gin.Context) {
		if os.Getenv("IS_LOCAL") == "true" {
			ctx.Header("Access-Control-Allow-Origin", "*")
		} else {
			ctx.Header("Access-Control-Allow-Origin", "https://go55.dev")
		}
		ctx.Next()
	})

	r.GET("/h", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	{
		ctrl := controllers.NewCatsController()
		r.GET("/v1/cats", ctrl.FindAll)
	}

	return r
}
