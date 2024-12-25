package encryption

import (
	"testing"
)

func TestEncrypt(t *testing.T) {
	t.Run("Should encrypt valid payload", func(t *testing.T) {
		encryption := NewAesCfbEncryption()
		payload := &AESPayload{
			Secret:  "00112233445566778899aabbccddeeff", // valid 16-byte key
			Payload: "Hello, World!",
		}

		encrypted, err := encryption.Encrypt(payload)
		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if encrypted == "" {
			t.Errorf("Expected encrypted string, got empty string")
		}
	})
}

func TestDecrypt(t *testing.T) {
	t.Run("Should decrypt valid payload", func(t *testing.T) {
		encryption := NewAesCfbEncryption()
		payload := &AESPayload{
			Secret:  "00112233445566778899aabbccddeeff", // valid 16-byte key
			Payload: "",                                 // We will fill this in with the encrypted version of "Hello, World!"
		}

		// First, let's encrypt the string "Hello, World!" so we know what to expect when we decrypt it.
		encryptPayload := &AESPayload{
			Secret:  "00112233445566778899aabbccddeeff", // valid 16-byte key
			Payload: "Hello, World!",
		}
		encrypted, err := encryption.Encrypt(encryptPayload)
		if err != nil {
			t.Errorf("Unexpected error during encryption: %v", err)
		}

		// Now we can fill in the Payload field with the encrypted string.
		payload.Payload = encrypted

		decrypted, err := encryption.Decrypt(payload)
		if err != nil {
			t.Errorf("Unexpected error during decryption: %v", err)
		}

		if decrypted != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got '%v'", decrypted)
		}
	})
}
