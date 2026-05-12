package auth

import (
	"testing"
)

func TestMakeAndValidatePassHash(t *testing.T) {
	passHash, err := HashPassword("test123")
	if err != nil {
		t.Fatalf("password hashing failed: %v", err)
	}
	match, err := CheckPasswordHash("test123", passHash)
	if err != nil {
		t.Fatalf("check hash failed: %v", err)
	}
	if match != true {
		t.Errorf("expected true, got %v", match)
	}
}

func TestWrongHash(t *testing.T) {
	passHash, err := HashPassword("somethingElse")
	if err != nil {
		t.Fatalf("password hashing failed: %v", err)
	}
	match, err := CheckPasswordHash("test123", passHash)
	if err != nil {
		t.Fatalf("check hash failed: %v", err)
	}
	if match != false {
		t.Errorf("expected false, got %v", match)
	}
}
