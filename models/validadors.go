package models

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
)

// Checks if the email is already registered in the database
func emailVerify(ctx context.Context, email string) error {

	count, err := userCollection.CountDocuments(ctx, bson.M{"email": email})
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error occured while checking for the email")
	}
	if count > 0 {
		return fmt.Errorf("this email has already been registered")
	}

	return nil
}

// Checks if the phone is already registered in the database
func phoneVerify(ctx context.Context, phone string) error {

	count, err := userCollection.CountDocuments(ctx, bson.M{"phone": phone})
	if err != nil {
		log.Println(err.Error())
		return fmt.Errorf("error occured while checking for the phone")
	}
	if count > 0 {
		return fmt.Errorf("this phone has already been registered")
	}

	return nil
}

func validateEmail(email string) error {
	if email == "" {
		return fmt.Errorf("email must be a non-empty string")
	}
	if parts := strings.Split(email, "@"); len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("malformed email string: %q", email)
	}
	return nil
}

func validatePhone(phone string) error {
	if phone == "" {
		return fmt.Errorf("phone number must be a non-empty string")
	}
	if !regexp.MustCompile(`\+.*[0-9A-Za-z]`).MatchString(phone) {
		return fmt.Errorf("phone number must be a valid, E.164 compliant identifier")
	}
	return nil
}

// strings.Compare()
// status, err := regexp.MatchString(k, strings.ToUpper(SensorName.String))
// if err != nil {
// 	log.Println(err.Error())
// }

//userCollection.DeleteOne()

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

// password := HashPassword(*user.Password)
// time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
// token, refreshToken, _ := GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
// id := primitive.NewObjectID()
// newuser := entity.User{
// 	ID:            id,
// 	First_name:    user.First_name,
// 	Last_name:     user.Last_name,
// 	Password:      &password,
// 	Email:         user.Email,
// 	Phone:         user.Phone,
// 	Token:         &token,
// 	Refresh_token: &refreshToken,
// 	User_type:     user.User_type,
// 	Created_at:    time,
// 	Updated_at:    time,
// 	User_id:       id.Hex(),
// }

// verificar se a chave secreta está sendo utilizada em algum lugar

// // Validators.

// func validateDisplayName(val string) error {
// 	if val == "" {
// 		return fmt.Errorf("display name must be a non-empty string")
// 	}
// 	return nil
// }

// func validatePhotoURL(val string) error {
// 	if val == "" {
// 		return fmt.Errorf("photo url must be a non-empty string")
// 	}
// 	return nil
// }

// func validatePassword(val string) error {
// 	if len(val) < 6 {
// 		return fmt.Errorf("password must be a string at least 6 characters long")
// 	}
// 	return nil
// }

// func validateUID(uid string) error {
// 	if uid == "" {
// 		return fmt.Errorf("uid must be a non-empty string")
// 	}
// 	if len(uid) > 128 {
// 		return fmt.Errorf("uid string must not be longer than 128 characters")
// 	}
// 	return nil
// }

// func validateProviderUserInfo(p *UserProvider) error {
// 	if p.UID == "" {
// 		return fmt.Errorf("user provider must specify a uid")
// 	}
// 	if p.ProviderID == "" {
// 		return fmt.Errorf("user provider must specify a provider ID")
// 	}
// 	return nil
// }

// func validateProvider(providerID string, providerUID string) error {
// 	if providerID == "" {
// 		return fmt.Errorf("providerID must be a non-empty string")
// 	} else if providerUID == "" {
// 		return fmt.Errorf("providerUID must be a non-empty string")
// 	}
// 	return nil
// }

// // End of validators

// userID, _ := ctx.Get("UID")
// entry := new(logEntry.LoggingEntry)
// payload := entity.Payload{Name: "GetUser", UID: userID.(string)}
// values := ctx.Request.URL.Query()
// var user *auth.UserRecord
// var err error
// if _, ok := values["email"]; ok {
// 	user, err = models.GetUserByEmail(ctx, values["email"][0])

// } else if _, ok := values["phoneNumber"]; ok {
// 	user, err = models.GetUserByPhone(ctx, "+"+values["phoneNumber"][0])

// } else {
// 	err = fmt.Errorf("invalid search parameters")
// 	payload.Message = err.Error()
// 	service.Logger.Log(entry.New(payload, ctx, logging.Error, http.StatusBadRequest))
// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	fmt.Println(err.Error())
// 	return
// }

// var user entity.User
// // Call BindJSON to bind the received JSON to user
// err := ctx.BindJSON(&user)
// if err != nil {
// 	err := fmt.Errorf("delete user: %v", err)
// 	payload.Message = err.Error()
// 	service.Logger.Log(entry.New(payload, ctx, logging.Error, http.StatusBadRequest))
// 	ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
// 	fmt.Println(err.Error())
// 	return
// }
// err = service.Client.DeleteUser(ctx, user.Uid)
// if err != nil {
// 	err := fmt.Errorf("delete user: %v", err)
// 	payload.Message = err.Error()
// 	service.Logger.Log(entry.New(payload, ctx, logging.Error, http.StatusInternalServerError))
// 	ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
// 		"error":  err.Error(),
// 		"detail": map[string]interface{}{"userID": user.Uid},
// 	})
// 	fmt.Println(err.Error())
// 	return
// } else {
// 	fmt.Printf("User successfully deleted: %v\n", user.Uid)
// 	ctx.JSON(http.StatusCreated, gin.H{
// 		"success": "User successfully deleted",
// 		"detail":  map[string]interface{}{"userID": user.Uid},
// 	})
// 	payload.Message = "success"
