package models

import (
	"context"
	"fmt"
	"golang-jwt/database"
	"golang-jwt/entity"
	"golang-jwt/helpers"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user entity.User
		defer cancel()
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		count, err := userCollection.CountDocuments(c, bson.M{"email": user.Email})
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
		}

		count, err = userCollection.CountDocuments(c, bson.M{"phone": user.Phone})
		if err != nil {
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone"})
		}

		if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "this email or phone number already exists"})
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refresh_token, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_Type, *&user.User_id)
		user.Token = &token
		user.Refresh_Token = &refresh_token
		resultInsertionNumber, err := userCollection.InsertOne(c, user)
		if err != nil {
			msg := fmt.Sprintf("user item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		ctx.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user entity.User
		var foundUser entity.User
		if err := ctx.BindJSON(&user); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(c, bson.M{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"email": "email or password is incorrect"})
			return
		}
		passwordIsValid, msg := VerifyPassword(*user.Password, *&foundUser.Password)
		defer cancel()
		//!!Pesquisar sobre o defer cancel() debaixo de cada if
	}
}

func GetUsers() {

}

func GetUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID := ctx.Param("user_id")

		if err := helpers.MatchUserTypeToUid(ctx, userID); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var c, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var user entity.User
		if err := userCollection.FindOne(c, bson.M{"user_id": userID}).Decode(&user); err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ctx.JSON(http.StatusOK, user)
	}
}

func HashPassword() {

}

func VerifyPassword(userPass, providedPass string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPass), []byte(userPass))
	check := true
	msg := ""
	if err != nil {
		msg = fmt.Sprintf("email of password is incorrect")
		check = false
	}
	return check, msg
}
