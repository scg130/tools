package tools

import (
	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswd(passwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func CompareHashAndPasswd(passwd, target string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(target), []byte(passwd))
	if err != nil {
		return false
	}
	return true
}
