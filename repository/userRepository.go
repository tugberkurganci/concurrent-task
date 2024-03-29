package repository

import (
	"database/sql"
	"fmt"
	"konzek-jun/loggerx"
	"konzek-jun/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=../mocks//repository/mockUserrepository.go -package=repository konzek-jun/repository UserRepository
type UserRepository interface {
	InsertUser(user models.User) (models.User, error)
	UpdateUser(user models.User) (models.User, error)
	FindByEmail(email string) (models.User, error)
	FindByUserID(userID string) (models.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{
		db: db,
	}
}

func (ur *userRepo) InsertUser(user models.User) (models.User, error) {
	user.Password = hashAndSalt([]byte(user.Password))
	err := ur.db.QueryRow("INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id", user.Name, user.Email, user.Password).Scan(&user.ID)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while inserting user: %v", err))
		return models.User{}, err
	}
	loggerx.Info("User inserted successfully")
	return user, nil
}

func (ur *userRepo) UpdateUser(user models.User) (models.User, error) {
	if user.Password != "" {
		user.Password = hashAndSalt([]byte(user.Password))
	} else {
		var tempUser models.User
		err := ur.db.QueryRow("SELECT password FROM users WHERE id = $1", user.ID).Scan(&tempUser.Password)
		if err != nil {
			loggerx.Error(fmt.Sprintf("Error while updating user: %v", err))
			return models.User{}, err
		}
		user.Password = tempUser.Password
	}

	_, err := ur.db.Exec("UPDATE users SET name = $1, email = $2, password = $3 WHERE id = $4", user.Name, user.Email, user.Password, user.ID)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while updating user: %v", err))
		return models.User{}, err
	}
	loggerx.Info("User updated successfully")
	return user, nil
}

func (ur *userRepo) FindByEmail(email string) (models.User, error) {
	var user models.User
	err := ur.db.QueryRow("SELECT id, name, email, password FROM users WHERE email = $1", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while finding user by email: %v", err))
		return models.User{}, err
	}
	loggerx.Info("User found by email successfully")
	return user, nil
}

func (ur *userRepo) FindByUserID(userID string) (models.User, error) {
	var user models.User
	err := ur.db.QueryRow("SELECT id, name, email, password FROM users WHERE id = $1", userID).Scan(&user.ID, &user.Name, &user.Email, &user.Password)
	if err != nil {
		loggerx.Error(fmt.Sprintf("Error while finding user by ID: %v", err))
		return models.User{}, err
	}
	loggerx.Info("User found by ID successfully")
	return user, nil
}

func hashAndSalt(pwd []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
		loggerx.Error(fmt.Sprintf("Error while hashing password: %v", err))
		panic(err)
	}
	return string(hash)
}
