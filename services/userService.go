package services

import (
	"errors"
	"fmt"
	"konzek-jun/dto"
	"konzek-jun/loggerx"
	"konzek-jun/models"
	"konzek-jun/repository"

	"github.com/mashingan/smapping"
)

type UserService interface {
	CreateUser(registerRequest dto.RegisterRequest) (*dto.UserResponse, error)
	UpdateUser(updateUserRequest dto.UpdateUserRequest) (*dto.UserResponse, error)
	FindUserByEmail(email string) (*dto.UserResponse, error)
	FindUserByID(userID string) (*dto.UserResponse, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (c *userService) UpdateUser(updateUserRequest dto.UpdateUserRequest) (*dto.UserResponse, error) {
	loggerx.Info("UpdateUser function called")

	user := models.User{}
	err := smapping.FillStruct(&user, smapping.MapFields(&updateUserRequest))
	if err != nil {
		loggerx.Error(fmt.Sprintf("Failed to map user: %s", err))
		return nil, err
	}

	user, err = c.userRepo.UpdateUser(user)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while updating user: %s", err))
		return nil, err
	}

	res := dto.NewUserResponse(user)
	loggerx.Info("User updated successfully")
	return &res, nil
}

func (c *userService) CreateUser(registerRequest dto.RegisterRequest) (*dto.UserResponse, error) {
	loggerx.Info("CreateUser function called")

	user, err := c.userRepo.FindByEmail(registerRequest.Email)
	if err == nil {
		loggerx.Error("User already exists")
		return nil, errors.New("user already exists")
	}

	err = smapping.FillStruct(&user, smapping.MapFields(&registerRequest))
	if err != nil {
		loggerx.Error(fmt.Sprintf("Failed to map user: %s", err))
		return nil, err
	}

	user, _ = c.userRepo.InsertUser(user)
	res := dto.NewUserResponse(user)
	loggerx.Info("User created successfully")
	return &res, nil
}

func (c *userService) FindUserByEmail(email string) (*dto.UserResponse, error) {
	loggerx.Info("FindUserByEmail function called")

	user, err := c.userRepo.FindByEmail(email)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while finding user by email: %s", err))
		return nil, err
	}

	userResponse := dto.NewUserResponse(user)
	loggerx.Info("User found by email successfully")
	return &userResponse, nil
}

func (c *userService) FindUserByID(userID string) (*dto.UserResponse, error) {
	loggerx.Info("FindUserByID function called")

	user, err := c.userRepo.FindByUserID(userID)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while finding user by ID: %s", err))
		return nil, err
	}

	userResponse := dto.UserResponse{}
	err = smapping.FillStruct(&userResponse, smapping.MapFields(&user))
	if err != nil {
		loggerx.Error(fmt.Sprintf("Failed to map user response: %s", err))
		return nil, err
	}

	loggerx.Info("User found by ID successfully")
	return &userResponse, nil
}
