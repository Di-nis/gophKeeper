package model

import (
	"fmt"
	"strings"
	"time"
)

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
	Alias    string `json:"-" db:"alias"`
	Info     string `json:"info" db:"info"`
}

func (c *Credentials) GetRaw() string    { return c.Login + c.Password + c.Info }
func (c *Credentials) SetAlias(a string) { c.Alias = a }

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

func (p *PaymentCard) GetRaw() string    { return p.Number }
func (p *PaymentCard) SetAlias(a string) { p.Alias = a }

// Binary - бинарные данные.
type Binary struct {
	UUID  UserID `json:"-" db:"user_id"`
	Data  []byte `json:"data" db:"data"`
	Alias string `json:"alias" db:"alias"`
	Info  string `json:"info" db:"info"`
}

func (b *Binary) GetRaw() string    { return string(b.Data) + b.Info }
func (b *Binary) SetAlias(a string) { b.Alias = a }

// Text - произвольные текстовые данные.
type Text struct {
	UUID  UserID `json:"-" db:"user_id"`
	Data  string `json:"data" db:"data"`
	Alias string `json:"alias" db:"alias"`
	Info  string `json:"info" db:"info"`
}

func (t *Text) GetRaw() string    { return t.Data + t.Info }
func (t *Text) SetAlias(a string) { t.Alias = a }

// Command - команда.
type Command struct {
	Method string
	Item   string
	Value  []string
}

// String должен уметь сериализовать переменную типа в строку.
func (c *Command) String() string {
	return fmt.Sprint(strings.Join(c.Value, ","))
}

// Set связывает переменную типа со значением флага
// и устанавливает правила парсинга для пользовательского типа.
func (c *Command) Set(flagValue string) error {
	c.Value = strings.Split(flagValue, " ")
	return nil
}
