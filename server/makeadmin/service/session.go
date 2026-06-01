package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	DefaultSessionTTLSeconds = 7200

	SessionTokenKeyPrefix = "makeadmin:token:"
	SessionTokenSetPrefix = "makeadmin:token:set:"
)

var (
	ErrSessionStore   = errors.New("makeadmin session store error")
	ErrTokenGenerator = errors.New("makeadmin token generator error")
)

type TokenGenerator interface {
	Generate() (string, error)
}

type RandomTokenGenerator struct{}

func (RandomTokenGenerator) Generate() (string, error) {
	token := make([]byte, 32)
	if _, err := rand.Read(token); err != nil {
		return "", errors.Join(ErrTokenGenerator, err)
	}
	return base64.RawURLEncoding.EncodeToString(token), nil
}

type SessionStore interface {
	Save(ctx context.Context, token string, identity Identity, ttlSeconds int) error
	Delete(ctx context.Context, token string) error
}

type UnavailableSessionStore struct{}

type redisSessionStore struct {
	client *redis.Client
	prefix string
}

func NewRedisSessionStore(client *redis.Client, prefix string) SessionStore {
	return redisSessionStore{client: client, prefix: prefix}
}

func (UnavailableSessionStore) Save(ctx context.Context, token string, identity Identity, ttlSeconds int) error {
	return ErrSessionStore
}

func (UnavailableSessionStore) Delete(ctx context.Context, token string) error {
	return ErrSessionStore
}

func (store redisSessionStore) Save(ctx context.Context, token string, identity Identity, ttlSeconds int) error {
	if store.client == nil {
		return ErrSessionStore
	}
	adminID := strconv.FormatUint(identity.AdminID, 10)
	ttl := time.Duration(ttlSeconds) * time.Second
	if err := store.client.Set(ctx, store.prefix+SessionTokenKeyPrefix+token, adminID, ttl).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	if err := store.client.SAdd(ctx, store.prefix+SessionTokenSetPrefix+adminID, token).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	return nil
}

func (store redisSessionStore) Delete(ctx context.Context, token string) error {
	if store.client == nil {
		return ErrSessionStore
	}
	if err := store.client.Del(ctx, store.prefix+SessionTokenKeyPrefix+token).Err(); err != nil {
		return errors.Join(ErrSessionStore, err)
	}
	return nil
}
