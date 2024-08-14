package models

type User struct {
	Email       string `json:"email" db:"email"`
	Password    string `json:"password" db:"password_hash"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"`
	Gender      string `json:"gender" db:"gender"`
}

type UserResponse struct {
	Id          string `json:"id" db:"id"`
	Email       string `json:"email" db:"email"`
	FirstName   string `json:"first_name" db:"first_name"`
	LastName    string `json:"last_name" db:"last_name"`
	DateOfBirth string `json:"date_of_birth" db:"date_of_birth"`
}

type LoginRequest struct {
	Email    string `json:"email" db:"email"`
	Password string `json:"password" db:"password"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    string `json:"expires_in"`
}

type StorageLogin struct {
	Id       string `json:"id" db:"id"`
	Password string `json:"password" db:"password_hash"`
	Role     string `json:"role" db:"role"`
}

type Success struct {
	Message string `json:"message"`
}

type Error struct {
	Error string `json:"error"`
}

type AddingAdmin struct {
	Password string `json:"password" db:"password_hash"`
	Email    string `json:"email" db:"email"`
}
