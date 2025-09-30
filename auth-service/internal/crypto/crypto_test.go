package crypto

import (
	"testing"
)

func TestVerifyPassword(t *testing.T) {
	password := "secure_password123"
	p := GetDefaultParams()

	encodedHash, err := HashPassword(password, &p)
	if err != nil {
		t.Errorf("Ran into error while testing: %s", err)
	}

	got, err := VerifyPassword(password, encodedHash)
	if err != nil {
		t.Errorf("Ran into error while testing: %s", err)
	}

	want := true

	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

	got, err = VerifyPassword(password+"wrong", encodedHash)
	if err != nil {
		t.Errorf("Ran into error while testing: %s", err)
	}

	want = false

	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}
}
