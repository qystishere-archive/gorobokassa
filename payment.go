package gorobokassa

import (
	"crypto/md5"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Payment struct {
	*Robokassa
	*PaymentParameters

	// Пользовательские параметры, что будут переданы магазину.
	Data map[string]interface{}
}

type PaymentParameters struct {
	// Сумма платежа в рублях.
	Sum float32
	// Описание платежа.
	Description string

	/*
		Опционально
	*/

	// Способ оплаты.
	Method *string
	// Номер платежа в системе магазина.
	ID *uint32
	// Язык общения с клиентом.
	Language *string
	// Кодировка страницы.
	Encoding *string
	// Email покупателя.
	Email *string
	// Срок, до которого действует счет.
	ExpireAt *time.Time
	// Курс, по которому считать оплату. (USD, EUR, KZT)
	Currency *string
	// IP пользователя.
	IP *string
}

func (p *Payment) Set(name string, value interface{}) {
	p.Data[name] = value
}

func (p *Payment) Get(name string) (value interface{}, ok bool) {
	value, ok = p.Data[name]
	return
}

func (p *Payment) Signature() string {
	parameters := make([]string, 0)
	parameters = append(parameters,
		p.Robokassa.parameters.MerchantLogin,
		fmt.Sprintf("%.2f", p.Sum),
	)
	if p.ID != nil {
		parameters = append(parameters, fmt.Sprintf("%d", *p.ID))
	}
	if p.Currency != nil {
		parameters = append(parameters, *p.Currency)
	}
	if p.IP != nil {
		parameters = append(parameters, *p.IP)
	}
	parameters = append(parameters, p.Robokassa.parameters.Password1)

	data := make([]string, 0)
	for k, v := range p.Data {
		data = append(data, fmt.Sprintf(
			"SHP_%s=%s",
			k,
			url.QueryEscape(fmt.Sprintf("%v", v)),
		))
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

func (p *Payment) QueryURL() string {
	var sb strings.Builder
	sb.WriteString(
		fmt.Sprintf(
			"?MerchantLogin=%s&OutSum=%.2f&Description=%s&SignatureValue=%s",
			p.Robokassa.parameters.MerchantLogin,
			p.Sum,
			p.Description,
			p.Signature(),
		),
	)
	if p.Method != nil {
		sb.WriteString(fmt.Sprintf("&IncCurrLabel=%s", *p.Method))
	}
	if p.ID != nil {
		sb.WriteString(fmt.Sprintf("&InvId=%d", *p.ID))
	}
	if p.Language != nil {
		sb.WriteString(fmt.Sprintf("&Culture=%s", *p.Language))
	}
	if p.Encoding != nil {
		sb.WriteString(fmt.Sprintf("&Encoding=%s", *p.Encoding))
	}
	if p.Email != nil {
		sb.WriteString(fmt.Sprintf("&Email=%s", *p.Email))
	}
	if p.ExpireAt != nil {
		sb.WriteString(
			fmt.Sprintf("&ExpirationDate=%s", (*p.ExpireAt).Format("2006-01-02 15:04:05")),
		)
	}
	if p.Currency != nil {
		sb.WriteString(fmt.Sprintf("&OutSumCurrency=%s", *p.Currency))
	}
	if p.IP != nil {
		sb.WriteString(fmt.Sprintf("&UserIp=%s", *p.IP))
	}

	data := make([]string, 0)
	for k, v := range p.Data {
		vs := url.QueryEscape(fmt.Sprintf("%v", v))
		data = append(data, fmt.Sprintf("SHP_%s=%s", k, vs))
	}
	sort.Strings(data)
	for _, raw := range data {
		sb.WriteString(fmt.Sprintf("&%s", raw))
	}

	if p.Robokassa.parameters.Test {
		sb.WriteString("&IsTest=1")
	}

	return sb.String()
}
