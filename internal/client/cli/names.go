// Package cli - ввод/вывод данных.
package cli

// Method - тип данных для имен методов.
type Method string

// Методы.
const (
	MethodRegister Method = "register"
	MethodLogin    Method = "login"
	MethodGet      Method = "get"
	MethodAdd      Method = "add"
	MethodDelete   Method = "delete"
	MethodSync     Method = "sync"
)

// item - тип данных для имен данных.
type Item string

// Типы данных.
const (
	ItemCredentials Item = "credentials"
	ItemPaymentCard Item = "payment_card"
	ItemBinary      Item = "binary"
	ItemText        Item = "text"
)
