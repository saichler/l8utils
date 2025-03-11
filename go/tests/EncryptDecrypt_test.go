package tests

import (
	. "github.com/saichler/shared/go/tests/infra"
	"github.com/saichler/types/go/aes"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	base := "Test Data To Encrypt"
	key := aes.GenerateAES256Key()
	encData, err := aes.Encrypt([]byte(base), key)
	if err != nil {
		Log.Fail("Failed to encrypt data:", err)
		return
	}
	decData, err := aes.Decrypt(encData, key)
	if err != nil {
		Log.Fail("Failed to decrypt data:", err)
		return
	}
	out := string(decData)
	if out != base {
		Log.Fail("Decrypted data is not equal to base:", out)
		return
	}
}
