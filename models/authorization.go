package models

import (
	"errors"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func CheckUserType(ctx *gin.Context, role string) (err error) {
	userType := ctx.GetString("userType")
	err = nil
	if userType != role {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	return err
}

func MatchUserTypeToUid(ctx *gin.Context, userId string) (err error) {
	userType := ctx.GetString("userType")
	uid := ctx.GetString("uid")
	err = nil
	if userType == "USER" && uid != userId {
		err = errors.New("unauthorized to access this resource")
		return err
	}
	err = CheckUserType(ctx, userType)
	return err
}

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func VerifyPassword(userPass, providedPass string) (bool, error) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPass), []byte(userPass))
	if err != nil {
		err = fmt.Errorf("user password is incorrect")
		return false, err
	}
	return true, nil
}
