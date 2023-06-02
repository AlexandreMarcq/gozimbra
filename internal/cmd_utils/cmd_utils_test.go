package cmd_utils_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/AlexandreMarcq/gozimbra/internal/cmd_utils"
	client "github.com/AlexandreMarcq/gozimbra/pkg"
	"github.com/AlexandreMarcq/gozimbra/test"
	"github.com/google/go-cmp/cmp"
)

func TestReadInput(t *testing.T) {
	fakeFile := new(bytes.Buffer)
	fakeFile.WriteString(`test@test.fr
test1@test.com`)

	fakeServer := test.NewFakeServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal("error while reading body")
		}

		var list string
		reg := regexp.MustCompile(`domain="(.*)" allServers`)
		if reg.Match(body) {
			domain := reg.FindStringSubmatch(string(body))
			if domain[1] == "domain.fr" {
				list = `<account name="a@domain.fr" limit="10" id="a" used="100"/>
			<account name="b@domain.fr" limit="20" id="b" used="200"/>`
			} else if domain[1] == "domain.com" {
				list = `<account name="c@domain.com" limit="30" id="c" used="300"/>
			<account name="d@domain.com" limit="40" id="d" used="400"/>`
			} else {
				w.WriteHeader(http.StatusBadRequest)
				t.Fatalf("wrong domain, got %v", domain[1])
			}
		} else {
			w.WriteHeader(http.StatusBadRequest)
			t.Fatal("no domain in body")
		}

		resp := fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">
	<soap:Header>
		<context
			xmlns="urn:zimbra"/>
		</soap:Header>
	<soap:Body>
		<GetQuotaUsageResponse searchTotal="6" more="0" xmlns="urn:zimbraAdmin">
			%s
		</GetQuotaUsageResponse>
	</soap:Body>
</soap:Envelope>`, list)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(resp))
	})
	fakeClient := client.NewClient(fakeServer.URL)
	fakeClient.Token = "TEST_TOKEN"

	tests := map[string]struct {
		client    *client.Client
		accounts  []string
		domains   []string
		inputFile io.Reader
		shouldErr bool
		want      []string
	}{
		"at least one parameter":      {client: nil, accounts: nil, domains: nil, inputFile: nil, shouldErr: true, want: []string{}},
		"one account is given":        {client: nil, accounts: []string{"test@test.fr"}, domains: nil, inputFile: nil, shouldErr: false, want: []string{"test@test.fr"}},
		"multiple accounts are given": {client: nil, accounts: []string{"test@test.fr", "test1@test.fr"}, domains: nil, inputFile: nil, shouldErr: false, want: []string{"test@test.fr", "test1@test.fr"}},
		"one domain is given":         {client: fakeClient, accounts: nil, domains: []string{"domain.fr"}, inputFile: nil, shouldErr: false, want: []string{"a@domain.fr", "b@domain.fr"}},
		"multiple domains are given":  {client: fakeClient, accounts: nil, domains: []string{"domain.fr", "domain.com"}, inputFile: nil, shouldErr: false, want: []string{"a@domain.fr", "b@domain.fr", "c@domain.com", "d@domain.com"}},
		"one file is given":           {client: nil, accounts: nil, domains: nil, inputFile: fakeFile, shouldErr: false, want: []string{"test@test.fr", "test1@test.com"}},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.shouldErr {
				_, err := cmd_utils.ReadInput(tc.client, tc.accounts, tc.domains, tc.inputFile)

				test.AssertError(t, err)
			} else {
				got, err := cmd_utils.ReadInput(tc.client, tc.accounts, tc.domains, tc.inputFile)

				test.AssertNoError(t, err)

				diff := cmp.Diff(tc.want, got)
				if diff != "" {
					t.Error(diff)
				}
			}
		})
	}
}
