package security

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestBcryptPasswordHasherHashAndVerify(t *testing.T) {
	hasher := NewBcryptPasswordHasher(bcrypt.MinCost)

	digest, err := hasher.Hash("makeadmin-secret")
	if err != nil {
		t.Fatalf("Hash() error = %v", err)
	}
	if digest.Hash == "" {
		t.Fatal("Hash() returned empty hash")
	}
	if digest.Salt != "" {
		t.Fatalf("Hash() salt = %q, want empty bcrypt salt field", digest.Salt)
	}

	matched, err := hasher.Verify("makeadmin-secret", digest)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if !matched {
		t.Fatal("Verify() matched = false, want true")
	}

	matched, err = hasher.Verify("wrong-secret", digest)
	if err != nil {
		t.Fatalf("Verify() wrong password error = %v", err)
	}
	if matched {
		t.Fatal("Verify() wrong password matched = true, want false")
	}
}

func TestBcryptPasswordHasherVerifyLegacyMD5Salt(t *testing.T) {
	hasher := NewBcryptPasswordHasher(bcrypt.MinCost)
	sum := md5.Sum([]byte("makeadmin-secret" + "abcde"))
	digest := PasswordDigest{
		Hash: hex.EncodeToString(sum[:]),
		Salt: "abcde",
	}

	matched, err := hasher.Verify("makeadmin-secret", digest)
	if err != nil {
		t.Fatalf("Verify() legacy error = %v", err)
	}
	if !matched {
		t.Fatal("Verify() legacy matched = false, want true")
	}
	if !hasher.NeedsUpgrade(digest) {
		t.Fatal("NeedsUpgrade() legacy = false, want true")
	}
}

func TestBcryptPasswordHasherRejectsInstallPlaceholder(t *testing.T) {
	hasher := NewBcryptPasswordHasher(bcrypt.MinCost)

	_, err := hasher.Verify("makeadmin-secret", PasswordDigest{
		Hash: "INSTALL_TIME_PASSWORD_BCRYPT_REPLACE_ME",
	})
	if !errors.Is(err, ErrPasswordPlaceholder) {
		t.Fatalf("Verify() error = %v, want ErrPasswordPlaceholder", err)
	}
}

func TestValidatePasswordLength(t *testing.T) {
	if !errors.Is(ValidatePassword("short"), ErrPasswordTooShort) {
		t.Fatal("ValidatePassword() short password did not return ErrPasswordTooShort")
	}
	if err := ValidatePassword("12345678"); err != nil {
		t.Fatalf("ValidatePassword() valid password error = %v", err)
	}
}
