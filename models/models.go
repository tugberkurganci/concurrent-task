package models

type Task struct {
	Id      int    `json:"id,omitempty" `
	Title   string `json:"title,omitempty" validate:"required,min=2"`
	Content string `json:"content,omitempty" validate:"required,min=2"`
	Status  bool   `json:"status,omitempty" validate:"required"`
}
type User struct {
	ID       int64  `gorm:"primary_key:auto_increment" json:"-"`
	Name     string `gorm:"type:varchar(100)" json:"name,omitempty" validate:"required,min=2"`
	Email    string `gorm:"type:varchar(100);unique;" json:"email,omitempty" validate:"required,email"`
	Password string `gorm:"type:varchar(100)" json:"password,omitempty" validate:"required,min=6"`
}
