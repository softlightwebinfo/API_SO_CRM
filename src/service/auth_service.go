package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"so-crm/src/entity"
)

type authService struct {
	Context *gin.Context
}

func NewAuthService(context *gin.Context) *authService {
	return &authService{
		Context: context,
	}
}
func AuthServiceLogin(context *gin.Context) entity.User {
	auth := NewAuthService(context)
	return auth.User()
}
func (auth *authService) User() entity.User {
	get, exists := auth.Context.Get("user")
	if !exists {
		auth.Context.JSON(http.StatusForbidden, nil)
	}
	return get.(entity.User)
}
