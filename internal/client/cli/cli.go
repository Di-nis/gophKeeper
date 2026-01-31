package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

// Cli - структура для хранения аргументов командной строки.
type Cli struct {
	Method string
	Item   string
	// TODO: сомнительно
	Alias  string
	Values []string
}

// New - конструктор.
func New() Cli {
	return Cli{}
}

// Parser - парсер аргументов командной строки.
func (с *Cli) Parser() {
	var (
		method, item string
		values       []string
	)

	flag.StringVar(&method, "method", "", "method request")
	flag.StringVar(&item, "item", "", "item type")
	flag.Func("values", "Comma-separated list of values", func(s string) error {
		values = strings.Split(s, ",")
		return nil
	})

	flag.Parse()

	if method != "" && item != "" {
		с.Method = method
		с.Item = item
		с.Values = values
		return
	}

	// Иначе используем интерактивный ввод
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите метод отправки: ")
	method, _ = reader.ReadString('\n')
	method = strings.TrimSpace(method)

	fmt.Print("Введите тип данных: ")
	item, _ = reader.ReadString('\n')
	item = strings.TrimSpace(item)

	fmt.Print("Введите значения через запятую без пробелов: ")
	var value string
	value, _ = reader.ReadString('\n')
	value = strings.TrimSpace(value)

	с.Method = method
	с.Item = item
	с.Values = strings.Split(value, ",")
}

// Print - вывод результатов запроса.
func (c *Cli) Print() {
	switch c.Method {
	case MethodAdd:
		fmt.Printf("alias %s - %s", c.Item, c.Alias)
	case MethodGet:
		fmt.Printf("%s: %s", c.Item, strings.Join(c.Values, ", "))
	case MethodDelete:
		fmt.Printf("данные %s удалены", c.Item)
	case MethodSync:
		fmt.Printf("данные %s синхронизированы", c.Item)
	}
}
