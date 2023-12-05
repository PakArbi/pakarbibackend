package pakarbibackend

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(passwordhash string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(passwordhash), 14)
	return string(bytes), err
}

func CompareHashPass(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CheckPasswordHash(passwordhash, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(passwordhash))
	return err == nil
}

func CheckEmail(email string) bool {
	emailRegex := `^[a-zA-Z0-9._%+-]+@std.ulbi.ac.id$`
	match, _ := regexp.MatchString(emailRegex, email)
	return match
}
