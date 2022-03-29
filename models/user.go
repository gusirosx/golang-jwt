package models

import (
	"context"
	"fmt"
	"golang-jwt/database"
	"golang-jwt/entity"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func Signup() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
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
		count, err := userCollection.CountDocuments(queryCtx, bson.M{"email": user.Email})
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		} else if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "this email has already been registered"})
			return
		}

		count, err = userCollection.CountDocuments(queryCtx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Println(err.Error())
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the phone"})
			return
		} else if count > 0 {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "this phone has already been registered"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		token, refreshToken, _ := GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)

		user.Created_at = time
		user.Updated_at = time
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()

		user.Token = &token
		user.Refresh_token = &refreshToken
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := "user item was not created"
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		ctx.JSON(http.StatusOK, resultInsertionNumber)
	}
}

// type User struct {
// 	ID            primitive.ObjectID `bson:"_id"`
// 	First_name    *string            `json:"first_name" validate:"required,min=2,max=100"`
// 	Last_name     *string            `json:"last_name" validate:"required,min=2,max=100"`
// 	Password      *string            `json:"Password" validate:"required,min=6"`
// 	Email         *string            `json:"email" validate:"email,required"`
// 	Phone         *string            `json:"phone" validate:"required"`
// 	Token         *string            `json:"token"`
// 	User_type     *string            `json:"user_type" validate:"required,eq=ADMIN|eq=USER"`
// 	Refresh_token *string            `json:"refresh_token"`
// 	Created_at    time.Time          `json:"created_at"`
// 	Updated_at    time.Time          `json:"updated_at"`
// 	User_id       string             `json:"user_id"`
// }

// {
// 	"First_name": "Gustavo",
// 	"Last_name": "Rodrigues",
// 	"Password": "123456",
// 	"Email":"gsr3@test.com",
// 	"Phone":"+5534900000002",
// 	"User_type":"ADMIN"
// }

// verificar se a chave secreta está sendo utilizada em algum lugar
func Login(email, password *string) (entity.User, error) {

	var user entity.User // Found User
	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := userCollection.FindOne(queryCtx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		return entity.User{}, fmt.Errorf("email is incorrect")
	}
	passwordIsValid, err := VerifyPassword(*password, *user.Password)
	if !passwordIsValid {
		return entity.User{}, err
	}
	if user.Email == nil {
		return entity.User{}, fmt.Errorf("user not found")
	}

	token, refreshToken, err := GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
	if err != nil {
		return entity.User{}, fmt.Errorf("unable to generate the user token's")
	}
	if err := UpdateAllTOkens(token, refreshToken, user.User_id); err != nil {
		return entity.User{}, fmt.Errorf("unable to update the user token")
	}
	err = userCollection.FindOne(queryCtx, bson.M{"user_id": user.User_id}).Decode(&user)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func GetUsers(ctx *gin.Context) (response *mongo.Cursor, err error) {
	recordPerPage, err := strconv.Atoi(ctx.Query("recordPerPage"))
	if err != nil || recordPerPage < 1 {
		recordPerPage = 10
	}

	startIndex, _ := strconv.Atoi(ctx.Query("startIndex"))
	matchStage := bson.D{{Key: "$match", Value: bson.D{{}}}}
	groupStage := bson.D{{Key: "$group", Value: bson.D{
		{Key: "_id", Value: bson.D{{Key: "_id", Value: "null"}}},
		{Key: "total_count", Value: bson.D{{Key: "$sum", Value: 1}}},
		{Key: "data", Value: bson.D{{Key: "$push", Value: "$$ROOT"}}}}}}
	projectStage := bson.D{
		{Key: "$project", Value: bson.D{
			{Key: "_id", Value: 0},
			{Key: "total_count", Value: 1},
			{Key: "user_items", Value: bson.D{
				{Key: "$slice", Value: []interface{}{"$data", startIndex, recordPerPage}}}}}}}

	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	response, err = userCollection.Aggregate(queryCtx, mongo.Pipeline{matchStage, groupStage, projectStage})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

func GetUser(UID string) (user entity.User, err error) {
	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	if err = userCollection.FindOne(queryCtx, bson.M{"user_id": UID}).Decode(&user); err != nil {
		return
	}
	return
}
