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

func (r *Robokassa) ParseNotification(formParams map[string]string) (*Notification, error) {
	for _, v := range requiredNotificationParams {
		_, ok := formParams[v]
		if !ok {
			return nil, ErrNoRequiredParameters
		}
	}

	sum, err := strconv.ParseFloat(formParams["OutSum"], 32)
	if err != nil {
		return nil, fmt.Errorf("%w: OutSum", ErrBadParameterFormat)
	}

	id, err := strconv.ParseUint(formParams["InvId"], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("%w: InvId", ErrBadParameterFormat)
	}

	fee, err := strconv.ParseFloat(formParams["Fee"], 32)
	if err != nil {
		return nil, fmt.Errorf("%w: Fee", ErrBadParameterFormat)
	}

	data := make(map[string]string, 0)
	for k, v := range formParams {
		// FIXME:
		if len(k) > 4 && strings.HasPrefix(strings.ToLower(k), "shp_") {
			data[k[4:]] = v
		}
	}

	notification := &Notification{
		Robokassa: r,
		Sum:       float32(sum),
		ID:        uint32(id),
		Fee:       float32(fee),
		Data:      data,
	}

	if email, ok := formParams["EMail"]; ok {
		notification.Email = &email
	}

	if method, ok := formParams["IncCurrLabel"]; ok {
		notification.Method = &method
	}

	if incSum, ok := formParams["IncSum"]; ok {
		incSum, err := strconv.ParseFloat(incSum, 32)
		if err != nil {
			return nil, fmt.Errorf("%w: IncSum", ErrBadParameterFormat)
		}
		incSum32 := float32(incSum)
		notification.IncSum = &incSum32
	}

	if notification.Signature() != formParams["SignatureValue"] {
		return nil, ErrBadSignature
	}

	return notification, nil
}
