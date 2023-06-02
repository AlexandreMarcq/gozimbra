package request

import "fmt"

type Request interface {
	ToXML() string
}

func AddEnvelope(header, body string) string {
	return fmt.Sprintf(`<soap:Envelope xmlns:soap="http://schemas.xmlsoap.org/soap/envelope/">%s%s</soap:Envelope>`, header, body)
}

func AddHeader(token string) string {
	if token == "" {
		return `<soap:Header><context xmlns="urn:zimbra"><format type="xml"/></context></soap:Header>`
	} else {
		return fmt.Sprintf(`<soap:Header><context xmlns="urn:zimbra"><authToken>%s</authToken><format type="xml"/></context></soap:Header>`, token)
	}
}

func AddBody(content string) string {
	return fmt.Sprintf(`<soap:Body>%s</soap:Body>`, content)
}
