package auth

import (
	"context"
	"errors"
	"os"

	"github.com/Khaym03/REG/internal/browser"
	"github.com/Khaym03/REG/internal/session"
)

type AuthService interface {
	Login(context.Context, session.Session, User) error
	Logout(context.Context, session.Session) error
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u User) Validate() error {
	if u.Username == "" {
		return errors.New("username is required")
	}

	if u.Password == "" {
		return errors.New("password is required")
	}

	return nil
}

type (
	PageFunc = browser.PageFunc
)

func LoadCredential() User {
	return User{
		Username: os.Getenv("REG_TEST_USERNAME"),
		Password: os.Getenv("REG_TEST_PASSWORD"),
	}
}
