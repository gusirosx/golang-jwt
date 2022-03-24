package routes

import (
	"golang-jwt/middleware"
	"golang-jwt/models"

	"github.com/gin-gonic/gin"
)

func UserRoutes(routes *gin.Engine) {
	routes.Use(middleware.Authenticate())
	routes.GET("users", models.GetUsers())
	routes.GET("users/:user_id", models.GetUser())
}
