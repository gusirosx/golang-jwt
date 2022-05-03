package models

import (
	"context"
	"fmt"
	"golang-jwt/entity"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

func Login(email, password *string) (entity.User, error) {

	var user entity.User // Found User
	var queryCtx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	err := collection.FindOne(queryCtx, bson.M{"email": email}).Decode(&user)
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

	if err := UpdateAllTOkens(user); err != nil {
		return entity.User{}, err
	}
	err = collection.FindOne(queryCtx, bson.M{"uid": user.UID}).Decode(&user)
	if err != nil {
		return entity.User{}, err
	}
	return user, nil
}
