package gorobokassa

import (
	"testing"
)

func TestNotifications(t *testing.T) {
	formParameters := map[string]string{
		"Fee": "5.0",
		"IncCurrLabel": "BANK0CEAN3R",
		"IncSum": "150.50",
		"InvId": "2",
		"IsTest": "1",
		"OutSum": "119",
		"PaymentMethod": "BankCard",
		"SHP_user_id": "1",
		"SignatureValue": "C56C15AF62D5A2ED9767081D855E9BAE",
		"crc": "C56C15AF62D5A2ED9767081D855E9BAE",
	}

	notification, err := robokassa.ParseNotification(formParameters)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if notification.Sum != "119" {
		t.Fatalf("sum must be = 119")
	}

	if notification.ID != 2 {
		t.Fatalf("id must be 2")
	}

	if value, ok := notification.Get("user_id"); !ok || value != "1" {
		t.Fatalf("data user_id must be = 1")
	}
}
