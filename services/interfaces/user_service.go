package interfaces

import (
	"go-boilerplate/models/request"
	"go-boilerplate/models/response"
)

type UserService interface {
	CreateUser(req *request.CreateUserRequest) (*response.UserResponse, error)
	GetUserByID(id uint) (*response.UserResponse, error)
	GetUsers(page, perPage int) (*response.PaginationResponse, error)
	UpdateUser(id uint, req *request.UpdateUserRequest) (*response.UserResponse, error)
	DeleteUser(id uint) error
	Login(req *request.LoginRequest) (*response.LoginResponse, error)
	Logout(userID uint) error
}
