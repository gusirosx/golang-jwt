package models

import (
	"context"
	"fmt"
	"golang-jwt/database"
	"golang-jwt/entity"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func CreateUser(user entity.User) error {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// e-mail check
	if err := emailVerify(ctx, *user.Email); err != nil {
		return err
	}
	// phone check
	if err := phoneVerify(ctx, *user.Phone); err != nil {
		return err
	}

	count, err := userCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error occured while checking for the phone")
	} else if count > 0 {
		return fmt.Errorf("this phone has already been registered")
	}
	password := HashPassword(*user.Password)
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	id := primitive.NewObjectID()
	user.UID = id.Hex()
	token, refreshToken, _ := GenerateAllTokens(user)

	newuser := entity.User{
		ID:           id,
		UserName:     user.UserName,
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		Password:     &password,
		Email:        user.Email,
		Phone:        user.Phone,
		UserType:     user.UserType,
		Picture:      user.Picture,
		Token:        &token,
		RefreshToken: &refreshToken,
		Created:      time,
		Updated:      time,
		UID:          user.UID,
	}
	_, err = userCollection.InsertOne(ctx, newuser)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to create user")
	}

	return nil
}

func UpdateUser(id string, user entity.User) error {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// Declare a primitive ObjectID from a hexadecimal string
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	recU, err := GetUser(id)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error during user recovering")
	}
	// e-mail check
	if !strings.EqualFold(*user.Email, *recU.Email) {
		if err := emailVerify(ctx, *user.Email); err != nil {
			return err
		}
	}
	// phone check
	if !strings.EqualFold(*user.Phone, *recU.Phone) {
		if err := phoneVerify(ctx, *user.Phone); err != nil {
			return err
		}
	}

	var updatedUser primitive.D
	password := HashPassword(*user.Password)
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatedUser = append(updatedUser,
		bson.E{Key: "userName", Value: user.UserName},
		bson.E{Key: "firstName", Value: user.FirstName},
		bson.E{Key: "lastName", Value: user.LastName},
		bson.E{Key: "password", Value: &password},
		bson.E{Key: "email", Value: user.Email},
		bson.E{Key: "phone", Value: user.Phone},
		bson.E{Key: "userType", Value: user.UserType},
		bson.E{Key: "updated", Value: time})
	opt := options.Update().SetUpsert(true)
	update := bson.D{{Key: "$set", Value: updatedUser}}
	_, err = userCollection.UpdateByID(ctx, idPrimitive, update, opt)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to update user")
	}

	if err := UpdateAllTOkens(*recU.Token, *recU.RefreshToken, recU.UID); err != nil {
		return fmt.Errorf("unable to update the user token")
	}
	return nil
}

func DeleteUser(id string) error {
	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Declare a primitive ObjectID from a hexadecimal string
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	// Call the DeleteOne() method by passing BSON
	res, err := userCollection.DeleteOne(queryCtx, bson.M{"_id": idPrimitive})
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to delete user")
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("there is no such user for be deleted")
	}
	return nil
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
	log.Println(UID)
	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	if err = userCollection.FindOne(queryCtx, bson.M{"uid": UID}).Decode(&user); err != nil {
		return
	}

	log.Println(user.UID)

	return
}

// CreateClaims
// ReadClaims
// UpdateClaims
// DeleteClaims
