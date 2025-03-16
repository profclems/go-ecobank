package ecobank

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func encodePayload(exp int64) string {
	claims := tokenClaims{Exp: exp}
	jsonData, _ := json.Marshal(claims)
	encoded := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(jsonData)
	return encoded
}

func TestGetTokenExpiry(t *testing.T) {
	t.Run("valid token", func(t *testing.T) {
		expTime := time.Now().Add(time.Hour).Unix()
		token := fmt.Sprintf("header.%s.signature", encodePayload(expTime))

		expiry, err := getTokenExpiry(token)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if expiry.Unix() != expTime {
			t.Errorf("expected expiry %v, got %v", expTime, expiry.Unix())
		}
	})

	t.Run("invalid JWT format", func(t *testing.T) {
		_, err := getTokenExpiry("invalid.token")
		if err == nil || !strings.Contains(err.Error(), "invalid JWT format") {
			t.Errorf("expected invalid JWT format error, got %v", err)
		}
	})

	t.Run("invalid base64 payload", func(t *testing.T) {
		_, err := getTokenExpiry("header.invalidbase64.signature")
		if err == nil {
			t.Error("expected error due to invalid base64 payload")
		}
	})

	t.Run("invalid JSON payload", func(t *testing.T) {
		invalidPayload := base64.URLEncoding.EncodeToString([]byte("invalid json"))
		_, err := getTokenExpiry(fmt.Sprintf("header.%s.signature", invalidPayload))
		if err == nil {
			t.Error("expected error due to invalid JSON payload")
		}
	})
}

func BenchmarkGetTokenExpiry(b *testing.B) {
	expTime := time.Now().Add(time.Hour).Unix()
	token := fmt.Sprintf("header.%s.signature", encodePayload(expTime))

	for i := 0; i < b.N; i++ {
		_, _ = getTokenExpiry(token)
	}
}
