package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultSessionTTLSeconds = 7200

	SessionStateKeyPrefix = "makeadmin:session:"
	SessionStateSetPrefix = "makeadmin:session:set:"
	TokenIssuer           = "go-makeadmin"
)

var (
	ErrSessionStore   = errors.New("makeadmin session store error")
	ErrTokenGenerator = errors.New("makeadmin token generator error")
	ErrTokenInvalid   = errors.New("makeadmin token invalid")
	ErrTokenExpired   = errors.New("makeadmin token expired")
)

type SessionToken struct {
	AccessToken string
	SessionID   string
}

type TokenClaims struct {
	SessionID string
	AdminID   uint64
	TenantID  uint64
	IssuedAt  int64
	ExpiresAt int64
	Issuer    string
}

type TokenCodec interface {
	Issue(identity Identity, ttlSeconds int) (SessionToken, error)
	Parse(token string) (TokenClaims, error)
}

type JWTTokenCodec struct {
	secret []byte
	now    func() time.Time
}

type jwtHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

type jwtPayload struct {
	SessionID string `json:"sid"`
	AdminID   uint64 `json:"adminId"`
	TenantID  uint64 `json:"tenantId"`
	IssuedAt  int64  `json:"iat"`
	ExpiresAt int64  `json:"exp"`
	Issuer    string `json:"iss"`
}

func NewJWTTokenCodec(secret string) JWTTokenCodec {
	if secret == "" {
		secret = "go-makeadmin-dev-secret"
	}
	return JWTTokenCodec{secret: []byte(secret), now: time.Now}
}

func (codec JWTTokenCodec) Issue(identity Identity, ttlSeconds int) (SessionToken, error) {
	if ttlSeconds <= 0 {
		ttlSeconds = DefaultSessionTTLSeconds
	}
	sessionID, err := randomSessionID()
	if err != nil {
		return SessionToken{}, err
	}
	now := codec.currentTime().Unix()
	payload := jwtPayload{
		SessionID: sessionID,
		AdminID:   identity.AdminID,
		TenantID:  identity.TenantID,
		IssuedAt:  now,
		ExpiresAt: now + int64(ttlSeconds),
		Issuer:    TokenIssuer,
	}
	token, err := codec.sign(payload)
	if err != nil {
		return SessionToken{}, err
	}
	return SessionToken{AccessToken: token, SessionID: sessionID}, nil
}

func (codec JWTTokenCodec) Parse(token string) (TokenClaims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return TokenClaims{}, ErrTokenInvalid
	}
	signingInput := parts[0] + "." + parts[1]
	expectedSignature := codec.signature(signingInput)
	if !hmac.Equal([]byte(parts[2]), []byte(expectedSignature)) {
		return TokenClaims{}, ErrTokenInvalid
	}

	var header jwtHeader
	if err := decodeJWTPart(parts[0], &header); err != nil {
		return TokenClaims{}, errors.Join(ErrTokenInvalid, err)
	}
	if header.Alg != "HS256" || header.Typ != "JWT" {
		return TokenClaims{}, ErrTokenInvalid
	}

	var payload jwtPayload
	if err := decodeJWTPart(parts[1], &payload); err != nil {
		return TokenClaims{}, errors.Join(ErrTokenInvalid, err)
	}
	if payload.SessionID == "" || payload.AdminID == 0 || payload.Issuer != TokenIssuer {
		return TokenClaims{}, ErrTokenInvalid
	}
	if payload.ExpiresAt <= codec.currentTime().Unix() {
		return TokenClaims{}, ErrTokenExpired
	}
	return TokenClaims{
		SessionID: payload.SessionID,
		AdminID:   payload.AdminID,
		TenantID:  payload.TenantID,
		IssuedAt:  payload.IssuedAt,
		ExpiresAt: payload.ExpiresAt,
		Issuer:    payload.Issuer,
	}, nil
}

func (codec JWTTokenCodec) sign(payload jwtPayload) (string, error) {
	headerBytes, err := json.Marshal(jwtHeader{Alg: "HS256", Typ: "JWT"})
	if err != nil {
		return "", errors.Join(ErrTokenGenerator, err)
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", errors.Join(ErrTokenGenerator, err)
	}
	signingInput := encodeJWTPart(headerBytes) + "." + encodeJWTPart(payloadBytes)
	return signingInput + "." + codec.signature(signingInput), nil
}

func (codec JWTTokenCodec) signature(signingInput string) string {
	mac := hmac.New(sha256.New, codec.secret)
	_, _ = mac.Write([]byte(signingInput))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}

