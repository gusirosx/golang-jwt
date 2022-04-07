package models

import (
	"context"
	"golang-jwt/entity"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails2 struct {
	Email      string
	First_name string
	Last_name  string
	Uid        string
	User_type  string
	jwt.StandardClaims
}

var SECRET_KEY string = os.Getenv("SECRET_KEY")

//var SecretKey []byte = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateAllTokens(email, firstName, lastName, userType, uid string) (signedToken string, signedRefreshToken string, err error) {
	claims := &entity.SignedDetails{
		UID:       uid,
		UserName:  "",
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     "",
		UserType:  userType,
		Picture:   "",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &entity.SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * 1).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println(err.Error())
		return
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println(err.Error())
		return
	}
	return token, refreshToken, err
}

func UpdateAllTOkens(signedToken, signedRefreshToken, userId string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateToken primitive.D

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateToken = append(updateToken,
		bson.E{Key: "token", Value: signedToken},
		bson.E{Key: "refresh_token", Value: signedRefreshToken},
		bson.E{Key: "updated_at", Value: UpdatedAt})
	opt := options.Update().SetUpsert(true)
	filter := bson.M{"user_id": userId}
	update := bson.D{{Key: "$set", Value: updateToken}}
	_, err = userCollection.UpdateOne(ctx, filter, update, opt)
	if err != nil {
		log.Println(err.Error())
		return
	}
	return
}

func ValidadeToken(signedToken string) (claims *entity.SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&entity.SignedDetails{},
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)
	if err != nil {
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*entity.SignedDetails)
	if !ok {
		msg = "token is invalid"
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is expired"
	}
	return
}
