package auth

import (
	auth_handler "auth/controllers/auth_handlers"
	"auth/helpers/constants"

	"github.com/labstack/echo/v4"

	"github.com/labstack/echo/v4/middleware"
)

func Register(app *echo.Echo, config *constants.AppConfig) {
	auth := app.Group("/auth")
	auth_handler := auth_handler.AuthConfig{AppConfig: config}
	auth.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(20)))
	auth.POST("/login", auth_handler.Login)
	auth.POST("", auth_handler.SignUp)
	auth.POST("/reset-password", auth_handler.ResetPassword)
	auth.GET("/:username/:code", auth_handler.VerifyRegistration)
	auth.POST("/forgot-password", auth_handler.ForgotPassword)
	auth.POST("/forgot-password-reset", auth_handler.ForgotPasswordReset)
}
