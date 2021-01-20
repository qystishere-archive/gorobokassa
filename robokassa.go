package gorobokassa

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type method string

const (
	MethodMD5 method = "md5"
)

var (
	ErrNoRequiredParameters = errors.New("not all required parameters were provided")
	ErrBadParameterFormat   = errors.New("bad parameter format")
	ErrBadSignature         = errors.New("bad signature")

	requiredNotificationParams = []string{
		"OutSum", "InvId", "EMail", "Fee", "SignatureValue",
	}
)

type Robokassa struct {
	parameters Parameters
}

type Parameters struct {
	// Индетификатор магазина.
	MerchantLogin string
	// Метод рассчета контрольной суммы.
	Method method
	// Пароль #1.
	Password1 string
	// Пароль #2.
	Password2 string
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

	data := make(map[string]interface{}, 0)
	for k, v := range formParams {
		if strings.HasPrefix(strings.ToLower(k), "shp_") {
			data[k] = v
		}
	}

	notification := &Notification{
		Robokassa: r,
		Sum:       float32(sum),
		ID:        uint32(id),
		Email:     formParams["EMail"],
		Fee:       float32(fee),
		Data:      data,
	}

	if notification.Signature() != formParams["SignatureValue"] {
		return nil, ErrBadSignature
	}

	return notification, nil
}
