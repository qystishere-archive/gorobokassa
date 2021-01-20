package gorobokassa

import (
	"os"
	"testing"
)

var robokassa *Robokassa

func TestMain(m *testing.M) {
	robokassa = New(Parameters{
		MerchantLogin: "login",
		Password1:     "password1",
		Password2:     "password2",
	})

	os.Exit(m.Run())
}
