package models

import (
	"golang-jwt/database"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type User struct {
	ID            primitive.ObjectID `bson: "_id"`
	First_name    *string            `json: "first_name" validate:"required, min=2, max=100"`
	Last_name     *string            `json: "last_name" validate:"required, min=2, max=100"`
	Password      *string            `json:"password" validade: required, min=6`
	Email         *string            `json:"email"`
	Phone         *string            `json:"phone"`
	Token         *string            `json:"token"`
	User_Type     *string            `json:"user_type"`
	Refresh_Token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_id       string             `json:"user_id"`
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func Login() {

}

func GetUsers() {

}

func GetUser() {

}

func HashPassword() {

}

func VerifyPassword() {

}

func Signup() {

}