func (codec JWTTokenCodec) currentTime() time.Time {
	if codec.now == nil {
		return time.Now()
	}
	return codec.now()
}

func randomSessionID() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", errors.Join(ErrTokenGenerator, err)
	}
	return base64.RawURLEncoding.EncodeToString(token), nil
}

func encodeJWTPart(data []byte) string {
	return base64.RawURLEncoding.EncodeToString(data)
}

func decodeJWTPart(part string, target interface{}) error {
	data, err := base64.RawURLEncoding.DecodeString(part)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

type SessionStore interface {
	Save(ctx context.Context, sessionID string, identity Identity, ttlSeconds int) error
	Delete(ctx context.Context, sessionID string) error
	FindAdminID(ctx context.Context, sessionID string) (uint64, error)
	Refresh(ctx context.Context, sessionID string, ttlSeconds int) error
}

type UnavailableSessionStore struct{}

type redisSessionStore struct {
	client *redis.Client
	prefix string
}

func NewRedisSessionStore(client *redis.Client, prefix string) SessionStore {
	return redisSessionStore{client: client, prefix: prefix}
}

func (UnavailableSessionStore) Save(ctx context.Context, sessionID string, identity Identity, ttlSeconds int) error {
	return ErrSessionStore
}

func (UnavailableSessionStore) Delete(ctx context.Context, sessionID string) error {
	return ErrSessionStore
}

func (UnavailableSessionStore) FindAdminID(ctx context.Context, sessionID string) (uint64, error) {
	return 0, ErrSessionStore
}

func (UnavailableSessionStore) Refresh(ctx context.Context, sessionID string, ttlSeconds int) error {
	return ErrSessionStore
}

func (store redisSessionStore) Save(ctx context.Context, sessionID string, identity Identity, ttlSeconds int) error {
	if store.client == nil {
		return ErrSessionStore
	}
	if sessionID == "" || identity.AdminID == 0 {
		return ErrSessionStore
	}
	adminID := strconv.FormatUint(identity.AdminID, 10)
	ttl := time.Duration(ttlSeconds) * time.Second
	if err := store.client.Set(ctx, store.prefix+SessionStateKeyPrefix+sessionID, adminID, ttl).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	sessionSetKey := store.prefix + SessionStateSetPrefix + adminID
	if err := store.client.SAdd(ctx, sessionSetKey, sessionID).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	if err := store.client.Expire(ctx, sessionSetKey, ttl).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	return nil
}

func (store redisSessionStore) Delete(ctx context.Context, sessionID string) error {
	if store.client == nil {
		return ErrSessionStore
	}
	if sessionID == "" {
		return nil
	}
	adminID, _ := store.FindAdminID(ctx, sessionID)
	if err := store.client.Del(ctx, store.prefix+SessionStateKeyPrefix+sessionID).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	if adminID > 0 {
		if err := store.client.SRem(ctx, store.prefix+SessionStateSetPrefix+strconv.FormatUint(adminID, 10), sessionID).Err(); err != nil {
			return errors.Join(ErrSessionStore, err)
		}
	}
	return nil
}

func (store redisSessionStore) FindAdminID(ctx context.Context, sessionID string) (uint64, error) {
	if store.client == nil {
		return 0, ErrSessionStore
	}
	value, err := store.client.Get(ctx, store.prefix+SessionStateKeyPrefix+sessionID).Result()
	if errors.Is(err, redis.Nil) {
		return 0, ErrTokenInvalid
	}
	if err != nil {
		return 0, errors.Join(ErrSessionStore, err)
	}
	adminID, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, errors.Join(ErrSessionStore, fmt.Errorf("invalid admin id %q: %w", value, err))
	}
	return adminID, nil
}

func (store redisSessionStore) Refresh(ctx context.Context, sessionID string, ttlSeconds int) error {
	if store.client == nil {
		return ErrSessionStore
	}
	if sessionID == "" || ttlSeconds <= 0 {
		return nil
	}
	adminID, err := store.FindAdminID(ctx, sessionID)
	if err != nil {
		return err
	}
	ttl := time.Duration(ttlSeconds) * time.Second
	if err := store.client.Expire(ctx, store.prefix+SessionStateKeyPrefix+sessionID, ttl).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	if err := store.client.Expire(ctx, store.prefix+SessionStateSetPrefix+strconv.FormatUint(adminID, 10), ttl).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	return nil
}
