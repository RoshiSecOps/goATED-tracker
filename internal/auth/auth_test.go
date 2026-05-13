package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
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

func TestJwtCreateAndValidate(t *testing.T) {
	userId := uuid.New()
	secret := "basic"
	expireIn := 60 * time.Second
	token, err := MakeJWT(userId, secret, expireIn)
	if err != nil {
		t.Fatalf("unable to create token: %v", err)
	}
	checkUserId, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("unable to validate token: %v", err)
	}
	if userId != checkUserId {
		t.Errorf("expected decoded ID to equal %v, but go %v", userId, checkUserId)
	}
}

func TestJwtExpiration(t *testing.T) {
	userId := uuid.New()
	secret := "basic"
	expireIn := -time.Second
	token, err := MakeJWT(userId, secret, expireIn)
	if err != nil {
		t.Fatalf("unable to create token: %v", err)
	}
	checkUserId, err := ValidateJWT(token, secret)
	if err == nil {
		t.Errorf("expected expired token error, got nil")
	}
	if checkUserId != uuid.Nil {
		t.Errorf("Expected id to be nil, got %v", checkUserId)
	}
}

func TestJwtWrongSecret(t *testing.T) {
	userId := uuid.New()
	secret := "basic"
	expireIn := -time.Second
	token, err := MakeJWT(userId, secret, expireIn)
	if err != nil {
		t.Fatalf("unable to create token: %v", err)
	}
	checkUserId, err := ValidateJWT(token, "wrong")
	if err == nil {
		t.Errorf("expected error, got nil")
	}
	if checkUserId != uuid.Nil {
		t.Errorf("expected uuid.Nil, got %v", checkUserId)
	}
}
