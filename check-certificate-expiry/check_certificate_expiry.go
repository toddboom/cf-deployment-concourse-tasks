package check_certificate_expiry

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
)

// CheckCertificateExiries
// GetAllCredentials(credHubClient)
// for each certificateCredential
//   certificate := credHubClien.Find()
//   ignore non-certificate credentials
//   name, expiry, err := RetrieveCertificateExpirationInfo(certicicate)
//   ignore expirations > 30 days

func ExtractCertificateExpiry(certificatePEM string) (time.Duration, error) {
	if certificatePEM == "" {
		return 0, errors.New("No certificate provided")
	}

	block, _ := pem.Decode([]byte(certificatePEM))
	if block == nil {
		panic("DON'T PANIC")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		panic("DON'T PANIC")
	}

	expiry := cert.NotAfter.Sub(time.Now())
	return expiry, nil
}

// Input: Credential
// Output: name, expiry, err
func RetrieveCertificateExpirationInfo(certificate credentials.Certificate) (string, time.Duration, error) {
}

func EvaluateCredentials(credentials []credentials.Credential) (map[string]time.Duration, error) {
	return nil, nil
}

// func FindExpiringCertificates(certificates []credentials.Certificate) (map[string]time.Duration, error) {
// 	expiring := map[string]time.Duration{}
// 	for _, certificate := range certificates {
// 		expiry, _ := ExtractCertificateExpiry(certificate.Value.Certificate)
// 		if expiry < 30*24*time.Hour {
// 			expiring[certificate.Name] = expiry
// 		}
// 	}

// 	return expiring, nil
// }
