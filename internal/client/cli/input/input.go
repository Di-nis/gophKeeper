// Package input - модуль для работы с вводом данных.
package input

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/Di-nis/gophKeeper/internal/client/cli"
)

// Input - структура для хранения аргументов командной строки.
type Input struct {
	Method cli.Method
	Item   cli.Item
	Values []string
}

// New - конструктор Input.
func New() Input {
	return Input{}
}

// Parser - парсер аргументов командной строки.
func (c *Input) Parser() {
	var (
		method, item, value string
	)
	// Иначе используем интерактивный ввод
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf(
		"Введите на выбор метод отправки или зарегистрируйтесь или авторизуйтесь, %s, %s, %s, %s или %s или %s: ",
		cli.MethodAdd,
		cli.MethodGet,
		cli.MethodDelete,
		cli.MethodSync,
		cli.MethodRegister,
		cli.MethodLogin,
	)
	method, _ = reader.ReadString('\n')
	method = strings.TrimSpace(method)

	if slices.Contains([]cli.Method{cli.MethodAdd, cli.MethodGet, cli.MethodDelete}, cli.Method(method)) {
		fmt.Printf(
			"Введите тип данных, %s, %s, %s или %s: ",
			cli.ItemCredentials,
			cli.ItemPaymentCard,
			cli.ItemBinary,
			cli.ItemText,
		)
		item, _ = reader.ReadString('\n')
		item = strings.TrimSpace(item)
	}

	if slices.Contains([]cli.Method{cli.MethodAdd, cli.MethodGet, cli.MethodDelete, cli.MethodRegister, cli.MethodLogin}, cli.Method(method)) {
		fmt.Print("Введите значения через запятую: ")
		value, _ = reader.ReadString('\n')
		value = strings.TrimSpace(value)
	}

	c.Method = cli.Method(method)
	c.Item = cli.Item(item)
	c.Values = strings.Split(value, ", ")
}
