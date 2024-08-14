package postgres

import (
	pb "auth-service/generated/users"
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strconv"
)

type UserRepo struct {
	Db *sqlx.DB
}

func (u *UserRepo) GetUserProfile(in *pb.UserID) (*pb.UserResponse, error) {
	var user pb.UserResponse

	query := `SELECT email, first_name, last_name, date_of_birth, gender FROM users WHERE id = $1 and deleted_at=0 and role != 'admin'`
	err := u.Db.QueryRow(query, in.Id).Scan(&user.Email, &user.FirstName, &user.LastName, &user.DataOfBirthday, &user.Gender)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %s not found", in.Id)
		}
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) CreateUser(in *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	res := &pb.CreateUserResponse{}

	query := `
		INSERT INTO users (email, password_hash, first_name, last_name, date_of_birth, gender, role) 
		VALUES ($1, $2, $3, $4, $5, $6, $7) returning id
	`

	err := u.Db.QueryRow(query, in.Email, in.Password, in.FirstName, in.LastName, in.DataOfBirthday, in.Gender, in.Role).
		Scan(&res.Id)
	if err != nil {
		return nil, err
	}

	res.Email = in.Email
	res.FirstName = in.FirstName
	res.LastName = in.LastName
	res.Gender = in.Gender
	res.DataOfBirthday = in.DataOfBirthday

	return res, nil
}

func (u *UserRepo) UpdateUser(in *pb.UpdateUserRequest) (*pb.UserResponse, error) {
	query := `
		UPDATE users SET first_name = $1, last_name = $2, date_of_birth = $3 
		WHERE id = $4 and deleted_at=0
		RETURNING email, first_name, last_name, date_of_birth, gender
	`

	var user pb.UserResponse
	err := u.Db.QueryRow(query, in.FirstName, in.LastName, in.DataOfBirthday, in.Id).
		Scan(&user.Email, &user.FirstName, &user.LastName, &user.DataOfBirthday, &user.Gender)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (u *UserRepo) DeleteUser(in *pb.UserID) (*pb.Message, error) {
	_, err := u.Db.Exec(`update users set
	deleted_at = date_part('epoch', current_timestamp)::INT
	where id = $1 and deleted_at = 0  and deleted_at=0`, in.Id)

	if err != nil {
		return nil, err
	}

	return &pb.Message{Message: "user deleted successfully"}, nil
}

func (u *UserRepo) GetAllUsers(in *pb.GetAllUsersRequest) (*pb.GetAllUsersResponse, error) {
	var users []*pb.UserResponse

	query := "SELECT email, first_name, last_name, date_of_birth, gender FROM users WHERE deleted_at=0 and role != 'admin'"
	args := []interface{}{}
	count := 0

	if in.Limit == 0 {
		in.Limit = 10
	}
	if in.FirstName != "" {
		count++
		a := strconv.Itoa(count)
		query += " AND first_name = $" + a
		args = append(args, in.FirstName)
	}
	if in.Gender != "" {
		count++
		a := strconv.Itoa(count)
		query += " AND gender = $" + a
		args = append(args, in.Gender)
	}
	count++
	a := strconv.Itoa(count)
	query += " LIMIT $" + a
	args = append(args, in.Limit)

	count++
	a = strconv.Itoa(count)
	query += " OFFSET $" + a
	args = append(args, in.Offset)

	rows, err := u.Db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user pb.UserResponse
		if err := rows.Scan(&user.Email, &user.FirstName, &user.LastName, &user.DataOfBirthday, &user.Gender); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &pb.GetAllUsersResponse{Users: users}, nil
}
