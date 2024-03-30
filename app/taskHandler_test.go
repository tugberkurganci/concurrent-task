package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	services "konzek-jun/mocks/service"
	"konzek-jun/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

var mockService *services.MockTaskService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = services.NewMockTaskService(ctrl)

	return func() { defer ctrl.Finish() }
}

func TestTaskHandler_GetAllTask(t *testing.T) {

	trd := setup(t)
	defer trd()
	td := NewTaskHandler(mockService, 5)
	router := fiber.New()
	router.Get("/api/tasks", td.GetAllTask)
	var FakeData = []models.Task{
		{Id: 1, Title: "Task 1", Content: "Description 1", Status: true},
		{Id: 2, Title: "Task 2", Content: "Description 2", Status: false},
		{Id: 3, Title: "Task 3", Content: "Description 3", Status: true},
	}

	mockService.EXPECT().TaskGetAll().Return(FakeData, nil)

	req := httptest.NewRequest("GET", "/api/tasks", nil)

	resp, _ := router.Test(req, 1)

	assert.Equal(t, 200, resp.StatusCode)
}

// @Summary Create a new task
// @Description Creates a new task
// @Tags Tasks
// @Accept json
// @Produce json
// @Param task body models.Task true "Task object to create"
// @Success 201 {object} map[string]interface{} "Empty response"
// @Failure 400 {object} globalerror.ErrorResponse "Bad request"
// @Failure 500 {object} globalerror.ErrorResponse "Internal server error"
// @Router /tasks [post]
func TestTaskHandler_CreateTask(t *testing.T) {
	trd := setup(t)
	defer trd()

	td := NewTaskHandler(mockService, 5)
	router := fiber.New()
	router.Post("/api/tasks", td.CreateTask)

	mockService.EXPECT().TaskInsert(gomock.Any()).Return(nil)

	task := models.Task{Title: "Test Task", Content: "Test Content", Status: true}

	jsonData, err := json.Marshal(task)
	if err != nil {
		t.Fatalf("JSON formatına dönüştürme hatası: %v", err)
	}

	jsonReader := bytes.NewReader(jsonData)

	req := httptest.NewRequest("POST", "/api/tasks", jsonReader)

	req.Header.Set("Content-Type", "application/json")

	resp, err := router.Test(req)
	if err != nil {
		t.Fatalf("İstek gerçekleştirilirken bir hata oluştu: %v", err)
	}

	if resp.StatusCode != 201 {
		t.Errorf("Beklenen durum kodu 201 değil, alınan: %d", resp.StatusCode)
	}
}

func TestTaskHandler_DeleteTask(t *testing.T) {
	trd := setup(t)
	defer trd()
	td := NewTaskHandler(mockService, 5)
	router := fiber.New()
	router.Delete("/api/tasks/:id", td.DeleteTask)
	mockService.EXPECT().TaskDelete(gomock.Any()).Return(nil)
	req := httptest.NewRequest("DELETE", "/api/tasks/1", nil)
	resp, _ := router.Test(req, 1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestTaskHandler_GetByID(t *testing.T) {
	trd := setup(t)
	defer trd()
	td := NewTaskHandler(mockService, 5)
	router := fiber.New()
	router.Get("/api/tasks/:id", td.GetByID)
	mockService.EXPECT().TaskGetByID(1).Return(models.Task{Id: 1, Title: "Test Task", Content: "Test Content", Status: true}, nil)
	req := httptest.NewRequest("GET", "/api/tasks/1", nil)
	resp, _ := router.Test(req, 1)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestUpdateTaskHandler(t *testing.T) {
	trd := setup(t)
	defer trd()

	td := NewTaskHandler(mockService, 5)
	mockService.EXPECT().TaskUpdate(gomock.Any()).Return(nil)
	router := fiber.New()
	router.Put("/api/tasks", td.UpdateTask)
	task := models.Task{
		Id:      1,
		Title:   "Updated Task Title",
		Content: "Updated Task Content",
		Status:  true,
	}

	taskJSON, err := json.Marshal(task)
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest("PUT", "/api/tasks", bytes.NewReader(taskJSON))
	req.Header.Set("Content-Type", "application/json")

	resp, err := router.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusOK, resp.StatusCode)

	fmt.Println("Test başarılı. Geçti mesajı alındı.")
}
