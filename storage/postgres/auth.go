package postgres

import (
	"auth-service/pkg/hashing"
	"auth-service/pkg/models"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
)

type AuthRepo struct {
	Db *sqlx.DB
}

func (a *AuthRepo) Register(user models.User) (string, error) {
	var id string

	err := a.Db.QueryRow(`insert into users
    (email, password_hash, first_name, last_name, date_of_birth, gender) 
	values ($1, $2, $3, $4, $5, $6) returning id`,
		user.Email, user.Password, user.FirstName, user.LastName, user.DateOfBirth, user.Gender).
		Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (a *AuthRepo) AddAdmin(admin models.AddingAdmin) error {
	var role string
	err := a.Db.QueryRow("select role from users where password_hash = $1", admin.Password).Scan(&role)
	if err != nil {
		return err
	}
	if role != "admin" {
		return errors.New("User already exists")
	}

	hash, err := hashing.HashPassword(admin.Password)
	if err != nil {
		return err
	}

	_, err = a.Db.Exec(`insert into users
    (role, password_hash, email) 
	values ($1, $2, $3)`,
		"admin", hash, admin.Email)
	return err
}

func (a *AuthRepo) Login(user models.LoginRequest) (models.StorageLogin, error) {
	var res models.StorageLogin

	err := a.Db.Get(&res, `select id, role, password_hash from users where email = $1 and deleted_at=0`, user.Email)
	if errors.Is(err, sql.ErrNoRows) {
		return models.StorageLogin{}, errors.New("user not found")
	}
	if err != nil {
		return models.StorageLogin{}, err
	}

	return res, nil
}
