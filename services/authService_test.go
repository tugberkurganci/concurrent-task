package services

import (
	"errors"
	"konzek-jun/mocks/repository"
	"konzek-jun/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var mockAuthRepo *repository.MockUserRepository
var mockAuthService AuthService

func setupAuth(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockAuthRepo = repository.NewMockUserRepository(ctrl)
	mockAuthService = NewAuthService(mockAuthRepo)

	return func() {
		service = nil
		ctrl.Finish()
	}
}

func TestAuthService_VerifyCredential_Success(t *testing.T) {
	// Test için hazırlıkları yap
	td := setupAuth(t)
	defer td()

	email := "test@example.com"
	password := "password"
	hashedPassword := "$2a$12$3AX3dyNLdk3D8EQri2w2f.mgU8pWDDn2Slehr7c1dUB1DP4WxH3L6"

	// Mock repository'den beklenen değerlerin ayarlanması
	mockAuthRepo.EXPECT().FindByEmail(email).Return(models.User{Email: email, Password: hashedPassword}, nil)

	// Servis fonksiyonunun çağrılması
	err := mockAuthService.VerifyCredential(email, password)

	// Hata kontrolü
	assert.NoError(t, err)
}

func TestAuthService_VerifyCredential_UserNotFound(t *testing.T) {
	// Test için hazırlıkları yap
	td := setupAuth(t)
	defer td()

	email := "test@example.com"
	password := "password"

	// Mock repository'den beklenen değerlerin ayarlanması
	mockAuthRepo.EXPECT().FindByEmail(email).Return(models.User{}, errors.New("user not found"))

	// Servis fonksiyonunun çağrılması
	err := mockAuthService.VerifyCredential(email, password)

	// Hata kontrolü
	assert.Error(t, err)
}

func TestAuthService_VerifyCredential_WrongPassword(t *testing.T) {
	// Test için hazırlıkları yap
	td := setupAuth(t)
	defer td()

	email := "test@example.com"
	password := "wrong_password"
	hashedPassword := "$2a$10$XkO/7pHBkHZvqK0b54R0YOMNc6q5aP/V0TbS3VIsffzY9j28W2PK6"

	// Mock repository'den beklenen değerlerin ayarlanması
	mockAuthRepo.EXPECT().FindByEmail(email).Return(models.User{Email: email, Password: hashedPassword}, nil)

	// Servis fonksiyonunun çağrılması
	err := mockAuthService.VerifyCredential(email, password)

	// Hata kontrolü
	assert.Error(t, err)
}
