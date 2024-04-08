package app

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	services "konzek-jun/mocks/service"
	"konzek-jun/models"
	"konzek-jun/repository"
	x "konzek-jun/services"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

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

func TestTaskStream(t *testing.T) {

	db, err := sql.Open("postgres", "dbname=konzek user=postgres password=test host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata oluştu: %v", err)
	}
	defer db.Close()

	clearDatabase(db)
	taskRepo := repository.NewTaskRepository(db)
	taskService := x.NewTaskService(taskRepo)
	taskHandler := NewTaskHandler(taskService, 5)

	router := fiber.New()
	router.Get("/api/tasks/:id", taskHandler.GetByID)
	router.Post("/api/tasks", taskHandler.CreateTask)
	router.Delete("/api/tasks/:id", taskHandler.DeleteTask)

	// Task oluştur
	task := models.Task{
		Content: "Test Content",
		Status:  true,
		Title:   "xxxxxxxx",
	}
	taskJSON, _ := json.Marshal(task)
	req := httptest.NewRequest(http.MethodPost, "/api/tasks", bytes.NewBuffer(taskJSON))
	req.Header.Set("Content-Type", "application/json")
	resp1, _ := router.Test(req)
	// Bekleme süresi (30 saniye)
	time.Sleep(20 * time.Second)

	req2 := httptest.NewRequest(http.MethodGet, "/api/tasks/8", nil)
	resp2, _ := router.Test(req2)

	req3 := httptest.NewRequest(http.MethodDelete, "/api/tasks/8", nil)
	resp3, _ := router.Test(req3)

	fmt.Println(resp2)
	assert.Equal(t, http.StatusCreated, resp1.StatusCode)
	assert.Equal(t, http.StatusOK, resp2.StatusCode)
	assert.Equal(t, http.StatusOK, resp3.StatusCode)

}
func clearDatabase(db *sql.DB) {
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		log.Fatalf("Veritabanını temizlerken hata oluştu: %v", err)
	}
}
