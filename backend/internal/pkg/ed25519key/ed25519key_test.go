package ed25519key

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPrivateKey(t *testing.T) {
	publicKey, privateKey, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	path := writePrivateKeyPEM(t, privateKey)

	got, err := LoadPrivateKey(path)
	if err != nil {
		t.Fatalf("load private key: %v", err)
	}

	if !bytes.Equal(got, privateKey) {
		t.Fatalf("private key mismatch")
	}

	if !bytes.Equal(got.Public().(ed25519.PublicKey), publicKey) {
		t.Fatalf("derived public key mismatch")
	}
}

func TestLoadPublicKey(t *testing.T) {
	publicKey, _, err := ed25519.GenerateKey(nil)
	if err != nil {
		t.Fatalf("generate key: %v", err)
	}

	path := writePublicKeyPEM(t, publicKey)

	got, err := LoadPublicKey(path)
	if err != nil {
		t.Fatalf("load public key: %v", err)
	}

	if !bytes.Equal(got, publicKey) {
		t.Fatalf("public key mismatch")
	}
}

func TestLoadPrivateKey_FileMissing(t *testing.T) {
	_, err := LoadPrivateKey(filepath.Join(t.TempDir(), "missing.pem"))
	if err == nil {
		t.Fatalf("expected error for missing file")
	}
}

func TestLoadPrivateKey_InvalidPEM(t *testing.T) {
	path := filepath.Join(t.TempDir(), "garbage.pem")
	if err := os.WriteFile(path, []byte("not a pem"), 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	_, err := LoadPrivateKey(path)
	if err == nil {
		t.Fatalf("expected error for invalid pem")
	}
}

func writePrivateKeyPEM(t *testing.T, key ed25519.PrivateKey) string {
	t.Helper()

	der, err := x509.MarshalPKCS8PrivateKey(key)
	if err != nil {
		t.Fatalf("marshal pkcs8: %v", err)
	}

	path := filepath.Join(t.TempDir(), "private.pem")
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der})
	if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	return path
}

func writePublicKeyPEM(t *testing.T, key ed25519.PublicKey) string {
	t.Helper()

	der, err := x509.MarshalPKIXPublicKey(key)
	if err != nil {
		t.Fatalf("marshal pkix: %v", err)
	}

	path := filepath.Join(t.TempDir(), "public.pem")
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: der})
	if err := os.WriteFile(path, pemBytes, 0o600); err != nil {
		t.Fatalf("write file: %v", err)
	}

	return path
}
