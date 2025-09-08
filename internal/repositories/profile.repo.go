package repositories

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/siddiq24/Tickitz-DB/internal/models"
)

type ProfileRepository interface {
	GetProfile(userID int) (*models.Profile, error)
	UpdateProfile(userID int, req models.UpdateProfileRequest) error
	CreateDefaultProfile(userID int) error
}

type profileRepository struct {
	db *pgxpool.Pool
}

func NewProfileRepository(db *pgxpool.Pool) ProfileRepository {
	return &profileRepository{db: db}
}

// ambil profile berdasarkan user_id
func (r *profileRepository) GetProfile(userID int) (*models.Profile, error) {
	var p models.Profile
	err := r.db.QueryRow(context.Background(),
		`SELECT user_id, COALESCE(avatar, ''), COALESCE(phone, ''), COALESCE(address, ''), COALESCE(birthdate, '1970-01-01') 
		 FROM profiles WHERE user_id=$1`, userID).
		Scan(&p.UserID, &p.Avatar, &p.Phone, &p.Address, &p.Birthdate)

	if err != nil {
		// jika belum ada profile, buat default
		_ = r.CreateDefaultProfile(userID)
		p = models.Profile{
			UserID:    userID,
			Avatar:    "",
			Phone:     "",
			Address:   "",
			Birthdate: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		return &p, nil
	}
	return &p, nil
}

// update profile
func (r *profileRepository) UpdateProfile(userID int, req models.UpdateProfileRequest) error {
	query := `
		UPDATE profiles
		SET avatar = COALESCE($1, avatar),
		    phone = COALESCE($2, phone),
		    address = COALESCE($3, address),
		    birthdate = COALESCE($4, birthdate)
		WHERE user_id = $5
	`
	_, err := r.db.Exec(context.Background(), query,
		req.Avatar, req.Phone, req.Address, req.Birthdate, userID,
	)
	return err
}

// buat default profile jika belum ada
func (r *profileRepository) CreateDefaultProfile(userID int) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO profiles (user_id, avatar, phone, address, birthdate) 
		 VALUES ($1, '', '', '', '1970-01-01') 
		 ON CONFLICT (user_id) DO NOTHING`, userID)
	return err
}
