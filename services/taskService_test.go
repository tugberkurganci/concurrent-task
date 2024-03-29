package services

import (
	"konzek-jun/mocks/repository"
	"konzek-jun/models"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var mockRepo *repository.MockTaskRepository
var service TaskService

var FakeData = []models.Task{
	{Id: 1, Title: "Task 1", Content: "Description 1", Status: true},
	{Id: 2, Title: "Task 2", Content: "Description 2", Status: false},
	{Id: 3, Title: "Task 3", Content: "Description 3", Status: true},
}

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockRepo = repository.NewMockTaskRepository(ctrl)
	service = NewTaskService(mockRepo)

	return func() {
		service = nil
		ctrl.Finish()
	}
}

func TestDefaultTaskService_TaskGetAll_Success(t *testing.T) {
	// Test için hazırlıkları yap
	td := setup(t)
	defer td()

	// Mock repository'den beklenen değerlerin ayarlanması
	mockRepo.EXPECT().GetAll().Return(FakeData, nil)

	// Servis fonksiyonunun çağrılması
	result, err := service.TaskGetAll()

	// Hata kontrolü
	if err != nil {
		t.Error(err)
	}
	assert.NotEmpty(t, result)
}

func TestDefaultTaskService_TaskInsert_Success(t *testing.T) {
	// Test için hazırlıkları yap
	defer setup(t)()

	// Mock repository'den beklenen değerlerin ayarlanması
	task := models.Task{Id: 1, Title: "Test Task", Content: "Test Description"}
	mockRepo.EXPECT().Insert(task).Return(int64(1), nil)

	// Servis fonksiyonunun çağrılması
	err := service.TaskInsert(task)

	// Hata kontrolü
	assert.NoError(t, err)
}

func TestDefaultTaskService_TaskDelete_Success(t *testing.T) {
	// Test için hazırlıkları yap
	defer setup(t)()

	// Mock repository'den beklenen değerlerin ayarlanması
	taskID := 1
	mockRepo.EXPECT().Delete(taskID).Return(nil)

	// Servis fonksiyonunun çağrılması
	err := service.TaskDelete(taskID)

	// Hata kontrolü
	assert.NoError(t, err)
}

func TestDefaultTaskService_TaskUpdate_Success(t *testing.T) {
	// Test için hazırlıkları yap
	defer setup(t)()

	// Mock repository'den beklenen değerlerin ayarlanması
	task := models.Task{Id: 1, Title: "Test Task", Content: "Test Description"}
	mockRepo.EXPECT().Update(task).Return(nil)

	// Servis fonksiyonunun çağrılması
	err := service.TaskUpdate(task)

	// Hata kontrolü
	assert.NoError(t, err)
}

func TestDefaultTaskService_TaskGetByID_Success(t *testing.T) {
	// Test için hazırlıkları yap
	defer setup(t)()

	// Mock repository'den beklenen değerlerin ayarlanması
	taskID := 1
	fakeTask := models.Task{Id: taskID, Title: "Test Task", Content: "Test Description"}
	mockRepo.EXPECT().GetByID(taskID).Return(fakeTask, nil)

	// Servis fonksiyonunun çağrılması
	task, err := service.TaskGetByID(taskID)

	// Hata kontrolü
	assert.NoError(t, err)

	// Sonuç kontrolü
	assert.Equal(t, fakeTask, task)
}
