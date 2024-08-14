package storage

import (
	pb "auth-service/generated/users"
	"auth-service/pkg/models"
	"auth-service/storage/postgres"
	"github.com/jmoiron/sqlx"
)

type AuthStorage interface {
	AddAdmin(admin models.AddingAdmin) error
	Register(user models.User) (string, error)
	Login(user models.LoginRequest) (models.StorageLogin, error)
}

type UserStorage interface {
	GetUserProfile(in *pb.UserID) (*pb.UserResponse, error)
	UpdateUser(in *pb.UpdateUserRequest) (*pb.UserResponse, error)
	DeleteUser(in *pb.UserID) (*pb.Message, error)
	GetAllUsers(in *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error)
	CreateUser(in *pb.CreateUserRequest) (*pb.CreateUserResponse, error)
}

func NewAuthRepo(db *sqlx.DB) AuthStorage {
	return &postgres.AuthRepo{db}
}

func NewUserRepo(db *sqlx.DB) UserStorage {
	return &postgres.UserRepo{Db: db}
}
