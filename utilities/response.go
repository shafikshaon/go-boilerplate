package utilities

import (
	"go-boilerplate/models"
	"go-boilerplate/models/response"
)

func ToUserResponse(user *models.User) *response.UserResponse {
	return &response.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
