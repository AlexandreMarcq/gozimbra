package domain_test

import (
	"fmt"
	"testing"

	"github.com/AlexandreMarcq/gozimbra/internal/domain"
	"github.com/AlexandreMarcq/gozimbra/test"
)

func TestGetQuotaUsageRequest(t *testing.T) {
	testDomain := "domain.fr"
	want := fmt.Sprintf(`<GetQuotaUsageRequest xmlns="urn:zimbraAdmin" domain="%s" allServers="1"/>`, testDomain)
	got := domain.NewGetQuotaUsageRequest(testDomain).ToXML()

	test.AssertXML(t, got, want)
}
