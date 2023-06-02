package test

import (
	"fmt"
	"testing"

	"github.com/AlexandreMarcq/gozimbra/internal/domain"
	client "github.com/AlexandreMarcq/gozimbra/pkg"
	"github.com/google/go-cmp/cmp"
)

func TestGetQuotaUsage(t *testing.T) {
	testToken := "TEST_TOKEN"
	testDomain := "domain.fr"
	quotas := []domain.Quota{
		{
			Account: "a@domain.fr",
			Used:    10,
			Total:   100,
		},
		{
			Account: "b@domain.fr",
			Used:    20,
			Total:   300,
		},
		{
			Account: "c@domain.fr",
			Used:    40,
			Total:   50,
		},
	}
	server := NewXMLServer(t, fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"><authToken>%s</authToken><format type="xml"/></context></soap:Header><soap:Body><GetQuotaUsageRequest xmlns="urn:zimbraAdmin" domain="%s" allServers="1"/></soap:Body></soap:Envelope>`, testToken, testDomain), fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"/></soap:Header><soap:Body><GetQuotaUsageResponse searchTotal="3" more="0" xmlns="urn:zimbraAdmin"><account name="%s" limit="%d" id="c5d699c0-81b0-4a07-8fbf-314cd8bd0e97" used="%d"/><account name="%s" limit="%d" id="9edd9aa5-6636-4374-8c52-8b65a08cec2e" used="%d"/><account name="%s" limit="%d" id="13f3f809-b9c2-4996-b93b-dde267001ee0" used="%d"/></GetQuotaUsageResponse></soap:Body></soap:Envelope>`, quotas[0].Account, quotas[0].Total, quotas[0].Used, quotas[1].Account, quotas[1].Total, quotas[1].Used, quotas[2].Account, quotas[2].Total, quotas[2].Used))
	defer server.Close()

	client := client.NewClient(server.URL)
	client.Token = testToken

	got, err := client.GetQuotaUsage(testDomain)
	if err != nil {
		t.Fatalf("error getting quota usage: %v", err)
	}

	diff := cmp.Diff(quotas, got)
	if diff != "" {
		t.Errorf("wrong values, want %v got %v", quotas, got)
	}
}
