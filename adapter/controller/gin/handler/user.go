package handler

import (
	"go-api-arch-clean-template/api"
	"net/http"

	"github.com/gin-gonic/gin"

	"go-api-arch-clean-template/adapter/controller/gin/presenter"
)

type UserHandler struct {
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

func (a *UserHandler) GetUserById(c *gin.Context, ID int) {
	c.JSON(http.StatusOK, &presenter.UserResponse{
		ApiVersion: api.Version,
		Data: presenter.User{
			Kind: "user",
			Id:   1,
			Name: "jun",
		},
	})
}
