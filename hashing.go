package pakarbibackend

import (
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

func HashPass(password string) (string, error) {
	bytess, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytess), err
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


