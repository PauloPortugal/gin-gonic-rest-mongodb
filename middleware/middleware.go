package middleware

import (
	"fmt"
	"net/http"

	"github.com/PauloPortugal/gin-gonic-rest-mongodb/model"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func AuthMiddleware(cfg *viper.Viper) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		auth := ctx.GetHeader("Authorization")
		claims := &model.Claims{}

		tkn, err := jwt.ParseWithClaims(auth, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(fmt.Sprint(cfg.Get("auth.secret"))), nil
		})
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		if tkn == nil || !tkn.Valid {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		ctx.Next()
	}
}
