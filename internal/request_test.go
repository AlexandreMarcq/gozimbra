package request_test

import (
	"fmt"
	"testing"

	request "github.com/AlexandreMarcq/gozimbra/internal"
)

func TestAddEnvelope(t *testing.T) {
	header := "HEADER"
	body := "BODY"
	want := fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">%s%s</soap:Envelope>`, header, body)

	got := request.AddEnvelope(header, body)
	if got != want {
		t.Errorf("incorrect envelope, want %s got %s", want, got)
	}
}

func TestAddHeader(t *testing.T) {
	t.Run("token is not empty", func(t *testing.T) {
		token := "TOKEN"
		want := fmt.Sprintf(`<soap:Header><context xmlns="urn:zimbra"><authToken>%s</authToken><format type="xml"/></context></soap:Header>`, token)

		got := request.AddHeader(token)
		if got != want {
			t.Errorf("incorrect header, want %s got %s", want, got)
		}
	})

	t.Run("token is empty", func(t *testing.T) {
		want := `<soap:Header><context xmlns="urn:zimbra"><format type="xml"/></context></soap:Header>`

		got := request.AddHeader("")
		if got != want {
			t.Errorf("incorrect header, want %s got %s", want, got)
		}
	})
}

func TestAddBody(t *testing.T) {
	content := "TEST"
	want := fmt.Sprintf(`<soap:Body>%s</soap:Body>`, content)

	got := request.AddBody(content)
	if got != want {
		t.Errorf("incorrect body, want %s got %s", want, got)
	}
}
