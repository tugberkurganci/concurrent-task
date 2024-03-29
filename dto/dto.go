package dto

import (
	"konzek-jun/models"
)

type UserCreateRequest struct {
	Name     string `gorm:"type:varchar(100)" json:"-"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=16"`
}

type LoginRequest struct {
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type RegisterRequest struct {
	Name     string `json:"name" form:"name" validate:"required,min=1"`
	Email    string `json:"email" form:"email" validate:"required,email"`
	Password string `json:"password" form:"password" validate:"required,min=6"`
}

type UpdateUserRequest struct {
	ID    int64  `json:"id" form:"id"`
	Name  string `json:"name" form:"name" validate:"required,min=1"`
	Email string `json:"email" form:"email" validate:"required,email"`
}

type UserResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token,omitempty"`
}

func NewUserResponse(user models.User) UserResponse {
	return UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
}
