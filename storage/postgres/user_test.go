package postgres

import (
	"auth-service/generated/users"
	"log"
	"testing"

	"github.com/google/uuid"
)

func TestGetUserProfile(t *testing.T) {
	db, err := Connectfortest()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userID := uuid.New().String()

	_, err = db.Exec(`INSERT INTO users(id, email, first_name, last_name, date_of_birth, gender, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, "test@example.com", "John", "Doe", "1990-01-01", "male", "patient")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := UserRepo{Db: db}

	userRes, err := userRepo.GetUserProfile(&users.UserID{Id: userID})
	if err != nil {
		log.Fatal(err)
	}

	if userRes.Email != "test@example.com" || userRes.FirstName != "John" || userRes.LastName != "Doe" {
		log.Fatal("User data does not match expected values")
	}

	_, err = db.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		log.Fatal(err)
	}
}

func TestCreateUser(t *testing.T) {
	db, err := Connectfortest()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := UserRepo{Db: db}

	createReq := &users.CreateUserRequest{
		Email:          "newuser@example.com",
		Password:       "password",
		FirstName:      "New",
		LastName:       "User",
		DataOfBirthday: "1995-01-01",
		Gender:         "female",
		Role:           "patient",
	}

	msg, err := userRepo.CreateUser(createReq)
	if err != nil {
		log.Fatal(err)
	}

	if msg.Id == "" {
		log.Fatalf("id not found")
	}

	_, err = db.Exec(`DELETE FROM users WHERE id = $1`, msg.Id)
	if err != nil {
		log.Fatal(err)
	}
}

func TestUpdateUser(t *testing.T) {
	db, err := Connectfortest()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	a := uuid.New().String()

	_, err = db.Exec(`INSERT INTO users(id, email, first_name, last_name, date_of_birth, gender, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		a, "updateuser@example.com", "OldFirstName", "OldLastName", "1980-01-01", "male", "patient")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := UserRepo{Db: db}

	updateReq := &users.UpdateUserRequest{
		Id:             a,
		FirstName:      "NewFirstName",
		LastName:       "NewLastName",
		DataOfBirthday: "1985-05-05",
	}
	_, err = userRepo.UpdateUser(updateReq)
	if err != nil {
		log.Fatal(err)
	}

	userRes, err := userRepo.GetUserProfile(&users.UserID{Id: a})
	if err != nil {
		log.Fatal(err, "Hello World")
	}
	if userRes.FirstName != "NewFirstName" || userRes.LastName != "NewLastName" {
		log.Fatal("User update failed, data does not match expected values")
	}

	_, err = db.Exec(`DELETE FROM users WHERE id = $1`, a)
	if err != nil {
		log.Fatal(err, "Hello World")
	}
}

func TestDeleteUser(t *testing.T) {
	db, err := Connectfortest()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userID := uuid.New().String()

	_, err = db.Exec(`INSERT INTO users(id, email, first_name, last_name, date_of_birth, gender, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		userID, "deleteuser@example.com", "Delete", "Me", "2000-01-01", "female", "patient")
	if err != nil {
		log.Fatal(err)
	}

	userRepo := UserRepo{Db: db}

	deleteReq := &users.UserID{Id: userID}
	_, err = userRepo.DeleteUser(deleteReq)
	if err != nil {
		log.Fatal(err)
	}

	userRes, err := userRepo.GetUserProfile(deleteReq)
	if err == nil || userRes != nil {
		log.Fatal("User was not deleted")
	}

	_, err = db.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		log.Fatal(err)
	}
}
