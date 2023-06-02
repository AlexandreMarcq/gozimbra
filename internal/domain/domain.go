package domain

import "fmt"

type Quota struct {
	Account string
	Used    uint
	Total   uint
}

type getQuotaUsageRequest struct {
	domain string
}

func NewGetQuotaUsageRequest(domain string) *getQuotaUsageRequest {
	return &getQuotaUsageRequest{domain}
}

func (r *getQuotaUsageRequest) ToXML() string {
	return fmt.Sprintf(`<GetQuotaUsageRequest xmlns="urn:zimbraAdmin" domain="%s" allServers="1"/>`, r.domain)
}
