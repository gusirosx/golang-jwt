package routes

import (
	"golang-jwt/handlers"
	"golang-jwt/models"

	"github.com/gin-gonic/gin"
)

func AuthRoutes(routes *gin.Engine) {
	routes.POST("users/signup", handlers.Signup())
	routes.POST("users/login", models.Login())
}
