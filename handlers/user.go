package handlers

import (
	"golang-jwt/models"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

func GetUsers() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := models.CheckUserType(ctx, "ADMIN"); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		response, err := models.GetUsers(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while listing user items"})
			return
		}

		var allusers []bson.M
		if err = response.All(ctx, &allusers); err != nil {
			log.Println(err.Error())
		}
		ctx.JSON(http.StatusOK, allusers[0])
	}
}
