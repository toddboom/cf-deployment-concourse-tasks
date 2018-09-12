package check_certificate_expiry_test

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"

	"code.cloudfoundry.org/credhub-cli/credhub/credentials"
	"code.cloudfoundry.org/credhub-cli/credhub/credentials/values"
	. "github.com/cloudfoundry/cf-deployment-concourse-tasks/check-certificate-expiry"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CheckCertificateExpiry", func() {
	// cannot decode
	// cannot parse
	// already expired

	Context("#ExtractCertificateExpiry", func() {
		It("returns error when certificate is not provided", func() {
			_, err := ExtractCertificateExpiry("")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("No certificate provided"))
		})

		It("returns number of days until expiration for a certificate expiring in 2 days", func() {
			// todo: generate it on the fly
			certificateExpiringInTwoDays, err := generateCertificate(2 * 24 * time.Hour)
			Expect(err).NotTo(HaveOccurred())

			expiry, err := ExtractCertificateExpiry(certificateExpiringInTwoDays)
			Expect(err).ToNot(HaveOccurred())
			Expect(expiry).To(BeNumerically("<=", 2*24*time.Hour))
		})
	})

	Context("#FindExpiringCertificates", func() {
		It("returns a list of expiring certificates", func() {
			certificateExpiringInTwoDays, _ := generateCertificate(2 * 24 * time.Hour)
			certificateExpiringInTenDays, _ := generateCertificate(10 * 24 * time.Hour)
			certificateExpiringInAYear, _ := generateCertificate(365 * 24 * time.Hour)

			certificatesToTest := []credentials.Certificate{
				credentials.Certificate{
					Metadata: credentials.Metadata{Base: credentials.Base{Name: "expiringInTwoDays"}},
					Value: values.Certificate{
						Certificate: certificateExpiringInTwoDays,
					},
				},
				credentials.Certificate{
					Metadata: credentials.Metadata{Base: credentials.Base{Name: "expiringInTenDays"}},
					Value: values.Certificate{
						Certificate: certificateExpiringInTenDays,
					},
				},
				credentials.Certificate{
					Metadata: credentials.Metadata{Base: credentials.Base{Name: "expiringInAYear"}},
					Value: values.Certificate{
						Certificate: certificateExpiringInAYear,
					},
				},
			}

			expiringCerts, err := FindExpiringCertificates(certificatesToTest)
			Expect(err).NotTo(HaveOccurred())
			Expect(expiringCerts["expiringInTwoDays"]).To(BeNumerically("<=", 2*24*time.Hour))
			Expect(expiringCerts["expiringInTenDays"]).To(BeNumerically("<=", 10*24*time.Hour))
		})
	})
})

func generateCertificate(expiresIn time.Duration) (string, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", err
	}

	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"test"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(expiresIn),

		BasicConstraintsValid: true,
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return "", err
	}
	out := &bytes.Buffer{}
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	return out.String(), nil
}
