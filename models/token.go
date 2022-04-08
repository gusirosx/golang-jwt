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

var SECRET_KEY string = os.Getenv("SECRET_KEY")

//var SecretKey []byte = []byte(os.Getenv("JWT_SECRET_KEY"))

//func GenerateAllTokens(email, firstName, lastName, userType, uid string) (signedToken string, signedRefreshToken string, err error) {

func GenerateAllTokens(user entity.User) (string, string, error) {
	claims := &entity.SignedDetails{
		UID:       user.UID,
		UserName:  *user.UserName,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     *user.Email,
		Phone:     *user.Phone,
		UserType:  *user.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshClaims := &entity.SignedDetails{
		UID:       user.UID,
		UserName:  *user.UserName,
		FirstName: *user.FirstName,
		LastName:  *user.LastName,
		Email:     *user.Email,
		Phone:     *user.Phone,
		UserType:  *user.UserType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(365)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		log.Println(err.Error())
		return "", "", err
	}
	return token, refreshToken, nil
}

func UpdateAllTOkens(signedToken, signedRefreshToken, userId string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()
	var updateToken primitive.D

	UpdatedAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateToken = append(updateToken,
		bson.E{Key: "token", Value: signedToken},
		bson.E{Key: "refreshToken", Value: signedRefreshToken},
		bson.E{Key: "updated", Value: UpdatedAt})
	opt := options.Update().SetUpsert(true)
	filter := bson.M{"uid": userId}
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
