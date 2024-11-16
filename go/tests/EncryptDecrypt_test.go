package tests

import (
	"github.com/saichler/shared/go/share/aes"
	. "github.com/saichler/shared/go/share/interfaces"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	base := "Test Data To Encrypt"
	key := aes.GenerateAES256Key()
	encData, err := aes.Encrypt([]byte(base), key)
	if err != nil {
		Fail("Failed to encrypt data:", err)
		return
	}
	decData, err := aes.Decrypt(encData, key)
	if err != nil {
		Fail("Failed to decrypt data:", err)
		return
	}
	out := string(decData)
	if out != base {
		Fail("Decrypted data is not equal to base:", out)
		return
	}
}
