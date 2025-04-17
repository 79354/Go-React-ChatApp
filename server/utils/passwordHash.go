package utils

import(
	"errors"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string)(string, error){
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil{
		return "", errors.New("Error occurred while creating a Hash")
	}

	return string(hashedPassword), nil
}

func VerifyPassword(hashPass, password sting) error{
	err := bcrypt.CompareHashAndPassword(hashPass, password)
	if err != nil{
		return errors.New("The " + password + " and " + hashPass + " doesn't match")
	}
	return nil
}