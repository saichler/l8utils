package tests

import (
	"os"
	"testing"

	"github.com/saichler/l8utils/go/utils/certs"
)

func TestCreateCA(t *testing.T) {
	certName := "/tmp/test_ca"

	// Clean up any existing files
	os.Remove(certName + ".ca")
	os.Remove(certName + ".caKey")

	// Create CA
	ca, caKey, err := certs.CreateCA(certName, "TestOrg", "US", "CA",
		"TestCity", "TestStreet", "12345", "test@example.com", 1)

	if err != nil {
		t.Fatalf("Failed to create CA: %v", err)
	}

	if ca == nil {
		t.Error("CA certificate should not be nil")
	}

	if caKey == nil {
		t.Error("CA private key should not be nil")
	}

	// Verify files were created
	if _, err := os.Stat(certName + ".ca"); os.IsNotExist(err) {
		t.Error("CA file was not created")
	}

	if _, err := os.Stat(certName + ".caKey"); os.IsNotExist(err) {
		t.Error("CA key file was not created")
	}

	// Try to create again - should fail
	_, _, err = certs.CreateCA(certName, "TestOrg", "US", "CA",
		"TestCity", "TestStreet", "12345", "test@example.com", 1)

	if err == nil {
		t.Error("Creating CA again should fail")
	}

	// Clean up
	os.Remove(certName + ".ca")
	os.Remove(certName + ".caKey")
}

func TestCreateCrt(t *testing.T) {
	certName := "/tmp/test_crt"
	caName := "/tmp/test_ca_for_crt"

	// Clean up any existing files
	os.Remove(caName + ".ca")
	os.Remove(caName + ".caKey")
	os.Remove(certName + ".crt")
	os.Remove(certName + ".crtKey")

	// First create a CA
	ca, caKey, err := certs.CreateCA(caName, "TestOrg", "US", "CA",
		"TestCity", "TestStreet", "12345", "test@example.com", 1)

	if err != nil {
		t.Fatalf("Failed to create CA: %v", err)
	}

	// Create certificate
	err = certs.CreateCrt(certName, "TestOrg", "US", "CA",
		"TestCity", "TestStreet", "12345", "test@example.com",
		"127.0.0.1", "TestSecret", 8443, 1, ca, caKey)

	if err != nil {
		t.Fatalf("Failed to create certificate: %v", err)
	}

	// Verify files were created
	if _, err := os.Stat(certName + ".crt"); os.IsNotExist(err) {
		t.Error("Certificate file was not created")
	}

	if _, err := os.Stat(certName + ".crtKey"); os.IsNotExist(err) {
		t.Error("Certificate key file was not created")
	}

	// Try to create again - should fail
	err = certs.CreateCrt(certName, "TestOrg", "US", "CA",
		"TestCity", "TestStreet", "12345", "test@example.com",
		"127.0.0.1", "TestSecret", 8443, 1, ca, caKey)

	if err == nil {
		t.Error("Creating certificate again should fail")
	}

	// Clean up
	os.Remove(caName + ".ca")
	os.Remove(caName + ".caKey")
	os.Remove(certName + ".crt")
	os.Remove(certName + ".crtKey")
}

func TestCreateLayer8CA(t *testing.T) {
	certName := "/tmp/test_layer8_ca"

	// Clean up any existing files
	os.Remove(certName + ".ca")
	os.Remove(certName + ".caKey")

	// Create Layer8 CA
	ca, caKey, err := certs.CreateLayer8CA(certName)

	if err != nil {
		t.Fatalf("Failed to create Layer8 CA: %v", err)
	}

	if ca == nil {
		t.Error("CA certificate should not be nil")
	}

	if caKey == nil {
		t.Error("CA private key should not be nil")
	}

	// Verify files were created
	if _, err := os.Stat(certName + ".ca"); os.IsNotExist(err) {
		t.Error("CA file was not created")
	}

	if _, err := os.Stat(certName + ".caKey"); os.IsNotExist(err) {
		t.Error("CA key file was not created")
	}

	// Clean up
	os.Remove(certName + ".ca")
	os.Remove(certName + ".caKey")
}

func TestCreateLayer8Crt(t *testing.T) {
	certName := "/tmp/test_layer8_crt"

	// Clean up any existing files
	os.Remove(certName + ".ca")
	os.Remove(certName + ".caKey")
	os.Remove(certName + ".crt")
	os.Remove(certName + ".crtKey")

	// Create Layer8 certificate
	err := certs.CreateLayer8Crt(certName, "127.0.0.1", 8443)

	if err != nil {
		t.Fatalf("Failed to create Layer8 certificate: %v", err)
	}

	// Verify files were created
	if _, err := os.Stat(certName + ".ca"); os.IsNotExist(err) {
		t.Error("CA file was not created")
	}

	if _, err := os.Stat(certName + ".caKey"); os.IsNotExist(err) {
		t.Error("CA key file was not created")
	}

	if _, err := os.Stat(certName + ".crt"); os.IsNotExist(err) {
		t.Error("Certificate file was not created")
	}

	if _, err := os.Stat(certName + ".crtKey"); os.IsNotExist(err) {
		t.Error("Certificate key file was not created")
	}

	// Clean up
	os.Remove(certName + ".ca")
	os.Remove(certName + ".caKey")
	os.Remove(certName + ".crt")
	os.Remove(certName + ".crtKey")
}
