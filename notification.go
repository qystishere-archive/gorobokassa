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
	Sum float32
	// Номер платежа в магазине.
	ID uint32
	// Email покупателя.
	Email string
	// Комиссия ROBOKASSA за проведение операции. Для физ. лиц всегда 0.
	Fee float32

	// Пользовательские параметры.
	Data map[string]interface{}
}

func (n *Notification) Signature() string {
	parameters := make([]string, 0)
	parameters = append(parameters,
		fmt.Sprintf("%.2f", n.Sum),
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

	return fmt.Sprintf(
		"%x",
		md5.Sum([]byte(
			strings.Join(parameters, ":"),
		)),
	)
}