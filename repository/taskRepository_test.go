package repository_test

import (
	"database/sql"
	"log"
	"testing"

	"konzek-jun/models"
	"konzek-jun/repository"

	_ "github.com/lib/pq"
)

func TestTaskRepository(t *testing.T) {

	db, err := sql.Open("postgres", "dbname=konzek user=postgres password=test host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata oluştu: %v", err)
	}
	defer db.Close()

	clearDatabase(db)

	taskRepo := repository.NewTaskRepository(db)

	// Insert metodunu test et
	t.Run("Insert", func(t *testing.T) {
		task := models.Task{Title: "Test Task", Content: "Test Content", Status: true}
		id, err := taskRepo.Insert(task)
		if err != nil {
			t.Errorf("Task eklenirken hata oluştu: %v", err)
		}
		// Insert işlemi başarılıysa, dönen ID'nin 0'dan büyük olduğunu kontrol et
		if id <= 0 {
			t.Errorf("Beklenen ID 0'dan büyük değil, alınan: %d", id)
		}
	})

	// GetAll metodunu test et
	t.Run("GetAll", func(t *testing.T) {
		// Test için örnek veri ekle
		_, err := db.Exec("INSERT INTO tasks (title, content, status) VALUES ($1, $2, $3)", "Test Task 1", "Test Content 1", true)
		if err != nil {
			t.Fatalf("Veritabanına örnek veri eklerken hata oluştu: %v", err)
		}

		// GetAll metodunu test et
		tasks, err := taskRepo.GetAll()
		if err != nil {
			t.Errorf("Task'leri getirirken hata oluştu: %v", err)
		}
		if len(tasks) != 1 {
			t.Errorf("Beklenen task sayısı 1 değil, alınan: %d", len(tasks))
		}
	})

	// Delete metodunu test et
	t.Run("Delete", func(t *testing.T) {
		// Test için örnek veri ekle
		_, err := db.Exec("INSERT INTO tasks (title, content, status) VALUES ($1, $2, $3)", "Test Task 1", "Test Content 1", true)
		if err != nil {
			t.Fatalf("Veritabanına örnek veri eklerken hata oluştu: %v", err)
		}

		// Delete metodunu test et
		err = taskRepo.Delete(1)
		if err != nil {
			t.Errorf("Task silinirken hata oluştu: %v", err)
		}
	})

	// GetByID metodunu test et
	t.Run("GetByID", func(t *testing.T) {
		// Test için örnek veri ekle
		_, err := db.Exec("INSERT INTO tasks (title, content, status) VALUES ($1, $2, $3)", "Test Task 1", "Test Content 1", true)
		if err != nil {
			t.Fatalf("Veritabanına örnek veri eklerken hata oluştu: %v", err)
		}

		// GetByID metodunu test et
		task, err := taskRepo.GetByID(1)
		if err != nil {
			t.Errorf("Task getirilirken hata oluştu: %v", err)
		}
		if task.Id != 0 {
			t.Errorf("Beklenen task ID'si 1 değil, alınan: %d", task.Id)
		}
	})

	// Update metodunu test et
	t.Run("Update", func(t *testing.T) {
		// Test için örnek veri ekle
		_, err := db.Exec("INSERT INTO tasks (title, content, status) VALUES ($1, $2, $3)", "Test Task 1", "Test Content 1", true)
		if err != nil {
			t.Fatalf("Veritabanına örnek veri eklerken hata oluştu: %v", err)
		}

		// Update metodunu test et
		task := models.Task{Id: 1, Title: "Updated Task", Content: "Updated Content", Status: false}
		err = taskRepo.Update(task)
		if err != nil {
			t.Errorf("Task güncellenirken hata oluştu: %v", err)
		}
	})
}

// Test veritabanını temizle
func clearDatabase(db *sql.DB) {
	_, err := db.Exec("DELETE FROM tasks")
	if err != nil {
		log.Fatalf("Veritabanını temizlerken hata oluştu: %v", err)
	}
}
