package output

import (
	"fmt"
	"strings"

	"github.com/Di-nis/gophKeeper/internal/client/cli"
)

// Output - структура для вывода данных.
type Output struct {
	Method cli.Method
	Item   cli.Item
	Values []string
}

// New - конструктор Output.
func New() Output {
	return Output{}
}

// SetMethod устанавливает значение поля Method.
func (o *Output) SetMethod(method cli.Method) *Output {
	o.Method = method
	return o
}

// SetItem устанавливает значение поля Item.
func (o *Output) SetItem(item cli.Item) *Output {
	o.Item = item
	return o
}

// SetValues устанавливает значение поля Values.
func (o *Output) SetValues(values ...string) *Output {
	o.Values = append(o.Values, values...)
	return o
}

// PrintSuccess - вывод успешных результатов запроса.
func (o *Output) PrintSuccess() {
	switch o.Method {
	case cli.MethodRegister:
		fmt.Print("вы успешно зарегистрированы")
	case cli.MethodLogin:
		fmt.Print("вы успешно авторизованы")
	case cli.MethodAdd:
		fmt.Printf("данные сохранены, хэш %s - %s", o.Item, o.Values[0])
	case cli.MethodGet:
		fmt.Printf("%s: %s", o.Item, strings.Join(o.Values, ", "))
	case cli.MethodDelete:
		fmt.Printf("данные %s удалены", o.Item)
	case cli.MethodSync:
		fmt.Printf("данные %s синхронизированы", o.Item)
	}
}

// PrintError - вывод ошибок.
func (o *Output) PrintError() {
	fmt.Print("ошибка выполнения, повторите попыптку")
}
