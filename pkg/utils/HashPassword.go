package utils

import (
	"os"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

var cost int

func init(){
	if val,err := strconv.Atoi(os.Getenv("HASH_SALT")); err != nil {
		cost =bcrypt.DefaultCost
	} else {
		cost = val
	}
}

func GeneratePassword(password string)(string,error) {
	hash,err := bcrypt.GenerateFromPassword([]byte(password), cost)
	return string(hash),err
}

func CompareHashPassword(hash,password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
