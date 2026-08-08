package jwt

import (
	"testing"
	"time"
)

func newTestJWT() *JWT {
	return NewJWT("test-secret-key", 1)
}

func TestNewJWT(t *testing.T) {
	j := NewJWT("secret", 24)
	if j == nil {
		t.Fatal("NewJWT returned nil")
	}
	if len(j.secret) == 0 {
		t.Error("secret not stored")
	}
	if j.expire != 24*time.Hour {
		t.Errorf("expire = %v, want %v", j.expire, 24*time.Hour)
	}
}

func TestGenTokenAndParseToken(t *testing.T) {
	j := newTestJWT()
	userID := int64(1001)
	username := "alice"

	token, err := j.GenToken(userID, username)
	if err != nil {
		t.Fatalf("GenToken error: %v", err)
	}
	if token == "" {
		t.Fatal("GenToken returned empty token")
	}

	claims, err := j.ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken error: %v", err)
	}
	if claims.UserID != userID {
		t.Errorf("UserID = %d, want %d", claims.UserID, userID)
	}
	if claims.Username != username {
		t.Errorf("Username = %q, want %q", claims.Username, username)
	}
	if claims.Issuer != "li" {
		t.Errorf("Issuer = %q, want %q", claims.Issuer, "li")
	}
	if claims.ExpiresAt == nil || claims.ExpiresAt.Time.Before(time.Now()) {
		t.Error("ExpiresAt should be in the future")
	}
}

func TestParseToken_InvalidToken(t *testing.T) {
	j := newTestJWT()
	if _, err := j.ParseToken("not-a-valid-token"); err == nil {
		t.Error("expected error for invalid token, got nil")
	}
}

func TestParseToken_WrongSecret(t *testing.T) {
	j1 := NewJWT("secret-A", 1)
	token, err := j1.GenToken(1, "bob")
	if err != nil {
		t.Fatalf("GenToken error: %v", err)
	}
	j2 := NewJWT("secret-B", 1)
	if _, err := j2.ParseToken(token); err == nil {
		t.Error("expected error when parsing with wrong secret, got nil")
	}
}

func TestParseToken_AlgConfusion(t *testing.T) {
	// 使用 RS256 签名的 token 不应被 HS256 密钥解析通过（算法混淆防护）。
	// 这里用一个明显非 HS256 结构的伪造 token 验证拒绝逻辑。
	j := newTestJWT()
	fake := "eyJhbGciOiJSUzI1NiJ9.eyJ1c2VyX2lkIjoxfQ.signature"
	if _, err := j.ParseToken(fake); err == nil {
		t.Error("expected error for algorithm confusion attempt, got nil")
	}
}
