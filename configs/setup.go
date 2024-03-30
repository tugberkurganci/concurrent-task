package configs

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

var db *sql.DB

func ConnectDB() *sql.DB {
	// PostgreSQL bağlantı dizesini al
	dbURI := EnvPostgresURI()

	// PostgreSQL veritabanına bağlan
	conn, err := sql.Open("postgres", dbURI)
	if err != nil {
		log.Fatalf("Veritabanına bağlanırken hata oluştu: %v\n", err)
	}

	// Bağlantıyı test et
	err = conn.Ping()
	if err != nil {
		log.Fatalf("Veritabanına ping atılırken hata oluştu: %v\n", err)
	}

	createTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT,
		status BOOLEAN
	)
`

	// Tabloyu oluştur
	_, err = conn.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Tablo oluşturulurken hata oluştu: %v\n", err)
	}

	createUserTableSQL := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) NOT NULL,
		email VARCHAR(100) UNIQUE NOT NULL,
		password VARCHAR(100) NOT NULL
	)
`
	_, err = conn.Exec(createUserTableSQL)
	if err != nil {
		log.Fatalf("users tablosu oluşturulurken hata oluştu: %v\n", err)
	}
	// Bağlantı başarıyla oluşturuldu, geri döndür
	return conn
}

// EnvPostgresURI, PostgreSQL bağlantı dizesini döndürür
func EnvPostgresURI() string {
	// PostgreSQL bağlantı dizesini oluştur
	host := "postgres"
	port := "5432"
	user := "postgres"
	password := "test"
	dbname := "konzek"

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
}

func GetDB() *sql.DB {
	if db == nil {
		db = ConnectDB()
	}
	return db
}
