package security

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
	DefaultBcryptCost = 12
)

var (
	ErrPasswordEmpty       = errors.New("makeadmin password is empty")
	ErrPasswordTooShort    = errors.New("makeadmin password is too short")
	ErrPasswordTooLong     = errors.New("makeadmin password is too long")
	ErrPasswordPlaceholder = errors.New("makeadmin password hash is an install-time placeholder")
	ErrPasswordUnsupported = errors.New("makeadmin password hash is unsupported")
)

type PasswordDigest struct {
	Hash string
	Salt string
}

type PasswordHasher interface {
	Hash(plain string) (PasswordDigest, error)
	Verify(plain string, digest PasswordDigest) (bool, error)
	NeedsUpgrade(digest PasswordDigest) bool
}

type BcryptPasswordHasher struct {
	cost int
}

func NewBcryptPasswordHasher(cost int) BcryptPasswordHasher {
	if cost == 0 {
		cost = DefaultBcryptCost
	}
	return BcryptPasswordHasher{cost: cost}
}

func (hasher BcryptPasswordHasher) Hash(plain string) (PasswordDigest, error) {
	if err := ValidatePassword(plain); err != nil {
		return PasswordDigest{}, err
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(plain), hasher.cost)
	if err != nil {
		return PasswordDigest{}, err
	}
	return PasswordDigest{Hash: string(hash)}, nil
}

func (hasher BcryptPasswordHasher) Verify(plain string, digest PasswordDigest) (bool, error) {
	if isInstallTimePlaceholder(digest.Hash) {
		return false, ErrPasswordPlaceholder
	}
	if isBcryptHash(digest.Hash) {
		err := bcrypt.CompareHashAndPassword([]byte(digest.Hash), []byte(plain))
		if err == nil {
			return true, nil
		}
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	if isLegacyMD5Hash(digest.Hash, digest.Salt) {
		sum := md5.Sum([]byte(plain + digest.Salt))
		expected := hex.EncodeToString(sum[:])
		matched := subtle.ConstantTimeCompare([]byte(expected), []byte(digest.Hash)) == 1
		return matched, nil
	}
	return false, ErrPasswordUnsupported
}

func (hasher BcryptPasswordHasher) NeedsUpgrade(digest PasswordDigest) bool {
	if isLegacyMD5Hash(digest.Hash, digest.Salt) {
		return true
	}
	if !isBcryptHash(digest.Hash) {
		return false
	}
	cost, err := bcrypt.Cost([]byte(digest.Hash))
	return err == nil && cost < hasher.cost
}

func ValidatePassword(plain string) error {
	switch length := len(plain); {
	case length == 0:
		return ErrPasswordEmpty
	case length < MinPasswordLength:
		return ErrPasswordTooShort
	case length > MaxPasswordLength:
		return ErrPasswordTooLong
	default:
		return nil
	}
}

func isBcryptHash(hash string) bool {
	return strings.HasPrefix(hash, "$2a$") ||
		strings.HasPrefix(hash, "$2b$") ||
		strings.HasPrefix(hash, "$2x$") ||
		strings.HasPrefix(hash, "$2y$")
}

func isLegacyMD5Hash(hash string, salt string) bool {
	return len(hash) == 32 && salt != "" && isLowerHex(hash)
}

func isInstallTimePlaceholder(hash string) bool {
	return strings.Contains(hash, "INSTALL_TIME") || strings.Contains(hash, "REPLACE_ME")
}

func isLowerHex(value string) bool {
	for _, char := range value {
		if (char < '0' || char > '9') && (char < 'a' || char > 'f') {
			return false
		}
	}
	return true
}
