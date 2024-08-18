package handlers

import (
	"net/http"
	"regexp"

	"boilerplate/internal/models"
	"boilerplate/internal/services"
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
	// Validate email
	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	if !emailRegex.MatchString(email) {
		return c.String(http.StatusBadRequest, "Invalid email format")
	}

	// Generate and send OTP
	otp := services.GenerateOTP()
	err := services.SendOTP(c.Request().Context(), email, otp)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to send OTP")
	}

	// Store OTP in database
	user, err := models.GetOrCreateUser(c.Request().Context(), email)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create user")
	}
	err = models.CreateOTP(c.Request().Context(), user.ID, otp)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to create OTP")
	}
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
