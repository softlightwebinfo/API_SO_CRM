package middleware

import (
	"encoding/json"
	"fmt"
	"net/http"
	"so-crm/src/entity"
	"so-crm/src/service"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func AuthorizeJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		user := entity.User{}
		const BearerSchema = "Bearer"
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BearerSchema):]

		token, err := service.JWTAuthService().ValidateToken(strings.Trim(tokenString, " "))
		if token.Valid {
			claims := token.Claims.(jwt.MapClaims)
			dbByte, _ := json.Marshal(claims["user"])

			err := json.Unmarshal(dbByte, &user)
			if err != nil {
				fmt.Println(err.Error())
				c.AbortWithStatus(http.StatusUnauthorized)
			}
			c.Set("user", user)
		} else {
			fmt.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
