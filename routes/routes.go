package routes

import (
	"golang-jwt/middleware"
	"golang-jwt/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	port = ":8080"
)

func RoutesSetup() {
	// Set Gin to production mode
	//gin.SetMode(gin.ReleaseMode)

	// Set up a http server
	router := gin.Default()

	// Initialize the routes
	initializeRoutes(router)

	// Run the http server
	if err := router.Run(port); err != nil {
		log.Fatalln("could not run server: ", err.Error())
	} else {
		log.Println("Server listening on port: ", port)
	}
}

func initializeRoutes(router *gin.Engine) {
	// Handle the index route
	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success": "Up and running..."})
	})
	// Handle the no route case
	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Page not found"})
	})
	//Group user related routes together
	userRoutes := router.Group("/users")
	AuthRoutes(userRoutes)
	UserRoutes(userRoutes)
}

func AuthRoutes(routes *gin.RouterGroup) {
	// Handle signup requests at /users/signup
	routes.POST("/signup", models.Signup())
	// Handle login requests at /users/login
	routes.POST("/login", models.Login())
}

func UserRoutes(routes *gin.RouterGroup) {
	routes.Use(middleware.Authenticate())
	// Handle GET users requests at /users
	routes.GET("", models.GetUsers())
	// Handle GET user requests at /users/ID
	routes.GET("/:user_id", models.GetUser())
}