package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/koki-develop/gogogo/backend/pkg/entities"
	"github.com/koki-develop/gogogo/backend/pkg/infrastructures/s3"
)

type CatsController struct{}

func NewCatsController() *CatsController {
	return &CatsController{}
}

func (ctrl *CatsController) FindAll(ctx *gin.Context) {
	s3cl := s3.New()

	data, err := s3cl.Download("gogogo-cats", "cats.json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
		return
	}

	var cats entities.Cats
	if err := json.NewDecoder(data).Decode(&cats); err != nil {
		ctx.JSON(http.StatusInternalServerError, map[string]string{"message": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, cats)
}
