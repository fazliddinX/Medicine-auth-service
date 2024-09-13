package postgres

import (
	"auth-service/pkg/models"
	"github.com/google/uuid"
	"log"
	"testing"
)

func TestRegister(t *testing.T) {
	db, err := Connectfortest()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Create a repository instance
	authRepo := AuthRepo{Db: db}

	// Define a test user
	testUser := models.User{
		Email:       "testregister@example.com",
		Password:    "hashedpassword",
		FirstName:   "John",
		LastName:    "Doe",
		DateOfBirth: "1990-01-01",
		Gender:      "male",
	}

	// Call the Register method
	userID, err := authRepo.Register(testUser)
	if err != nil {
		t.Fatalf("unexpected error during registration: %v", err)
	}

	if userID == "" {
		t.Fatalf("expected a user ID, got empty string")
	}

	// Clean up the test data
	_, err = db.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		log.Fatal(err)
	}
}

func TestLogin(t *testing.T) {
	db, err := Connectfortest()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Insert a test user into the database
	userID := uuid.New().String()
	_, err = db.Exec(`INSERT INTO users(id, email, password_hash, first_name, last_name, date_of_birth, gender, role)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		userID, "testlogin@example.com", "hashedpassword", "John", "Doe", "1990-01-01", "male", "patient")
	if err != nil {
		log.Fatal(err)
	}

	// Create a repository instance
	authRepo := AuthRepo{Db: db}

	// Define a login request
	loginReq := models.LoginRequest{
		Email:    "testlogin@example.com",
		Password: "hashedpassword",
	}

	// Call the Login method
	storageLogin, err := authRepo.Login(loginReq)
	if err != nil {
		t.Fatalf("unexpected error during login: %v", err)
	}

	if storageLogin.Role != "patient" || storageLogin.Password != "hashedpassword" {
		t.Fatalf("login data does not match expected values")
	}

	// Clean up the test data
	_, err = db.Exec(`DELETE FROM users WHERE id = $1`, userID)
	if err != nil {
		log.Fatal(err)
	}
}
