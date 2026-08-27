package security

import "testing"

func TestVaultEncryptsAndDecrypts(t *testing.T) {
	vault, err := NewVault("test-vault-seed")
	if err != nil {
		t.Fatal(err)
	}
	credential, err := vault.Encrypt("office", "100001", "secret-value")
	if err != nil {
		t.Fatal(err)
	}
	if credential.Ciphertext == "secret-value" || !vault.Verify("100001", credential, "secret-value") {
		t.Fatalf("credential was not protected: %#v", credential)
	}
	decrypted, err := vault.Decrypt("100001", credential)
	if err != nil || decrypted != "secret-value" {
		t.Fatalf("decrypted value: %q %v", decrypted, err)
	}
}

func TestVaultRejectsWrongAccount(t *testing.T) {
	vault, _ := NewVault("test-vault-seed")
	credential, _ := vault.Encrypt("office", "100001", "secret-value")
	if vault.Verify("100002", credential, "secret-value") {
		t.Fatal("wrong account must not verify")
	}
}
