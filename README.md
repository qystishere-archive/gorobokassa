# Go Robokassa
[![GoDoc](http://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/github.com/qystishere/gorobokassa)

<hr>

### Использование

```go
package main

import (
	"fmt"

	"github.com/qystishere/gorobokassa"
)

func main() {
	shop := gorobokassa.New(gorobokassa.Parameters{
		// Идентификатор магазина.
		MerchantLogin: "shop",
		// Пароль №1.
		Password1: "password1",
		// Пароль №2.
		Password2: "password",
		// Тестовая среда? (пароли тоже нужно использовать тестовые в этом случае)
		Test: true,
	})

	id := uint32(1)
	// Создание платежа
	payment := shop.NewPayment(gorobokassa.PaymentParameters{
		// Сумма в рублях.
		Sum: 100.50,
		// Комментарий к платежу.
		Description: "description",
		// ID платежа в системе магазина.
		ID: &id,
	})
	fmt.Println(payment.Signature())

	// Входящие POST данные на ResultURL.
	formParams := make(map[string]string, 0)

	// Разбор уведомления на ResultURL.
	notification, err := shop.ParseNotification(formParams)
	fmt.Println(notification, err)
}

```