package models

import "time"

// Role определяет уровень доступа пользователя.
// Два значения: user (обычный) и admin (тренер с расширенными правами).
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// User представляет зарегистрированного пользователя.
type User struct {
	ID           int       `json:"id"         db:"id"`
	Login        string    `json:"login"      db:"login"`
	Email        string    `json:"email"      db:"email"`
	PasswordHash string    `json:"-"          db:"password_hash"` // не сериализуется в JSON
	Role         Role      `json:"role"       db:"role"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
