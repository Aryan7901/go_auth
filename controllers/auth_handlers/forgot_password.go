package auth_handler

import (
	AppErrors "auth/helpers/errors"
	utils "auth/helpers/utils"
	validate "auth/helpers/validate"
	"bytes"
	"html/template"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (conf *AuthConfig) ForgotPassword(c echo.Context) error {

	payload := Email{}
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	errs := validate.ValidateStruct(&payload)
	if errs != nil {
		return c.JSON(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	var otp_expiry int64

	conf.Db.Get(&otp_expiry, get_otp_expiry, payload.Email)

	if otp_expiry > time.Now().Unix() {
		return c.String(http.StatusBadRequest, "Please wait for the previous OTP to expire")
	}
	otp := utils.GenerateOTP()
	hashed_otp, err3 := utils.HashString(otp)
	if err3 != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}

	_, err4 := conf.Db.Exec(create_otp, payload.Email, hashed_otp)
	if err4 != nil {
		return c.String(http.StatusBadRequest, "Email not registered")
	}
	t, err := template.ParseFS(conf.Templates, "templates/forgot_password.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}
	var body bytes.Buffer
	err5 := t.Execute(&body, map[string]string{"OTP": otp})
	if err5 != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}

	subject := "Subject:Password Reset OTP\n"
	go utils.SendEmail([]string{payload.Email}, subject, body.Bytes())

	return c.String(http.StatusOK, "OTP sent to your email")
}
