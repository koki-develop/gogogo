package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koki-develop/gogogo/backend/pkg/controllers"
)

func New() *gin.Engine {
	r := gin.Default()

	r.GET("/h", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{"status": "ok"})
	})

	catscontroller := controllers.NewCatsController()
	r.GET("/v1/cats", catscontroller.FindAll)

	return r
}
