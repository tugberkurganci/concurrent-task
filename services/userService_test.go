package services

import (
	"errors"
	"konzek-jun/dto"
	"konzek-jun/mocks/repository"
	"konzek-jun/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var mockRepository *repository.MockUserRepository
var mockService UserService

var FakeUser = dto.RegisterRequest{

	Name:  "John Doe",
	Email: "john@example.com",
}

func setupUser(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockRepository = repository.NewMockUserRepository(ctrl)
	mockService = NewUserService(mockRepository)

	return func() {
		service = nil
		ctrl.Finish()
	}
}

func TestUserService_CreateUser_Success(t *testing.T) {
	// Test için hazırlıkları yap
	td := setupUser(t)
	defer td()

	// Mock repository'den beklenen değerlerin ayarlanması
	mockRepository.EXPECT().FindByEmail(FakeUser.Email).Return(models.User{}, errors.New("some error"))
	mockRepository.EXPECT().InsertUser(gomock.Any()).Return(models.User{
		Name:  "John Doe",
		Email: "john@example.com"}, nil)

	// Servis fonksiyonunun çağrılması
	result, err := mockService.CreateUser(FakeUser)

	// Hata kontrolü
	assert.NoError(t, err)
	assert.Equal(t, result.Email, FakeUser.Email)
}

func TestUserService_FindUserByEmail_Success(t *testing.T) {
	// Test için hazırlıkları yap
	td := setupUser(t)
	defer td()

	// Mock repository'den beklenen değerlerin ayarlanması
	mockRepository.EXPECT().FindByEmail(gomock.Any()).Return(models.User{ID: 1, Email: "x@x.com"}, nil)

	// Servis fonksiyonunun çağrılması
	result, err := mockService.FindUserByEmail("x@x.com")

	// Hata kontrolü
	assert.NoError(t, err)
	assert.Equal(t, result.Email, "x@x.com")
}

func TestUserService_FindUserById_Success(t *testing.T) {
	// Test için hazırlıkları yap
	td := setupUser(t)
	defer td()

	// Mock repository'den beklenen değerlerin ayarlanması
	mockRepository.EXPECT().FindByUserID(gomock.Any()).Return(models.User{ID: 1, Email: "x@x.com"}, nil)

	// Servis fonksiyonunun çağrılması
	result, err := mockService.FindUserByID("1")

	// Hata kontrolü
	assert.NoError(t, err)
	assert.Equal(t, result.Email, "x@x.com")
}

func TestUserService_UserUpdate_Success(t *testing.T) {
	// Test için hazırlıkları yap
	td := setupUser(t)
	defer td()

	// Mock repository'den beklenen değerlerin ayarlanması
	mockRepository.EXPECT().UpdateUser(gomock.Any()).Return(models.User{ID: 1, Email: "x@x.com"}, nil)

	// Servis fonksiyonunun çağrılması
	result, err := mockService.UpdateUser(dto.UpdateUserRequest{ID: 1, Email: "x@x.com"})

	// Hata kontrolü
	assert.NoError(t, err)
	assert.Equal(t, result.Email, "x@x.com")
}
