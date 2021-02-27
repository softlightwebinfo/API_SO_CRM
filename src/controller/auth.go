package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"so-crm/src/dto"
	"so-crm/src/entity"
	"so-crm/src/models"
	"so-crm/src/service"
)

type LoginController struct {
	LoginService service.LoginService
	JWtService   service.JWTService
	User         models.AuthUser
	Db           *gorm.DB
}

func (controller *LoginController) Login(ctx *gin.Context) (string, *entity.User) {
	var credential dto.LoginCredentials
	err := ctx.BindJSON(&credential)
	if err != nil {
		return "no data found", nil
	}
	isUserAuthenticated, user := controller.LoginService.LoginUser(credential.Email, credential.Password)
	if isUserAuthenticated {
		return controller.JWtService.GenerateToken(*user, true), user
	}
	return "", nil
}

func (controller *LoginController) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		token, user := controller.Login(ctx)
		if token != "" {
			ctx.JSON(http.StatusOK, gin.H{
				"token": token,
				"user":  user,
			})
		} else {
			ctx.JSON(http.StatusUnauthorized, nil)
		}
	}
}
