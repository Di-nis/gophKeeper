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
func (c *Cli) Parser() {
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
		c.Method = method
		c.Item = item
		c.Values = values
		return
	}

	// Иначе используем интерактивный ввод
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Введите на выбор метод отправки, %s, %s, %s или %s: ", MethodAdd, MethodGet, MethodDelete, MethodSync)
	method, _ = reader.ReadString('\n')
	method = strings.TrimSpace(method)

	fmt.Printf("Введите тип данных, %s, %s, %s или %s: ", ItemCredentials, ItemPaymentCard, ItemBinary, ItemText)
	item, _ = reader.ReadString('\n')
	item = strings.TrimSpace(item)

	fmt.Print("Введите значения через запятую без пробелов: ")
	var value string
	value, _ = reader.ReadString('\n')
	value = strings.TrimSpace(value)

	c.Method = method
	c.Item = item
	c.Values = strings.Split(value, ",")
}

// Print - вывод результатов запроса.
func (c *Cli) Print() {
	switch c.Method {
	case MethodAdd:
		fmt.Printf("данные сохранены, хэш %s - %s", c.Item, c.Alias)
	case MethodGet:
		fmt.Printf("%s: %s", c.Item, strings.Join(c.Values, ", "))
	case MethodDelete:
		fmt.Printf("данные %s удалены", c.Item)
	case MethodSync:
		fmt.Printf("данные %s синхронизированы", c.Item)
	}
}
