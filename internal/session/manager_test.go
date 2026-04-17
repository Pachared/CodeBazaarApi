package session

import (
	"strings"
	"testing"
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/models"
)

func TestManagerSignAndParse(t *testing.T) {
	manager, err := NewManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	token, expiresAt, err := manager.Sign(&models.User{
		ID:       "usr_123",
		Name:     "Pachara",
		Email:    "PACHARA@example.com",
		Role:     "buyer",
		Provider: "google",
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	if token == "" {
		t.Fatal("Sign() returned an empty token")
	}

	if !expiresAt.After(time.Now()) {
		t.Fatal("Sign() returned an expired token")
	}

	claims, err := manager.Parse(token)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	if claims.UserID != "usr_123" {
		t.Fatalf("claims.UserID = %q, want %q", claims.UserID, "usr_123")
	}

	if claims.Email != "pachara@example.com" {
		t.Fatalf("claims.Email = %q, want %q", claims.Email, "pachara@example.com")
	}
}

func TestManagerRejectsTamperedToken(t *testing.T) {
	manager, err := NewManager("test-secret", time.Hour)
	if err != nil {
		t.Fatalf("NewManager() error = %v", err)
	}

	token, _, err := manager.Sign(&models.User{
		ID:       "usr_456",
		Name:     "Buyer",
		Email:    "buyer@example.com",
		Role:     "buyer",
		Provider: "google",
	})
	if err != nil {
		t.Fatalf("Sign() error = %v", err)
	}

	tamperedToken := token + "tampered"
	_, err = manager.Parse(tamperedToken)
	if err == nil {
		t.Fatal("Parse() unexpectedly accepted a tampered token")
	}

	if !strings.Contains(err.Error(), "invalid") {
		t.Fatalf("Parse() error = %v, want invalid token error", err)
	}
}
