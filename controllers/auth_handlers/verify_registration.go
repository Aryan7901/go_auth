package auth_handler

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (conf *AuthConfig) VerifyRegistration(c echo.Context) error {
	username := c.Param("username")
	activation_code := c.Param("code")
	verification_status := VerifyRegistration{}
	err := conf.Db.Get(&verification_status, get_user, username)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid Username")
	}
	if verification_status.Expiry < time.Now().Unix() {
		return c.String(http.StatusBadRequest, "Activation Code Expired")
	}
	if verification_status.IsVerified {
		return c.String(http.StatusBadRequest, "Already Verified")
	}
	err2 := bcrypt.CompareHashAndPassword([]byte(verification_status.ActivationCode), []byte(activation_code))
	if err2 != nil {
		return c.String(http.StatusBadRequest, "Invalid Activation Code")
	}

	_, err3 := conf.Db.Exec(update_user, username)
	if err3 != nil {
		return c.String(http.StatusBadRequest, err2.Error())
	}
	return c.String(http.StatusOK, "Verification Successful!")
}
