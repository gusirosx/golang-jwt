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

// Create an unexported global variable to hold the database connection pool.
var client *mongo.Client = database.MongoInstance()

// Create an unexported global variable to hold the collection connection pool.
var collection *mongo.Collection = database.OpenCollection(client, "user")

// Get all users from the DB by its id
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
	response, err = collection.Aggregate(queryCtx, mongo.Pipeline{matchStage, groupStage, projectStage})
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

// Get one user from the DB by its id
func GetUser(UID string) (entity.User, error) {
	var user entity.User
	// Get a primitive ObjectID from a hexadecimal string
	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// Call the FindOne() method by passing BSON
	if err := collection.FindOne(queryCtx, bson.M{"uid": UID}).Decode(&user); err != nil {
		return entity.User{}, err
	}

	return user, nil
}

// Create one user into DB
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

	count, err := collection.CountDocuments(ctx, bson.M{"phone": user.Phone})
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error occured while checking for the phone")
	} else if count > 0 {
		return fmt.Errorf("this phone has already been registered")
	}
	password := HashPassword(*user.Password)
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	// get a unique userID
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
	_, err = collection.InsertOne(ctx, newuser)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to create user")
	}

	return nil
}

// Update one user from the DB by its id
func UpdateUser(id string, user entity.User) error {

	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	// Declare a primitive ObjectID from a hexadecimal string
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	recoveredUser, err := GetUser(id)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error during user recovering")
	}
	// e-mail check
	if !strings.EqualFold(*user.Email, *recoveredUser.Email) {
		if err := emailVerify(ctx, *user.Email); err != nil {
			return err
		}
	}
	// phone check
	if !strings.EqualFold(*user.Phone, *recoveredUser.Phone) {
		if err := phoneVerify(ctx, *user.Phone); err != nil {
			return err
		}
	}

	var updatedUser primitive.D
	password := HashPassword(*user.Password)
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	// Update user info with the new token's
	token, refreshToken, err := GenerateAllTokens(user)
	if err != nil {
		return fmt.Errorf("unable to generate the user token's")
	}

	updatedUser = append(updatedUser,
		bson.E{Key: "userName", Value: user.UserName},
		bson.E{Key: "firstName", Value: user.FirstName},
		bson.E{Key: "lastName", Value: user.LastName},
		bson.E{Key: "password", Value: &password},
		bson.E{Key: "email", Value: user.Email},
		bson.E{Key: "phone", Value: user.Phone},
		bson.E{Key: "userType", Value: user.UserType},
		bson.E{Key: "token", Value: token},
		bson.E{Key: "refreshtoken", Value: refreshToken},
		bson.E{Key: "updated", Value: time})
	opt := options.Update().SetUpsert(true)
	update := bson.D{{Key: "$set", Value: updatedUser}}
	_, err = collection.UpdateByID(ctx, idPrimitive, update, opt)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to update user")
	}

	// filter := bson.M{"uid": user.UID}
	// update := bson.D{{Key: "$set", Value: updateToken}}
	// _, err = userCollection.UpdateOne(ctx, filter, update, opt)
	// if err != nil {
	// 	return fmt.Errorf("unable to update the user token's")
	// }

	return nil
}

// Delete one user from the DB by its id
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
	res, err := collection.DeleteOne(queryCtx, bson.M{"_id": idPrimitive})
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to delete user")
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("there is no such user for be deleted")
	}
	return nil
}
