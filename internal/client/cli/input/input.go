package input

import (
	"bufio"
	"flag"
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
		method, item string
		values       []string
	)

	flag.StringVar(&method, "method", "", "method request")
	flag.StringVar(&item, "item", "", "item type")
	flag.Func("values", "Comma-separated list of values", func(s string) error {
		values = strings.Split(s, ", ")
		return nil
	})

	flag.Parse()

	if method != "" && len(values) > 0 {
		c.Method = cli.Method(method)
		c.Item = cli.Item(item)
		c.Values = values
		return
	}

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

	if !slices.Contains([]cli.Method{cli.MethodRegister, cli.MethodLogin}, cli.Method(method)) {
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

	fmt.Print("Введите значения через запятую: ")
	var value string
	value, _ = reader.ReadString('\n')
	value = strings.TrimSpace(value)

	c.Method = cli.Method(method)
	c.Item = cli.Item(item)
	c.Values = strings.Split(value, ", ")
}
