package auth_handler

import (
	"auth/helpers/constants"
	AppErrors "auth/helpers/errors"
	utils "auth/helpers/utils"
	validate "auth/helpers/validate"
	"bytes"
	"html/template"
	"net/http"
	"strings"
	"time"

	emailverifier "github.com/AfterShip/email-verifier"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var (
	verifier = emailverifier.NewVerifier().EnableAutoUpdateDisposable()
)

func (conf *AuthConfig) SignUp(c echo.Context) error {
	payload := SignUp{}
	err := c.Bind(&payload)
	if err != nil {
		return c.String(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	errs := validate.ValidateStruct(&payload)
	if errs != nil {
		return c.JSON(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	email := strings.TrimSpace(payload.Email)
	ret, err := verifier.Verify(email)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Email Verification failed")
	} else if !ret.Syntax.Valid {
		return c.String(http.StatusBadRequest, "Invalid Email")
	} else if ret.Disposable {
		return c.String(http.StatusBadRequest, "Provided Email is disposable")
	}
	result_channel := utils.HashStringAsync(payload.Password)
	verification_status := VerificationStatus{}
	error := conf.Db.Get(&verification_status, check_verification, email)
	if error == nil && verification_status.IsVerified {
		return c.String(http.StatusBadRequest, "Already Registered")
	}
	if error == nil && verification_status.Expiry > time.Now().Unix() {
		return c.String(http.StatusBadRequest, "Please wait for the previous activation code to expire")
	}
	result := <-result_channel
	if result.Error != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}
	password_hash := result.Hash
	uid := uuid.New().String()
	hashed_uid, err := utils.HashString(uid)
	if err != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}

	_, err2 := conf.Db.Exec(user_create, payload.UserName, password_hash, payload.Email, hashed_uid)
	if err2 != nil {
		return c.String(http.StatusBadRequest, "Username already taken /Email recently registered")
	}
	subject := "Subject:Confirm your Queerky registration\n"
	t, err := template.ParseFS(conf.Templates, "templates/verification.html")
	if err != nil {
		return c.String(http.StatusInternalServerError, AppErrors.InternalServerError)
	}
	var body bytes.Buffer
	t.Execute(&body, struct {
		Host     string
		Username string
		Uid      string
	}{
		Host:     constants.Server,
		Username: payload.UserName,
		Uid:      uid,
	})
	go utils.SendEmail([]string{email}, subject, body.Bytes())
	// err3 := utils.SendEmail([]string{email}, subject, body.Bytes())
	// if err3 != nil {
	// 	return c.String(http.StatusInternalServerError, err3.Error())
	// }
	return c.String(http.StatusOK, "Registration Successful!,Please Verify Your Account")
}
