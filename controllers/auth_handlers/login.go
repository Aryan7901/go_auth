package auth_handler

import (
	"auth/helpers/constants"
	AppErrors "auth/helpers/errors"
	validate "auth/helpers/validate"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

func (conf *AuthConfig) Login(c echo.Context) error {
	payload := Login{}
	if err := c.Bind(&payload); err != nil {
		return c.String(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	errs := validate.ValidateStruct(&payload)
	if errs != nil {
		return c.JSON(http.StatusBadRequest, AppErrors.InvalidPayload)
	}
	user_details := UserCredentials{}
	err2 := conf.Db.Get(&user_details, get_user_details, payload.Email)

	if err2 != nil {
		return c.String(http.StatusBadRequest, "Invalid Email/Password")
	} else if !user_details.IsVerified {
		return c.String(http.StatusBadRequest, "User is Unverified")
	} else if err3 := bcrypt.CompareHashAndPassword([]byte(user_details.Password), []byte(payload.Password)); err3 != nil {
		return c.String(http.StatusBadRequest, "Invalid Email/Password")
	} else {
		token := jwt.New(jwt.SigningMethodHS256)

		// Set claims
		claims := token.Claims.(jwt.MapClaims)
		claims["username"] = user_details.UserName
		claims["email"] = user_details.Email
		claims["exp"] = time.Now().Add(time.Hour * 2).Unix()

		// Generate encoded token and send it as response.
		t, err4 := token.SignedString([]byte(constants.JwtSecret))
		if err4 != nil {
			return c.String(http.StatusInternalServerError, err4.Error())
		}
		return c.JSON(http.StatusOK, map[string]string{
			"token": t,
		})
	}
}
