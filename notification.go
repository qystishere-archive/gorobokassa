package gorobokassa

import (
	"crypto/md5"
	"fmt"
	"sort"
	"strings"
)

type Notification struct {
	*Robokassa
	// Сумма, оплаченная покупателем.
	Sum string
	// Номер платежа в магазине.
	ID uint32
	// Комиссия ROBOKASSA за проведение операции. Для физ. лиц всегда 0.
	Fee string

	/*
		Опционально
	*/

	// Email покупателя.
	Email *string
	// Метод платежа.
	Method *string
	// Конкретный метод платежа. (банк?)
	MethodLabel *string
	// Сумма с учетом комиссии.
	IncSum *string

	// Пользовательские параметры.
	Data map[string]string
}

func (n *Notification) Get(name string) (value string, ok bool) {
	value, ok = n.Data[name]
	return
}

func (n *Notification) Signature() string {
	parameters := make([]string, 0)
	parameters = append(parameters,
		n.Sum,
		fmt.Sprintf("%d", n.ID),
		n.Robokassa.parameters.Password2,
	)

	data := make([]string, 0)
	for k, v := range n.Data {
		data = append(data, fmt.Sprintf("SHP_%s=%s", k, v))
	}
	sort.Strings(data)
	for _, raw := range data {
		parameters = append(parameters, raw)
	}

	return strings.ToUpper(fmt.Sprintf(
		"%x",
		md5.Sum([]byte(
			strings.Join(parameters, ":"),
		)),
	))
}
