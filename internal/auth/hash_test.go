package auth

import (
	"testing"
)

func TestHashPasswordAndCheckPasswordHash(t *testing.T) {
	password := "supersecret"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = CheckPasswordHash(password, hash)
	if err != nil {
		t.Errorf("CheckPasswordHash failed for correct password: %v", err)
	}
}

func TestCheckPasswordHash_WrongPassword(t *testing.T) {
	password := "supersecret"
	wrongPassword := "notsecret"
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	err = CheckPasswordHash(wrongPassword, hash)
	if err == nil {
		t.Error("Expected error for wrong password, got nil")
	}
}
