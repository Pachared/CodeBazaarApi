package session

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/Pachared/CodeBazaarApi/internal/models"
)

var (
	ErrInvalidToken = errors.New("invalid session token")
	ErrExpiredToken = errors.New("expired session token")
)

type Claims struct {
	UserID    string `json:"userId"`
	Email     string `json:"email"`
	Name      string `json:"name"`
	Role      string `json:"role"`
	Provider  string `json:"provider"`
	ExpiresAt int64  `json:"expiresAt"`
	IssuedAt  int64  `json:"issuedAt"`
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) (*Manager, error) {
	normalizedSecret := strings.TrimSpace(secret)
	if normalizedSecret == "" {
		return nil, errors.New("session secret is required")
	}
	if ttl <= 0 {
		return nil, errors.New("session ttl must be greater than zero")
	}

	return &Manager{
		secret: []byte(normalizedSecret),
		ttl:    ttl,
	}, nil
}

func (m *Manager) Sign(user *models.User) (string, time.Time, error) {
	if user == nil {
		return "", time.Time{}, errors.New("user is required")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(m.ttl)

	claims := Claims{
		UserID:    strings.TrimSpace(user.ID),
		Email:     strings.TrimSpace(strings.ToLower(user.Email)),
		Name:      strings.TrimSpace(user.Name),
		Role:      strings.TrimSpace(user.Role),
		Provider:  strings.TrimSpace(user.Provider),
		ExpiresAt: expiresAt.Unix(),
		IssuedAt:  now.Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}

	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)
	signature := m.sign([]byte(encodedPayload))
	encodedSignature := base64.RawURLEncoding.EncodeToString(signature)

	return encodedPayload + "." + encodedSignature, expiresAt, nil
}

func (m *Manager) Parse(token string) (*Claims, error) {
	parts := strings.Split(strings.TrimSpace(token), ".")
	if len(parts) != 2 {
		return nil, ErrInvalidToken
	}

	payloadPart := strings.TrimSpace(parts[0])
	signaturePart := strings.TrimSpace(parts[1])
	if payloadPart == "" || signaturePart == "" {
		return nil, ErrInvalidToken
	}

	expectedSignature := m.sign([]byte(payloadPart))
	actualSignature, err := base64.RawURLEncoding.DecodeString(signaturePart)
	if err != nil {
		return nil, ErrInvalidToken
	}

	if subtle.ConstantTimeCompare(actualSignature, expectedSignature) != 1 {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(payloadPart)
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if strings.TrimSpace(claims.UserID) == "" || strings.TrimSpace(claims.Email) == "" {
		return nil, ErrInvalidToken
	}

	if time.Now().UTC().Unix() >= claims.ExpiresAt {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func (m *Manager) sign(payload []byte) []byte {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write(payload)
	return mac.Sum(nil)
}
