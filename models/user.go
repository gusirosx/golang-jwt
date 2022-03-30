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
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")

func CreateUser(user entity.User) error {
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

func UpdateUser(id string, user entity.User) error {

	// Declare a primitive ObjectID from a hexadecimal string
	idPrimitive, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	// count, err := userCollection.CountDocuments(queryCtx, bson.M{"email": user.Email})
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return fmt.Errorf("error occured while checking for the email")
	// } else if count > 0 {
	// 	return fmt.Errorf("this email has already been registered")
	// }
	// count, err = userCollection.CountDocuments(queryCtx, bson.M{"phone": user.Phone})
	// if err != nil {
	// 	log.Println(err.Error())
	// 	return fmt.Errorf("error occured while checking for the phone")
	// } else if count > 0 {
	// 	return fmt.Errorf("this phone has already been registered")
	// }

	var updatedUser primitive.D
	password := HashPassword(*user.Password)
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	updatedUser = append(updatedUser,
		bson.E{Key: "first_name", Value: user.First_name},
		bson.E{Key: "last_name", Value: user.Last_name},
		bson.E{Key: "password", Value: &password},
		bson.E{Key: "email", Value: user.Email},
		bson.E{Key: "phone", Value: user.Phone},
		bson.E{Key: "user_type", Value: user.User_type},
		bson.E{Key: "updated_at", Value: time})
	opt := options.Update().SetUpsert(true)
	update := bson.D{{Key: "$set", Value: updatedUser}}
	_, err = userCollection.UpdateByID(queryCtx, idPrimitive, update, opt)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("unable to update user")
	}

	recU, err := GetUser(id)
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error during user recovering")
	}
	if err := UpdateAllTOkens(*recU.Token, *recU.Refresh_token, recU.User_id); err != nil {
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

// Handle PUT requests at /users
// Handle DELETE requests at /users
// CreateUser
// UpdateUser
// DeleteUser
