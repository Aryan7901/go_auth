package utils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func HashString(password string) (string, error) {
	var passwordBytes = []byte(password)
	password_hash_bytes, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	return string(password_hash_bytes), err
}

type HashResult struct {
	Hash  string
	Error error
}

func HashStringAsync(password string) chan HashResult {
	ch := make(chan HashResult)
	go func(password string, ch chan HashResult) {
		hashed_password, err := HashString(password)
		if err != nil {
			ch <- HashResult{"", errors.New("does not compute")}
		} else {
			ch <- HashResult{hashed_password, nil}
		}
	}(password, ch)
	return ch
}
