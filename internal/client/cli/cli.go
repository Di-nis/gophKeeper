package cli

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/Di-nis/gophKeeper/internal/model"
)

// Parser - парсер аргументов командной строки.
func Parser() *model.Command {
	var (
		method, item, value string
		command             *model.Command
	)

	flag.StringVar(&method, "method", "", "method request")
	flag.StringVar(&item, "item", "", "item type")
	flag.StringVar(&value, "value", "", "values")

	flag.Parse()

	// Если аргументы CLI переданы, используем их
	if method != "" && item != "" {
		command = &model.Command{
			Method: method,
			Item:   item,
			Value:  strings.Split(value, "_"),
		}
		return command
	}

	// Иначе используем интерактивный ввод
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Введите метод отправки: ")
	method, _ = reader.ReadString('\n')
	method = strings.TrimSpace(method)

	fmt.Print("Введите тип данных: ")
	item, _ = reader.ReadString('\n')
	item = strings.TrimSpace(item)

	fmt.Print("Введите значения через пробел: ")
	value, _ = reader.ReadString('\n')
	value = strings.TrimSpace(value)

	command = &model.Command{
		Method: method,
		Item:   item,
		Value:  strings.Split(value, " "),
	}
	return command
}
