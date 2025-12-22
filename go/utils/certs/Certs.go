// © 2025 Sharon Aicler (saichler@gmail.com)
//
// Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package certs provides TLS/SSL certificate generation utilities for secure communication.
// It supports creating self-signed Certificate Authorities (CA) and signed certificates
// with RSA key pairs for server and client authentication.
//
// Key features:
//   - Self-signed CA creation with 4096-bit RSA keys
//   - Certificate generation signed by CA
//   - PEM-encoded file output for certificates and private keys
//   - Pre-configured Layer 8 certificate generation helpers
package certs

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"math/big"
	"net"
	"os"
	"time"
)

// CreateCA generates a self-signed Certificate Authority with the specified metadata.
// Creates {filenamePrefix}.ca and {filenamePrefix}.caKey files containing the PEM-encoded
// certificate and private key. Returns an error if the CA already exists.
func CreateCA(filenamePrefix, org, country, county, city, street, zipcode, email string, years int) (*x509.Certificate, *rsa.PrivateKey, error) {
	_, e := os.Stat(filenamePrefix + ".ca")
	if e != nil {
		ca := &x509.Certificate{
			SerialNumber: big.NewInt(2025),
			Subject: pkix.Name{
				CommonName:    "www.layer8vibe.dev",
				Organization:  []string{org},
				Country:       []string{country},
				Province:      []string{county},
				Locality:      []string{city},
				StreetAddress: []string{street},
				PostalCode:    []string{zipcode},
			},
			NotBefore:             time.Now(),
			NotAfter:              time.Now().AddDate(years, 0, 0),
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
			BasicConstraintsValid: true,
			EmailAddresses:        []string{email},
		}

		caKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return nil, nil, err
		}

		caData, err := x509.CreateCertificate(rand.Reader, ca, ca, &caKey.PublicKey, caKey)
		if err != nil {
			return nil, nil, err
		}

		caPEM := &bytes.Buffer{}
		pem.Encode(caPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: caData,
		})

		err = ioutil.WriteFile(filenamePrefix+".ca", caPEM.Bytes(), 0777)
		if err != nil {
			return nil, nil, err
		}

		caKeyPEM := &bytes.Buffer{}
		err = pem.Encode(caKeyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(caKey),
		})
		if err != nil {
			return nil, nil, err
		}
		err = ioutil.WriteFile(filenamePrefix+".caKey", caKeyPEM.Bytes(), 0777)
		return ca, caKey, err
	} else {
		return nil, nil, errors.New("Certificate Authority " + filenamePrefix + " already exists!")
	}
}

// CreateCrt generates a certificate signed by the provided CA for the specified IP address.
// Creates {filenamePrefix}.crt and {filenamePrefix}.crtKey files. The port is used as
// the certificate serial number. Returns an error if the certificate already exists.
func CreateCrt(filenamePrefix, org, country, county, city, street, zipcode, email, ip, secret string, port int64, years int, ca *x509.Certificate, caKey *rsa.PrivateKey) error {
	_, e := os.Stat(filenamePrefix + ".crt")
	if e != nil {
		ipAddress := net.ParseIP(ip)
		crt := &x509.Certificate{
			SerialNumber: big.NewInt(port),
			Subject: pkix.Name{
				CommonName:    "www.layer8vibe.dev",
				Organization:  []string{org},
				Country:       []string{country},
				Province:      []string{county},
				Locality:      []string{city},
				StreetAddress: []string{street},
				PostalCode:    []string{zipcode},
			},
			EmailAddresses: []string{email},
			IPAddresses:    []net.IP{ipAddress},
			NotBefore:      time.Now(),
			NotAfter:       time.Now().AddDate(years, 0, 0),
			SubjectKeyId:   []byte(secret),
			ExtKeyUsage:    []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:       x509.KeyUsageDigitalSignature,
		}

		key, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return err
		}

		crtData, err := x509.CreateCertificate(rand.Reader, crt, ca, &key.PublicKey, caKey)
		if err != nil {
			return err
		}

		crtPEM := new(bytes.Buffer)
		pem.Encode(crtPEM, &pem.Block{
			Type:  "CERTIFICATE",
			Bytes: crtData,
		})

		err = ioutil.WriteFile(filenamePrefix+".crt", crtPEM.Bytes(), 0777)
		if err != nil {
			return err
		}

		keyPEM := new(bytes.Buffer)
		pem.Encode(keyPEM, &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		})
		err = ioutil.WriteFile(filenamePrefix+".crtKey", keyPEM.Bytes(), 0777)
		return err
	} else {
		return errors.New("Certificate " + filenamePrefix + " already exists!")
	}
}

// CreateLayer8CA creates a CA with pre-configured Layer 8 organization details.
func CreateLayer8CA(certName string) (*x509.Certificate, *rsa.PrivateKey, error) {
	return CreateCA(certName, "Layer8", "USA", "Santa Clara",
		"San Jose", "1993 Curtner Ave", "95124", "saichler@gmail.com", 10)
}

// CreateLayer8Crt creates both a Layer 8 CA and a certificate for the specified host/port.
func CreateLayer8Crt(certName, host string, port int64) error {
	ca, caKey, err := CreateLayer8CA(certName)
	if err != nil {
		return err
	}
	return CreateCrt(certName, "Layer8", "USA", "Santa Clara",
		"San Jose", "1993 Curtner Ave", "95124", "saichler@gmail.com", host, "Layer8Secret", port, 10, ca, caKey)
}
