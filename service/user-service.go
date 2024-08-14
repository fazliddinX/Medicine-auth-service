package service

import (
	pb "auth-service/generated/users"
	"auth-service/pkg/hashing"
	"auth-service/storage"
	"context"
	"log/slog"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	log *slog.Logger
	storage.UserStorage
}

func NewUserService(log *slog.Logger, storage storage.UserStorage) *UserService {
	return &UserService{log: log, UserStorage: storage}
}

func (s *UserService) GetUserProfile(ctx context.Context, in *pb.UserID) (*pb.UserResponse, error) {
	res, err := s.UserStorage.GetUserProfile(in)
	if err != nil {
		s.log.Error("Error in GetUserProfile", "error", err)
		return nil, err
	}

	return res, nil
}

func (s *UserService) UpdateUser(ctx context.Context, in *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	res, err := s.UserStorage.UpdateUser(in)
	if err != nil {
		s.log.Error("Error in UpdateUser", "error", err)
		return nil, err
	}

	return res, nil
}

func (s *UserService) DeleteUser(ctx context.Context, in *pb.UserID) (*pb.Message, error) {
	res, err := s.UserStorage.DeleteUser(in)
	if err != nil {
		s.log.Error("Error in DeleteUser", "error", err)
		return nil, err
	}

	return res, nil
}

func (s *UserService) GetAllUsers(ctx context.Context, in *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	res, err := s.UserStorage.GetAllUsers(in)
	if err != nil {
		s.log.Error("Error in GetAllUsers", "error", err)
		return nil, err
	}

	return res, nil
}

func (s *UserService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hash, err := hashing.HashPassword(in.Password)
	if err != nil {
		s.log.Error("Error in CreateUser", "error", err)
		return nil, err
	}

	in.Password = hash

	res, err := s.UserStorage.CreateUser(in)
	if err != nil {
		s.log.Error("Error in CreateUser", "error", err)
		return nil, err
	}

	return res, nil
}
