package repository_test

import (
	"database/sql"
	"log"
	"strconv"
	"testing"

	"konzek-jun/models"
	"konzek-jun/repository"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository(t *testing.T) {
	// Veritabanı bağlantısını oluştur
	db, err := sql.Open("postgres", "dbname=konzek user=postgres password=test host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata oluştu: %v", err)
	}
	defer db.Close()

	// Veritabanını temizle
	clearDatabasev2(db)

	// UserRepositoryDb'yi oluştur
	userRepo := repository.NewUserRepo(db)

	// Insert metodunu test et
	t.Run("Insert", func(t *testing.T) {
		// Test için gerekli örnek kullanıcı verisini oluştur
		user := models.User{Name: "Test User", Email: "test@example.com", Password: "testpass"}
		insertedUser, err := userRepo.InsertUser(user)
		if err != nil {
			t.Errorf("Kullanıcı eklenirken hata oluştu: %v", err)
		}
		// Insert işlemi başarılıysa, dönen kullanıcının ID'sinin 0'dan büyük olduğunu kontrol et
		assert.Equal(t, user.Name, insertedUser.Name)
		assert.Equal(t, user.Email, insertedUser.Email)
	})

	t.Run("GetByEmail", func(t *testing.T) {
		// Test için bir örnek e-posta adresi belirle
		email := "test@example.com"

		// Test için bir örnek kullanıcı verisi oluştur ve veritabanına kaydet
		user := models.User{Name: "Test User", Email: email, Password: "testpass"}
		_, err := userRepo.InsertUser(user)
		if err != nil {
			t.Errorf("Kullanıcı eklenirken hata oluştu: %v", err)
		}

		// Belirtilen e-posta adresine sahip kullanıcıyı getir ve sonucu kontrol et
		foundUser, err := userRepo.FindByEmail(email)
		if err != nil {
			t.Errorf("E-posta adresine göre kullanıcı getirilirken hata oluştu: %v", err)
		}

		// Kullanıcı bulundu mu?
		assert.NotNil(t, foundUser)
		// Bulunan kullanıcının e-posta adresi beklenen e-posta adresi ile eşleşiyor mu?
		assert.Equal(t, email, foundUser.Email)
	})
	t.Run("GetByUserID", func(t *testing.T) {
		// Test için bir örnek kullanıcı verisi oluştur ve veritabanına kaydet
		user := models.User{Name: "Test User", Email: "test@example.com", Password: "testpass"}
		insertedUser, err := userRepo.InsertUser(user)
		if err != nil {
			t.Errorf("Kullanıcı eklenirken hata oluştu: %v", err)
		}

		// Kaydedilen kullanıcının ID'sini kullanarak kullanıcıyı getir ve sonucu kontrol et
		foundUser, err := userRepo.FindByUserID(strconv.FormatInt(insertedUser.ID, 10))
		if err != nil {
			t.Errorf("Kullanıcı ID'sine göre kullanıcı getirilirken hata oluştu: %v", err)
		}

		// Bulunan kullanıcının ID'si beklenen ID ile eşleşiyor mu?
		if foundUser.ID != insertedUser.ID {
			t.Errorf("Beklenen kullanıcı ID'si ile eşleşmiyor. Beklenen: %d, Alınan: %d", insertedUser.ID, foundUser.ID)
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		// Test için bir örnek kullanıcı verisi oluştur ve veritabanına kaydet
		user := models.User{Name: "Test User", Email: "test@example.com", Password: "testpass"}
		insertedUser, err := userRepo.InsertUser(user)
		if err != nil {
			t.Errorf("Kullanıcı eklenirken hata oluştu: %v", err)
		}

		// Kullanıcının adını ve e-posta adresini güncelle
		insertedUser.Name = "Updated Name"
		insertedUser.Email = "updated@example.com"
		_, err = userRepo.UpdateUser(insertedUser)
		if err != nil {
			t.Errorf("Kullanıcı güncellenirken hata oluştu: %v", err)
		}

		// Güncellenen kullanıcıyı tekrar getir ve sonucu kontrol et
		updatedUser, err := userRepo.FindByUserID(strconv.FormatInt(insertedUser.ID, 10))
		if err != nil {
			t.Errorf("Kullanıcı güncellenirken hata oluştu: %v", err)
		}

		// Kullanıcı bulundu mu?
		assert.NotNil(t, updatedUser)
		// Güncellenen kullanıcının adı ve e-posta adresi beklenen değerlerle eşleşiyor mu?
		assert.Equal(t, "Updated Name", updatedUser.Name)
		assert.Equal(t, "updated@example.com", updatedUser.Email)
	})

	// Diğer testleri buraya ekle...

}

func clearDatabasev2(db *sql.DB) {
	_, err := db.Exec("DELETE FROM users")
	if err != nil {
		log.Fatalf("Veritabanını temizlerken hata oluştu: %v", err)
	}
}
