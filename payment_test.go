package gorobokassa

import (
	"testing"
)

func TestPayments(t *testing.T) {
	id := uint32(123)
	email := "123@123.ru"
	ip := "127.0.0.1"

	payment := robokassa.NewPayment(PaymentParameters{
		Sum:         100.10,
		ID:          &id,
		Email:       &email,
		IP:          &ip,
	})
	if payment.Signature() != "A759441945E710BAF2EE24FBF5738506" {
		t.Fatalf("bad signature")
	}

	payment.Set("user_id", 1)
	if payment.Signature() != "BB5688776EEC32D6FA73AE96D55704CB" {
		t.Fatalf("bad signature")
	}
}
