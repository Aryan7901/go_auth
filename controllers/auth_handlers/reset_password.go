package auth_handler

import (
	AppErrors "auth/helpers/errors"
	utils "auth/helpers/utils"
	validate "auth/helpers/validate"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (conf *AuthConfig) ResetPassword(c echo.Context) error {
	payload := ResetPassword{}
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	if payload.Password == payload.NewPassword {
		return c.String(http.StatusBadRequest, "New Password cannot be same as old password")
	}
	errs := validate.ValidateStruct(&payload)
	if errs != nil {
		return c.JSON(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	user_details := Login{}
	err2 := conf.Db.Get(&user_details, getEmailPassword, payload.Email)
	if err2 != nil {
		return c.String(http.StatusBadRequest, err2.Error())
	}
	err3 := bcrypt.CompareHashAndPassword([]byte(user_details.Password), []byte(payload.Password))
	if err3 != nil {
		return c.String(http.StatusBadRequest, "Invalid Email/Password")
	}
	hashed_new_password, err4 := utils.HashString(payload.NewPassword)
	if err4 != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}
	_, err5 := conf.Db.Exec( /* sql */ `UPDATE user_auth SET password_hash=$1 WHERE email=$2`, hashed_new_password, payload.Email)
	if err5 != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}
	return c.String(http.StatusOK, "Password Reset Successful")
}
