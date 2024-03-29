package services

import (
	"errors"
	"fmt"
	"konzek-jun/loggerx"
	"konzek-jun/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	VerifyCredential(email string, password string) error
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{
		userRepo: userRepo,
	}
}

func (c *authService) VerifyCredential(email string, password string) error {
	loggerx.Info("Verifying user credential")

	user, err := c.userRepo.FindByEmail(email)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while finding user by email: %s", err))
		return err
	}

	isValidPassword := comparePassword(user.Password, []byte(password))
	if !isValidPassword {
		return errors.New("failed to login. check your credential")
	}

	loggerx.Info("User credential verified successfully")
	return nil
}

func comparePassword(hashedPwd string, plainPassword []byte) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPassword)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while comparing passwords:%s", err))
		return false
	}
	return true
}
