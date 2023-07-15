package auth_handler

import (
	AppErrors "auth/helpers/errors"
	utils "auth/helpers/utils"
	validate "auth/helpers/validate"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (conf *AuthConfig) ForgotPasswordReset(c echo.Context) error {
	payload := ForgotPassword{}
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	errs := validate.ValidateStruct(&payload)
	if errs != nil {
		return c.JSON(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	otp_details := ForgotPasswordReset{}
	err2 := conf.Db.Get(&otp_details, get_otp, payload.Email)
	if err2 != nil {
		return c.String(http.StatusBadRequest, "Could not find OTP")
	}
	if otp_details.OTPExpiry < time.Now().Unix() {
		return c.String(http.StatusBadRequest, "OTP Expired")
	}
	err3 := bcrypt.CompareHashAndPassword([]byte(otp_details.OTP), []byte(payload.OTP))
	if err3 != nil {
		return c.String(http.StatusBadRequest, "Invalid OTP")
	}
	hashed_new_password, err4 := utils.HashString(payload.Password)
	if err4 != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}
	_, err5 := conf.Db.Exec(update_password, hashed_new_password, payload.Email)
	if err5 != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}
	go conf.Db.Exec(delete_password, payload.Email)
	return c.String(http.StatusOK, "Password Reset Successful")

}
