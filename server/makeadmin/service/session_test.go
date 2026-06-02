package service

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestJWTTokenCodecIssueAndParse(t *testing.T) {
	now := time.Unix(1000, 0)
	codec := NewJWTTokenCodec("test-secret")
	codec.now = func() time.Time { return now }

	token, err := codec.Issue(Identity{
		AdminID:  7,
		TenantID: 3,
	}, 3600)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	if token.AccessToken == "" || token.SessionID == "" {
		t.Fatalf("Issue() token = %#v, want access token and session id", token)
	}

	claims, err := codec.Parse(token.AccessToken)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if claims.SessionID != token.SessionID || claims.AdminID != 7 || claims.TenantID != 3 {
		t.Fatalf("claims = %#v, want issued identity", claims)
	}
	if claims.IssuedAt != 1000 || claims.ExpiresAt != 4600 || claims.Issuer != TokenIssuer {
		t.Fatalf("time/issuer claims = %#v, want iat=1000 exp=4600 issuer=%q", claims, TokenIssuer)
	}
}

func TestJWTTokenCodecRejectsTamperedToken(t *testing.T) {
	codec := NewJWTTokenCodec("test-secret")
	token, err := codec.Issue(Identity{AdminID: 1}, 3600)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	tampered := token.AccessToken
	if strings.HasSuffix(tampered, "a") {
		tampered = strings.TrimSuffix(tampered, "a") + "b"
	} else {
		tampered += "a"
	}

	_, err = codec.Parse(tampered)
	if !errors.Is(err, ErrTokenInvalid) {
		t.Fatalf("Parse() tampered error = %v, want ErrTokenInvalid", err)
	}
}

func TestJWTTokenCodecRejectsExpiredToken(t *testing.T) {
	now := time.Unix(1000, 0)
	codec := NewJWTTokenCodec("test-secret")
	codec.now = func() time.Time { return now }
	token, err := codec.Issue(Identity{AdminID: 1}, 10)
	if err != nil {
		t.Fatalf("Issue() error = %v", err)
	}
	codec.now = func() time.Time { return now.Add(11 * time.Second) }

	_, err = codec.Parse(token.AccessToken)
	if !errors.Is(err, ErrTokenExpired) {
		t.Fatalf("Parse() expired error = %v, want ErrTokenExpired", err)
	}
}
