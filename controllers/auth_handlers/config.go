package auth_handler

import "auth/helpers/constants"

type AuthConfig struct {
	*constants.AppConfig
}

type SignUp struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8" db:"password_hash"`
	UserName string `json:"user_name" validate:"required"`
}
type UserCredentials struct {
	Email      string `json:"email" validate:"required,email"`
	Password   string `json:"password" validate:"required,min=8" db:"password_hash"`
	UserName   string `json:"user_name" validate:"required"`
	IsVerified bool   `json:"is_verified" db:"is_verified"`
}

type Login struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8" db:"password_hash"`
}

type ResetPassword struct {
	Email       string `json:"email" validate:"required,email"`
	Password    string `json:"password" validate:"required,min=8" db:"password_hash"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type Email struct {
	Email string `json:"email" validate:"required,email"`
}
type VerifyRegistration struct {
	IsVerified     bool   `db:"is_verified" validate:"required"`
	ActivationCode string `db:"activation_code" validate:"required"`
	Expiry         int64  `db:"activation_expiry" validate:"required"`
}
type VerificationStatus struct {
	IsVerified bool  `db:"is_verified" validate:"required"`
	Expiry     int64 `db:"activation_expiry" validate:"required"`
}

type ForgotPassword struct {
	Email    string `json:"email" validate:"required,email"`
	OTP      string `json:"otp" validate:"required"`
	Password string `json:"password" validate:"required,min=8"`
}
type ForgotPasswordReset struct {
	OTP       string `db:"otp" validate:"required"`
	OTPExpiry int64  `db:"otp_expiry" validate:"required"`
}
