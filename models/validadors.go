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

	count, err := collection.CountDocuments(ctx, bson.M{"email": email})
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

	count, err := collection.CountDocuments(ctx, bson.M{"phone": phone})
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

// // End of validators
