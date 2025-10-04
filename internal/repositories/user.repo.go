package repositories

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"github.com/siddiq24/Tickitz-DB/internal/models"
	"github.com/siddiq24/Tickitz-DB/internal/utils"
)

type UserRepository interface {
	CreateUser(username, password, role, email string) (models.Profile, error)
	GetUserByUsernameOrEmail(identifier string) (models.Profile, error)
	BlacklistToken(ctx context.Context, token string, duration time.Duration) error
	IsTokenBlacklisted(ctx context.Context, token string) (bool, error)
	UpdatePassword(userID int, oldPassword, newPassword string) error
}

type userRepository struct {
	db  *pgxpool.Pool
	rdb *redis.Client
}

func NewUserRepository(db *pgxpool.Pool, rdb *redis.Client) UserRepository {
	return &userRepository{db: db, rdb: rdb}
}

func (r *userRepository) CreateUser(username, password, role, email string) (models.Profile, error) {
	var exist bool
	sql := `SELECT EXISTS(SELECT 1 FROM users WHERE username=$1)`
	err := r.db.QueryRow(context.Background(), sql, username).Scan(&exist)
	if err != nil {
		return models.Profile{}, err
	}
	if exist {
		return models.Profile{}, errors.New("user already exist")
	}

	var user models.Profile
	err = r.db.QueryRow(context.Background(),
		`INSERT INTO users (username, password, role, email) 
		 VALUES ($1, $2, $3, $4) 
		 RETURNING id, username, role, point`,
		username, password, role, email,
	).Scan(&user.UserID, &user.Username, &user.Role, &user.Point)

	if err != nil {
		return models.Profile{}, err
	}
	return user, nil
}

func (r *userRepository) GetUserByUsernameOrEmail(identifier string) (models.Profile, error) {
	var user models.Profile
	err := r.db.QueryRow(context.Background(),
		`SELECT id, username, password, role, point FROM users WHERE username=$1 OR email=$1`,
		identifier,
	).Scan(&user.UserID, &user.Username, &user.Password, &user.Role, &user.Point)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Profile{}, errors.New("user not found")
		}
		log.Println(err)
		return models.Profile{}, fmt.Errorf("database error")
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

func (r *userRepository) UpdatePassword(userID int, oldPassword, newPassword string) error {
	var currentHashed string
	err := r.db.QueryRow(context.Background(),
		`SELECT password FROM users WHERE id=$1`, userID).Scan(&currentHashed)
	if err != nil {
		return err
	}

	// cek password lama
	if !utils.CheckPasswordHash(oldPassword, currentHashed) {
		return errors.New("old password is incorrect")
	}

	// hash password baru
	newHashed, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(context.Background(),
		`UPDATE users SET password=$1 WHERE id=$2`, newHashed, userID)
	if err != nil {
		return err
	}

	return nil
}
