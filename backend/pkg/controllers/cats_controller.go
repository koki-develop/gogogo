package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
)

type CatsController struct{}

func NewCatsController() *CatsController {
	return &CatsController{}
}

func (ctrl *CatsController) FindAll(ctx *gin.Context) {
	cats := entities.Cats{
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
		{URL: "https://koki.me/images/profile.png"},
	}

	ctx.JSON(http.StatusOK, cats)
}
