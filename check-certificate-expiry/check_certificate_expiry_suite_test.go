package check_certificate_expiry_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestCheckCertificateExpiry(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CheckCertificateExpiry Suite")
}
