package user

import (
	"context"
	"time"
)

type Service interface {
	// CreateUser create new user as given.
	//
	// CreateUser expects the given  user ID in the given
	// user already assigned.
	CreateUser(ctx context.Context, users User) (int64, error)

	// UpdatePassword, for the given user ID, updates user's
	// password with the new password. Before updating, it
	// checks whether the current password are correct.
	UpdatePassword(ctx context.Context, userID int64, newPassword string, currentPassword string) error

	// ResetPassword, for the given user ID, updates user's
	// password with the new password without checking the
	// current password.
	ResetPassword(ctx context.Context, userID int64, newPassword string) error

	// LoginBasic checks the given email and password with
	// the actual data. It returns a token and the encapsulated
	// data if the login process is success.
	LoginBasic(ctx context.Context, email string, password string) (string, TokenData, error)

	// ValidateToken validates the given token and returns the
	// data encapsulated in the token if the given token is
	// valid.
	ValidateToken(ctx context.Context, token string) (TokenData, error)

	// RefreshToken validates the given token and returns a
	// new token with the same encapsulated data but refreshed.
	//
	// RefreshToken is used to avoid expired token.
	RefreshToken(ctx context.Context, token string) (string, error)
}

type User struct {
	ID          int64
	Email       string
	Name        string
	Password    string
	PhoneNumber string
	CreateTime  time.Time
	UpdateTime  time.Time
}

// TokenData is the data that are encapsulated in a token.
type TokenData struct {
	UserID int64
	Email  string
}
