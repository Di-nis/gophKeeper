package model

import "time"

// UserID - тип для UUID пользователя.
type UserID string

// Role — тип данных для роли пользователя.
type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

// Auth - модель авторизации.
type Auth struct {
	ID           UserID    `db:"id"`
	Login        string    `json:"login" db:"login"`
	Password     string    `json:"password" db:"password"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// // User - пользователь.
// type User struct {
// 	ID        UserID    `db:"id"`
// 	CreatedAt time.Time `db:"created_at"`
// 	UpdatedAt time.Time `db:"updated_at"`
// }

// Credentials - учетные данные (логин/пароль).
type Credentials struct {
	UUID     UserID `json:"-" db:"user_id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
	Alias    string `json:"alias" db:"alias"`
	Info     string `json:"info" db:"info"`
}

// PaymentCard - данные банковской карты.
type PaymentCard struct {
	UUID     UserID `json:"-" db:"user_id"`
	Number   string `json:"number" db:"number"`
	ExpMonth int32  `json:"exp_month" db:"exp_month"`
	ExpYear  int32  `json:"exp_year" db:"exp_year"`
	CVV      string `json:"cvv" db:"cvv"`
	Holder   string `json:"holder" db:"holder"`
	Alias    string `json:"alias" db:"alias"`
	Info     string `json:"info" db:"info"`
}

// Binary - бинарные данные.
type Binary struct {
	UUID  UserID `json:"-" db:"user_id"`
	Data  []byte `json:"data" db:"data"`
	Alias string `json:"alias" db:"alias"`
	Info  string `json:"info" db:"info"`
}

// Text - произвольные текстовые данные.
type Text struct {
	UUID  UserID `json:"-" db:"user_id"`
	Data  string `json:"data" db:"data"`
	Alias string `json:"alias" db:"alias"`
	Info  string `json:"info" db:"info"`
}
