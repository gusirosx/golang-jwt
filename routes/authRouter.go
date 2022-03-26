package routes

import (
	"golang-jwt/models"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.Engine) {
	routes.POST("users/signup", models.Signup())
	routes.POST("users/login", models.Login())
}
