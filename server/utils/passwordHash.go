package utils

import(
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string)(string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", errors.New("error occurred while creating a Hash")
	}

	return string(hashedPassword), nil
}

func VerifyPassword(hashPass, password string) error{
	err := bcrypt.CompareHashAndPassword([]byte(hashPass), []byte(password))
	if err != nil{
		return errors.New("The " + password + " and " + hashPass + " doesn't match")
	}
	return nil
}