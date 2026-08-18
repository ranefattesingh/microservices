package service

import (
	"testing"
)

// TestEncryptPassword tests password encryption
func TestEncryptPassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == "" {
		t.Fatalf("expected hash, got empty string")
	}

	// Hash should not be equal to plain password
	if hash == password {
		t.Fatalf("expected hashed password to differ from plain text")
	}

	// Hashing same password twice should produce different hashes
	hash2, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if hash == hash2 {
		t.Fatalf("expected different hashes for same password (bcrypt uses salt)")
	}
}

// TestComparePassword tests password comparison
func TestComparePassword(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("failed to encrypt password: %v", err)
	}

	// Correct password should match
	err = ComparePassword(hash, password)
	if err != nil {
		t.Fatalf("expected password match, got error: %v", err)
	}
}

// TestComparePasswordWrong tests password mismatch
func TestComparePasswordWrong(t *testing.T) {
	password := "mySecurePassword123"
	wrongPassword := "differentPassword456"

	hash, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("failed to encrypt password: %v", err)
	}

	// Wrong password should not match
	err = ComparePassword(hash, wrongPassword)
	if err == nil {
		t.Fatalf("expected error for wrong password, got nil")
	}
}

// TestComparePasswordEmpty tests empty password comparison
func TestComparePasswordEmpty(t *testing.T) {
	password := "mySecurePassword123"

	hash, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("failed to encrypt password: %v", err)
	}

	// Empty password should not match
	err = ComparePassword(hash, "")
	if err == nil {
		t.Fatalf("expected error for empty password, got nil")
	}
}

// TestComparePasswordSpecialChars tests passwords with special characters
func TestComparePasswordSpecialChars(t *testing.T) {
	password := "P@ssw0rd!#$%^&*()"

	hash, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("failed to encrypt password: %v", err)
	}

	err = ComparePassword(hash, password)
	if err != nil {
		t.Fatalf("expected password match with special chars, got error: %v", err)
	}
}

// TestComparePasswordUnicode tests passwords with unicode characters
func TestComparePasswordUnicode(t *testing.T) {
	password := "pässwörd中文🔐"

	hash, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("failed to encrypt password: %v", err)
	}

	err = ComparePassword(hash, password)
	if err != nil {
		t.Fatalf("expected password match with unicode, got error: %v", err)
	}
}

// TestEncryptLongPassword tests encryption of password at bcrypt limit (72 bytes)
func TestEncryptLongPassword(t *testing.T) {
	// Bcrypt has a 72-byte limit
	password := ""
	for i := 0; i < 72; i++ {
		password += "a"
	}

	hash, err := EncryptPassword(password)
	if err != nil {
		t.Fatalf("failed to encrypt 72-byte password: %v", err)
	}

	err = ComparePassword(hash, password)
	if err != nil {
		t.Fatalf("expected password match at 72-byte limit, got error: %v", err)
	}
}
