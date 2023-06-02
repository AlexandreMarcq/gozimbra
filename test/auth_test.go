package test

import (
	"fmt"
	"testing"

	client "github.com/AlexandreMarcq/gozimbra/pkg"
)

func TestAuth(t *testing.T) {
	token := "TEST_TOKEN"
	testUser := "user"
	testPass := "password"
	server := NewXMLServer(t, fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"><format type="xml"/></context></soap:Header><soap:Body><AuthRequest xmlns="urn:zimbraAdmin" name="%s" password="%s"/></soap:Body></soap:Envelope>`, testUser, testPass),
		fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/"><soap:Header><context xmlns="urn:zimbra"><change token="CHANGE_TOKEN"/></context></soap:Header><soap:Body><AuthResponse xmlns="urn:zimbraAdmin"><lifetime>LIFETIME</lifetime><authToken>%s</authToken></AuthResponse></soap:Body></soap:Envelope>`, token))
	defer server.Close()

	c := client.NewClient(server.URL)

	err := c.Auth(testUser, testPass)
	if err != nil {
		t.Fatalf("error while authenticating: %s", err)
	}

	if c.Token == "" {
		t.Fatal("token is empty")
	}

	if c.Token != token {
		t.Fatalf("wrong token, wanted %s got %s", token, c.Token)
	}
}
