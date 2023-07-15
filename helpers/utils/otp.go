package utils

import (
	"math/rand"
	"time"
)

func GenerateOTP() string {
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	digits := "0123456789"
	OTP := ""
	for i := 0; i < 6; i++ {
		OTP += string(digits[random.Intn(len(digits))])
	}
	return OTP
}
