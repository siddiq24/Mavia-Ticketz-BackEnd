package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type UserRepository interface {
	CreateUser(username, password, role, email string) (models.User, error)
	GetUserByUsernameOrEmail(identifier string) (models.User, error)
}

type userRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(username, password, role, email string) (models.User, error) {
	var exist bool
	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`
	err := r.db.QueryRow(context.Background(), sql, username).Scan(&exist)
	if err != nil {
		return models.User{}, err
	}
	if exist {
		return models.User{}, errors.New("user already exist")
	}

	var user models.User
	err = r.db.QueryRow(context.Background(),
		`INSERT INTO users (username, password, role, email) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, username, role, email`,
		username, password, role, email,
	).Scan(&user.ID, &user.Username, &user.Role, &user.Email)

	if err != nil {
		return models.User{}, err
	}
	return user, nil

}

func (r *userRepository) GetUserByUsernameOrEmail(identifier string) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(context.Background(),
		`SELECT id, username, password, role, email FROM users WHERE username=$1 OR email=$1`,
		identifier,
	).Scan(&user.ID, &user.Username, &user.Password, &user.Role, &user.Email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, errors.New("user not found")
		}
		return models.User{}, fmt.Errorf("database error: %w", err)
	}

	return user, nil
}
