package service

import (
	"auth-service/pkg/hashing"
	"auth-service/pkg/models"
	"auth-service/pkg/token"
	"auth-service/storage"
	"log/slog"
	"time"
)

type AuthService interface {
	AddAdmin(admin models.AddingAdmin) error
	Register(user models.User) (models.UserResponse, error)
	Login(user models.LoginRequest) (models.LoginResponse, error)
}

func NewAuthService(log *slog.Logger, st storage.AuthStorage) AuthService {
	return &authService{log, st}
}

type authService struct {
	logger *slog.Logger
	storage.AuthStorage
}

func (s *authService) AddAdmin(admin models.AddingAdmin) error {
	err := s.AuthStorage.AddAdmin(admin)
	return err
}

func (a *authService) Register(user models.User) (models.UserResponse, error) {

	hash, err := hashing.HashPassword(user.Password)
	if err != nil {
		a.logger.Error("Failed to hash password", "error", err)
		return models.UserResponse{}, err
	}

	user.Password = hash

	id, err := a.AuthStorage.Register(user)
	if err != nil {
		a.logger.Error("Error registering user", "error", err)
		return models.UserResponse{}, err
	}

	req := models.UserResponse{
		Id:          id,
		Email:       user.Email,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		DateOfBirth: user.DateOfBirth,
	}

	return req, nil
}

func (s *authService) Login(user models.LoginRequest) (models.LoginResponse, error) {
	res, err := s.AuthStorage.Login(user)
	if err != nil {
		s.logger.Error("Error login user", "error", err)
		return models.LoginResponse{}, err
	}

	check := hashing.CheckPasswordHash(res.Password, user.Password)
	if !check {
		s.logger.Warn("invalid password")
		return models.LoginResponse{}, err
	}

	accessToken, err := token.GenerateAccessToken(res.Id, res.Role)
	if err != nil {
		s.logger.Error("Error generating access token", "error", err)
		return models.LoginResponse{}, err
	}

	refreshToken, err := token.GenerateRefreshToken(res.Id, res.Role)
	if err != nil {
		s.logger.Error("Error generating refresh token", "error", err)
		return models.LoginResponse{}, err
	}

	loggiResponse := models.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    time.Now().Add(time.Minute * 10).String(),
	}

	return loggiResponse, nil
}
