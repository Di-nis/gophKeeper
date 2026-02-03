package model

import (
	"time"

	"github.com/samborkent/uuidv7"
)

// UserID - тип для UUID пользователя.
type UserID uuidv7.UUID

// String - преобразование к строке.
func (u UserID) String() string {
	return uuidv7.UUID(u).String()
}

// Role — тип данных для роли пользователя.
type Role string

const (
	// RoleUser - роль "user"
	RoleUser Role = "user"
	// RoleAdmin - роль "admin"
	RoleAdmin Role = "admin"
)

// Auth - модель авторизации.
type Auth struct {
	ID           UserID    `db:"id"`
	Login        string    `json:"login" binding:"required" example:"your_username" db:"login"`
	Password     string    `json:"password" binding:"required" example:"your_password" db:"password"`
	PasswordHash string    `db:"password_hash"`
	Role         Role      `db:"role"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

// GetID - возвращает ID пользователя.
func (a *Auth) GetID() UserID {
	return a.ID
}

// GetLogin - возвращает логин пользователя.
func (a *Auth) GetLogin() string {
	return a.Login
}

// GetPasswordHash - возвращает хэш пароля пользователя.
func (a *Auth) GetPasswordHash() string {
	return a.PasswordHash
}

// GetRole - возвращает роль пользователя.
func (a *Auth) GetRole() Role {
	return a.Role
}

// SetID - устанавливает ID пользователя.
func (a *Auth) SetID(id uuidv7.UUID) *Auth {
	a.ID = UserID(id)
	return a
}

// SetLogin - устанавливает логин пользователя.
func (a *Auth) SetLogin(login string) *Auth {
	a.Login = login
	return a
}

// SetPasswordHash - устанавливает хэш пароля пользователя.
func (a *Auth) SetPassword(password string) *Auth {
	a.Password = password
	return a
}

// SetPasswordHash - устанавливает хэш пароля пользователя.
func (a *Auth) SetPasswordHash(passwordHash string) *Auth {
	a.PasswordHash = passwordHash
	return a
}

// SetRole - устанавливает роль пользователя.
func (a *Auth) SetRole(role Role) *Auth {
	a.Role = role
	return a
}

// Credentials - учетные данные (логин/пароль).
type Credentials struct {
	UUID     UserID `json:"-" db:"user_id"`
	Login    string `json:"login" db:"login"`
	Password string `json:"password" db:"password"`
	Alias    string `json:"-" db:"alias"`
	Info     string `json:"info" db:"info"`
}

// GetRaw - возвращает сырые данные.
func (c *Credentials) GetRaw() string { return c.Login + c.Password + c.Info }

// SetAlias - устанавливает алиас.
func (c *Credentials) SetAlias(a string) { c.Alias = a }

// PaymentCard - данные банковской карты.
type PaymentCard struct {
	UUID     UserID `json:"-" db:"user_id"`
	Number   string `json:"number" db:"number"`
	ExpMonth string `json:"exp_month" db:"exp_month"`
	ExpYear  string `json:"exp_year" db:"exp_year"`
	CVV      string `json:"cvv" db:"cvv"`
	Holder   string `json:"holder" db:"holder"`
	Alias    string `json:"alias" db:"alias"`
	Info     string `json:"info" db:"info"`
}

// GetRaw - возвращает сырые данные.
func (p *PaymentCard) GetRaw() string { return p.Number }

// SetAlias - устанавливает алиас.
func (p *PaymentCard) SetAlias(a string) { p.Alias = a }

// Binary - бинарные данные.
type Binary struct {
	UUID  UserID `json:"-" db:"user_id"`
	Data  []byte `json:"data" db:"data"`
	Alias string `json:"alias" db:"alias"`
	Info  string `json:"info" db:"info"`
}

// GetRaw - возвращает сырые данные.
func (b *Binary) GetRaw() string { return string(b.Data) + b.Info }

// SetAlias - устанавливает алиас.
func (b *Binary) SetAlias(a string) { b.Alias = a }

// Text - произвольные текстовые данные.
type Text struct {
	UUID  UserID `json:"-" db:"user_id"`
	Data  string `json:"data" db:"data"`
	Alias string `json:"alias" db:"alias"`
	Info  string `json:"info" db:"info"`
}

// GetRaw - возвращает сырые данные.
func (t *Text) GetRaw() string { return t.Data + t.Info }

// SetAlias - устанавливает алиас.
func (t *Text) SetAlias(a string) { t.Alias = a }

// Common - общие данные.
type Common struct {
	Credentials []*Credentials
	PaymentCard []*PaymentCard
	Binary      []*Binary
	Text        []*Text
}
