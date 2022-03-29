package models

import (
	"context"
	"fmt"
	"golang-jwt/database"
	"golang-jwt/entity"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func Signup(user entity.User) error {
	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	count, err := userCollection.CountDocuments(queryCtx, bson.M{"email": user.Email})
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error occured while checking for the email")
	} else if count > 0 {
		return fmt.Errorf("this email has already been registered")
	}
	count, err = userCollection.CountDocuments(queryCtx, bson.M{"phone": user.Phone})
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error occured while checking for the phone")
	} else if count > 0 {
		return fmt.Errorf("this phone has already been registered")
	}

	password := HashPassword(*user.Password)
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	token, refreshToken, _ := GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
	id := primitive.NewObjectID()
	newuser := entity.User{
		ID:            id,
		First_name:    user.First_name,
		Last_name:     user.Last_name,
		Password:      &password,
		Email:         user.Email,
		Phone:         user.Phone,
		Token:         &token,
		Refresh_token: &refreshToken,
		User_type:     user.User_type,
		Created_at:    time,
		Updated_at:    time,
		User_id:       id.Hex(),
	}

	_, err = userCollection.InsertOne(queryCtx, newuser)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to create user")
	}
	return nil
}

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
