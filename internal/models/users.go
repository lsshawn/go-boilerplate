package models

import (
	"context"
	"time"

	"boilerplate/internal/database"
)

type User struct {
	ID        int64
	Email     string
	CreatedAt time.Time
}

type OTP struct {
	ID        int64
	UserID    int64
	OTP       string
	CreatedAt time.Time
	ExpiresAt time.Time
}

func GetOrCreateUser(ctx context.Context, email string) (*User, error) {
	var user User
	err := database.DB.QueryRowContext(ctx, "SELECT id, email, created_at FROM users WHERE email = ?", email).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err == nil {
		return &user, nil
	}

	result, err := database.DB.ExecContext(ctx, "INSERT INTO users (email) VALUES (?)", email)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &User{ID: id, Email: email, CreatedAt: time.Now()}, nil
}

func CreateOTP(ctx context.Context, userID int64, otp string) error {
	_, err := database.DB.ExecContext(ctx, "INSERT INTO otps (user_id, otp, expires_at) VALUES (?, ?, ?)",
		userID, otp, time.Now().Add(5*time.Minute))
	return err
}

func ValidateOTP(ctx context.Context, userID int64, otp string) bool {
	var count int
	err := database.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM otps WHERE user_id = ? AND otp = ? AND expires_at > ?",
		userID, otp, time.Now()).Scan(&count)
	if err != nil {
		return false
	}
	return count > 0
}
