package test

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"testing"

	"github.com/AlexandreMarcq/gozimbra/internal/utils"
	client "github.com/AlexandreMarcq/gozimbra/pkg"
)

func TestGetAccount(t *testing.T) {
	token := "AUTH_TOKEN"

	t.Run("getting a value", func(t *testing.T) {
		attrs := make(utils.AttrsMap)
		attrs["name"] = "NAME"
		attrs["description"] = "DESCRIPTION"

		server := NewXMLServer(t, fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"><authToken>%s</authToken><format type="xml"/></context></soap:Header><soap:Body><GetAccountRequest xmlns="urn:zimbraAdmin" attrs="name,description"><account by="name">USER_ACCOUNT</account></GetAccountRequest></soap:Body></soap:Envelope>`, token), `<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"/></soap:Header><soap:Body><GetAccountResponse xmlns="urn:zimbraAdmin"><account name="USER_ACCOUNT"><a n="name">NAME</a><a n="description">DESCRIPTION</a></account></GetAccountResponse></soap:Body></soap:Envelope>`)
		defer server.Close()

		client := client.NewClient(server.URL)
		client.Token = token

		got, err := client.GetAccount("USER_ACCOUNT", []string{"name", "description"})
		if err != nil {
			t.Fatalf("error getting account information: %v", err)
		}

		if got["name"] != attrs["name"] || got["description"] != attrs["description"] {
			t.Fatalf("wrong values, wanted %v got %v", attrs, got)
		}
	})

	t.Run("getting an error on non-existing account", func(t *testing.T) {
		fakeAccount := "NO_ACCOUNT"
		errMessage := fmt.Sprintf("no such account: %s", fakeAccount)
		server := NewXMLServer(t, fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"><authToken>%s</authToken><format type="xml"/></context></soap:Header><soap:Body><GetAccountRequest xmlns="urn:zimbraAdmin" attrs="name,description"><account by="name">%s</account></GetAccountRequest></soap:Body></soap:Envelope>`, token, fakeAccount), fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"/></soap:Header><soap:Body><soap:Fault><faultcode>soap:Client</faultcode><faultstring>%s</faultstring><detail><Error xmlns="urn:zimbra"><Code>account.NO_SUCH_ACCOUNT</Code><Trace>TRACE</Trace></Error></detail></soap:Fault></soap:Body></soap:Envelope>`, errMessage))
		defer server.Close()

		client := client.NewClient(server.URL)
		client.Token = token

		_, err := client.GetAccount(fakeAccount, []string{"name", "description"})
		AssertError(t, err)

		if err.Error() != errMessage {
			t.Errorf("wrong error message, want \"%s\" got \"%s\"", errMessage, err)
		}
	})
}

func TestModifyAccount(t *testing.T) {
	token := "AUTH_TOKEN"
	account := "TEST_ACCOUNT"
	id := "TEST_ID"
	attrs := make(utils.AttrsMap)
	attrs["ATTR1"] = "VALUE1"
	server := NewFakeServer(t, func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("error reading body: %v", err)
		}

		getAccountReg := regexp.MustCompile(fmt.Sprintf(`<GetAccountRequest xmlns="urn:zimbraAdmin" attrs="zimbraId"><account by="name">%s</account></GetAccountRequest>`, account))
		modifyAccountReg := regexp.MustCompile(fmt.Sprintf(`<ModifyAccountRequest xmlns="urn:zimbraAdmin" id="%s"><a n="ATTR1">VALUE1</a></ModifyAccountRequest>`, id))

		if getAccountReg.Match(body) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"/></soap:Header><soap:Body><GetAccountResponse xmlns="urn:zimbraAdmin"><account name="%s"><a n="zimbraId">%s</a></account></GetAccountResponse></soap:Body></soap:Envelope>`, account, id)))
		} else if modifyAccountReg.Match(body) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"/></soap:Header><soap:Body><ModifyAccountResponse xmlns="urn:zimbraAdmin" id="%s"><a n="ATTR1">VALUE1</a></ModifyAccountResponse></soap:Body></soap:Envelope>`, id)))
		} else {
			w.WriteHeader(http.StatusBadRequest)
			t.Fatalf("wrong body: %v", string(body))
		}
	})
	defer server.Close()

	client := client.NewClient(server.URL)
	client.Token = token

	got, err := client.ModifyAccount(account, attrs)
	if err != nil {
		t.Fatalf("error modifying account: %s", err)
	}

	if got["ATTR1"] != "VALUE1" {
		t.Errorf("did not modify correct attribute, wanted 'VALUE1' got %s", got["ATTR1"])
	}
}
