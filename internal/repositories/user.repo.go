package repositories

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type UserRepository interface {
	CreateUser(username, password, role, email string) (models.User, error)
	GetUserByUsernameOrEmail(identifier string) (models.User, error)
	BlacklistToken(ctx context.Context, token string, duration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
}

type userRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewUserRepository(db *pgxpool.Pool, rdb *redis.Client) UserRepository {
	return &userRepository{db: db, rdb: rdb}
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

func (r *userRepository) BlacklistToken(ctx context.Context, token string, duration time.Duration) error {
	return r.rdb.Set(ctx, "blacklist:"+token, "true", duration).Err()
}

func (r *userRepository) IsTokenBlacklisted(ctx context.Context, token string) (bool, error) {
	val, err := r.rdb.Get(ctx, "blacklist:"+token).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return val == "true", nil
}
