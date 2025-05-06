package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/Gergenus/AuthService/internal/domain/models"
	"github.com/Gergenus/AuthService/internal/pkg/database"
	"github.com/lib/pq"
)

var (
	ErrUserEXists   = errors.New("user already exists")
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email exists")
)

type PostgresRepository struct {
	log *slog.Logger
	DB  database.PostgresDB
}

func NewPostgresRepository(log *slog.Logger, DB database.PostgresDB) *PostgresRepository {
	return &PostgresRepository{
		log: log,
		DB:  DB,
	}
}

func (p *PostgresRepository) SaveUser(ctx context.Context, username string, email string, passHash []byte) (uid int64, err error) {
	const op = "repository.SaveUser"
	var id int
	err = p.DB.DB.QueryRow("INSERT INTO users (username, email, hash_password) VALUES($1, $2, $3) RETURNING id", username, email, passHash).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Constraint, "username") {
					p.log.Error("username already exists", "error", err)
					return 0, fmt.Errorf("%s: %w", op, ErrUserEXists)
				}
				if strings.Contains(pgErr.Constraint, "email") {
					p.log.Error("email already exists", "error", err)
					return 0, fmt.Errorf("%s: %w", op, ErrEmailExists)
				}
			}
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return int64(id), nil
}

func (p *PostgresRepository) User(ctx context.Context, email string) (models.User, error) {
	const op = "repository.User"
	var user models.User
	err := p.DB.DB.QueryRow("SELECT id, username, email, hash_password FROM users WHERE email=$1", email).Scan(&user.ID, &user.Username, &user.Email, &user.HashPassword)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.log.Error("user not found", "error", ErrUserNotFound)
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		p.log.Error("internal error", "error", err)
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return user, nil
}

func (p *PostgresRepository) UserID(ctx context.Context, username string) (int64, error) {
	const op = "repository.UserID"
	var id int64
	err := p.DB.DB.QueryRow("SELECT id FROM users WHERE username = $1", username).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p.log.Error("user not found", "error", ErrUserNotFound)
			return 0, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}
		p.log.Error("internal error", "error", err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}
