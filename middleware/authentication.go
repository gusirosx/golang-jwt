package middleware

import (
	"golang-jwt/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken == "" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "No authorization header provided"})
			return
		}

		claims, err := models.ValidadeToken(clientToken)
		if err != "" {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err})
			return
		}
		ctx.Set("email", claims.Email)
		ctx.Set("firstName", claims.FirstName)
		ctx.Set("lastName", claims.LastName)
		ctx.Set("uid", claims.UID)
		ctx.Set("userType", claims.UserType)
		ctx.Next()
	}
}
