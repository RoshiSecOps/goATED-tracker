package auth

import (
	"net/http"
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

func TestGetBearerFunc(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Fatalf("dummy request creation failed: %v", err)
	}
	req.Header.Add("Authorization", "Bearer xyz")
	token, err := GetBearerToken(req.Header)
	if err != nil {
		t.Fatalf("unable to get header")
	}
	if token != "xyz" {
		t.Errorf("expected token to be xyz, got %v", token)
	}
}

func TestWrongHeader(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Fatalf("dummy request creation failed: %v", err)
	}
	req.Header.Add("AuthorizationZ", "Bearer xyz")
	token, err := GetBearerToken(req.Header)
	if err == nil {
		t.Fatalf("expected an error as header is not present, got token: %v", token)
	}
}

func TestWrongHeaderContent(t *testing.T) {
	req, err := http.NewRequest("GET", "https://google.com", nil)
	if err != nil {
		t.Fatalf("dummy request creation failed: %v", err)
	}
	req.Header.Add("Authorization", "Bearer")
	token, err := GetBearerToken(req.Header)
	if err == nil {
		t.Fatalf("expected error of malformed auth header, got token: %v", token)
	}
}
