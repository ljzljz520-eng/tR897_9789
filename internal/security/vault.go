package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
)

type Vault struct {
	key []byte
}

type EncryptedCredential struct {
	Domain     string
	Ciphertext string
}

func NewVault(seed string) (*Vault, error) {
	if strings.TrimSpace(seed) == "" {
		return nil, fmt.Errorf("vault seed is required")
	}
	digest := sha256.Sum256([]byte(seed))
	return &Vault{key: append([]byte(nil), digest[:]...)}, nil
}

func (v *Vault) block() (cipher.Block, error) {
	if v == nil || len(v.key) != 32 {
		return nil, fmt.Errorf("vault key is invalid")
	}
	return aes.NewCipher(v.key)
}

func (v *Vault) nonce(domain, account string) []byte {
	digest := sha256.Sum256(append(append([]byte(nil), v.key...), []byte(domain+":"+account)...))
	return append([]byte(nil), digest[:12]...)
}

func (v *Vault) Encrypt(domain, account, plaintext string) (EncryptedCredential, error) {
	if strings.TrimSpace(domain) == "" || strings.TrimSpace(account) == "" || plaintext == "" {
		return EncryptedCredential{}, fmt.Errorf("credential fields are required")
	}
	block, err := v.block()
	if err != nil {
		return EncryptedCredential{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return EncryptedCredential{}, err
	}
	ciphertext := gcm.Seal(nil, v.nonce(domain, account), []byte(plaintext), []byte(domain))
	return EncryptedCredential{Domain: domain, Ciphertext: base64.RawStdEncoding.EncodeToString(ciphertext)}, nil
}

func (v *Vault) Decrypt(account string, credential EncryptedCredential) (string, error) {
	if account == "" || credential.Domain == "" || credential.Ciphertext == "" {
		return "", fmt.Errorf("encrypted credential fields are required")
	}
	block, err := v.block()
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.RawStdEncoding.DecodeString(credential.Ciphertext)
	if err != nil {
		return "", fmt.Errorf("decode credential: %w", err)
	}
	plaintext, err := gcm.Open(nil, v.nonce(credential.Domain, account), ciphertext, []byte(credential.Domain))
	if err != nil {
		return "", fmt.Errorf("decrypt credential: %w", err)
	}
	return string(plaintext), nil
}

func (v *Vault) Verify(account string, credential EncryptedCredential, expected string) bool {
	value, err := v.Decrypt(account, credential)
	return err == nil && value == expected
}

func RedactedCredential(credential EncryptedCredential) string {
	if credential.Ciphertext == "" {
		return "<empty>"
	}
	return credential.Domain + ":<encrypted>"
}
