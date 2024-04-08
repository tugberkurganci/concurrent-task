package models

type Task struct {
	Id      int    `json:"id,omitempty" `
	Title   string `json:"title,omitempty" validate:"required,min=2"`
	Content string `json:"content,omitempty" validate:"required,min=2"`
	Status  bool   `json:"status,omitempty" validate:"required"`
}
type User struct {
	ID       int64  `json:"-"`
	Name     string `json:"name,omitempty" validate:"required,min=2"`
	Email    string ` json:"email,omitempty" validate:"required,email"`
	Password string ` json:"password,omitempty" validate:"required,min=6"`
}
