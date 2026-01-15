package model

import "time"

// UserID - тип для userID пользователя.
type UserID string

// Role — тип данных для роли пользователя.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Auth - auth.
type Auth struct {
	UserID       UserID    `db:"id"`
	Login        string    `db:"login"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// User - пользователь.
type User struct {
	ID        UserID    `db:"id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

// Credentials - учетные данные (логин/пароль).
type Credentials struct {
	UserID   UserID `json:"-" db:"user_id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
	Alias    string `json:"alias" db:"alias"`
	Info     string `json:"info" db:"info"`
}

// PaymentCard - данные банковской карты.
type PaymentCard struct {
	UserID   UserID `json:"-" db:"user_id"`
	Number   string `json:"number"`
	ExpMonth int    `json:"exp_month"`
	ExpYear  int    `json:"exp_year"`
	CVV      string `json:"cvv"`
	Holder   string `json:"holder"`
	Info     string `json:"info" db:"info"`
}

// Binary - бинарные данные.
type Binary struct {
	UserID UserID `json:"-" db:"user_id"`
	Data   []byte `json:"data" db:"data"`
	Info   string `json:"info" db:"info"`
}

// Text - произвольные текстовые данные.
type Text struct {
	UserID UserID `json:"-" db:"user_id"`
	Data   string `json:"data" db:"data"`
	Info   string `json:"info" db:"info"`
}
