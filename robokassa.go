package gorobokassa

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// type method string
//
// const (
// 	MethodMD5 method = "md5"
// )

var (
	ErrNoRequiredParameters = errors.New("not all required parameters were provided")
	ErrBadParameterFormat   = errors.New("bad parameter format")
	ErrBadSignature         = errors.New("bad signature")

	requiredNotificationParams = []string{
		"OutSum", "InvId", "Fee", "SignatureValue",
	}
)

type Robokassa struct {
	parameters Parameters
}

type Parameters struct {
	// Индетификатор магазина.
	MerchantLogin string
	// Метод рассчета контрольной суммы.
	// Method method
	// Пароль #1.
	Password1 string
	// Пароль #2.
	Password2 string

	// Тестовая среда?
	Test bool
}

func New(parameters Parameters) *Robokassa {
	return &Robokassa{
		parameters: parameters,
	}
}

func (r *Robokassa) NewPayment(pp PaymentParameters) *Payment {
	return &Payment{
		Robokassa:         r,
		PaymentParameters: &pp,
		Data:              map[string]interface{}{},
	}
}

func (r *Robokassa) ParseNotification(formParameters map[string]string) (*Notification, error) {
	for _, v := range requiredNotificationParams {
		_, ok := formParameters[v]
		if !ok {
			return nil, ErrNoRequiredParameters
		}
	}

	id, err := strconv.ParseUint(formParameters["InvId"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: InvId", ErrBadParameterFormat)
	}

	data := make(map[string]string, 0)
	for k, v := range formParameters {
		// FIXME:
		if len(k) > 4 && strings.HasPrefix(strings.ToLower(k), "shp_") {
			data[k[4:]] = v
		}
	}

	notification := &Notification{
		Robokassa: r,
		Sum:       formParameters["OutSum"],
		ID:        uint32(id),
		Fee:       formParameters["Fee"],
		Data:      data,
	}

	if email, ok := formParameters["EMail"]; ok {
		notification.Email = &email
	}

	if method, ok := formParameters["PaymentMethod"]; ok {
		notification.Method = &method
	}

	if methodLabel, ok := formParameters["IncCurrLabel"]; ok {
		notification.MethodLabel = &methodLabel
	}

	if incSum, ok := formParameters["IncSum"]; ok {
		notification.IncSum = &incSum
	}

	if notification.Signature() != formParameters["SignatureValue"] {
		return nil, ErrBadSignature
	}

	return notification, nil
}
