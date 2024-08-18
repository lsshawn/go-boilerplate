package handlers

import (
	"net/http"

	"boilerplate/views/account"

	"github.com/labstack/echo/v4"
)

func AccountPage(c echo.Context) error {
	return render(c, account.Index())
}

func RequestOTP(c echo.Context) error {
	email := c.FormValue("email")
	if email == "" {
		return c.String(http.StatusBadRequest, "Email is required")
	}
	// TODO: Implement OTP generation and sending
	return render(c, account.OTPForm(email))
}

func ValidateOTP(c echo.Context) error {
	email := c.FormValue("email")
	otp := c.FormValue("otp")
	if email == "" || otp == "" {
		return c.String(http.StatusBadRequest, "Email and OTP are required")
	}
	// TODO: Implement OTP validation
	return render(c, account.LoggedIn(email))
}

func Logout(c echo.Context) error {
	// TODO: Implement logout logic
	return render(c, account.LoginForm())
}
