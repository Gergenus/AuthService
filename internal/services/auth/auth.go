package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Gergenus/AuthService/internal/domain/models"
	jwtpkg "github.com/Gergenus/AuthService/internal/pkg/jwt"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCreadentials = errors.New("invalid credetials")
)

type Auth struct {
	log          *slog.Logger
	userSaver    UserSaver
	userProvider UserProvider
}

type UserSaver interface {
	SaveUser(ctx context.Context, username string, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	UserID(ctx context.Context, username string) (int64, error)
}

func NewAuth(log *slog.Logger, userSaver UserSaver, userProvider UserProvider) *Auth {
	return &Auth{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
	}
}

func (a *Auth) SignIn(ctx context.Context, email string, password string) (token string, err error) {
	const op = "auth.SignIn"

	log := a.log.With(slog.String("op", op))
	log.Info("SigningIn user")

	user, err := a.userProvider.User(ctx, email)
	if err != nil {
		a.log.Error("internal error", "error", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCreadentials)
	}
	if err := bcrypt.CompareHashAndPassword(user.HashPassword, []byte(password)); err != nil {
		a.log.Info("invalid credentials", "error", err)
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCreadentials)
	}

	log.Info("user logged in succesfully")
	tkn, err := jwtpkg.NewToken(user)
	if err != nil {
		log.Error("faild create jwt token", "error", err)
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return tkn, nil
}

func (a *Auth) GetUser(ctx context.Context, username string) (userID int64, err error) {
	const op = "auth.GetUser"

	log := a.log.With(slog.String("op", op))
	log.Info("getting user")

	id, err := a.userProvider.UserID(ctx, username)
	if err != nil {
		a.log.Error("internal error", "error", err)
		return 0, fmt.Errorf("%s: %w", op, ErrInvalidCreadentials)
	}

	return id, nil
}

func (a *Auth) RegisterNewUser(ctx context.Context, username string, email string, password string) (userID int64, err error) {
	const op = "auth.RegisterNewUser"

	log := a.log.With(slog.String("op", op))
	log.Info("registering user")

	HashPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", "error", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.userSaver.SaveUser(ctx, username, email, HashPassword)
	if err != nil {
		log.Error("failed to save user", "error", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("User has been created")
	return id, nil
}
