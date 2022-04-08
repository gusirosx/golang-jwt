package entity

import (
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID           primitive.ObjectID `bson:"_id"`
	UserName     *string            `json:"userName" validate:"required,min=6,max=20"`
	FirstName    *string            `json:"firstName" validate:"required,min=2,max=100"`
	LastName     *string            `json:"lastName" validate:"required,min=2,max=100"`
	Password     *string            `json:"password" validate:"required,min=6"`
	Email        *string            `json:"email" validate:"email,required"`
	Phone        *string            `json:"phone" validate:"required"`
	UserType     *string            `json:"userType" validate:"required,eq=ADMIN|eq=USER"`
	Picture      *string            `json:"picture"`
	Token        *string            `json:"token"`
	RefreshToken *string            `json:"refreshToken"`
	Created      time.Time          `json:"created"`
	Updated      time.Time          `json:"updated"`
	UID          string             `json:"uid"`
}

type SignedDetails struct {
	UID       string `json:"uid"`
	UserName  string `json:"userName"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	UserType  string `json:"userType"`
	jwt.StandardClaims
}

// 	 "permissions": {
// 	  "group_permissions": null,
// 	  "user_permissions": []
// 	},
