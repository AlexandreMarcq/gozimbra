package client

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	req "github.com/AlexandreMarcq/gozimbra/internal"
	"github.com/AlexandreMarcq/gozimbra/internal/utils"
)

type Client struct {
	url   string
	Token string
}

func NewClient(url string) *Client {
	return &Client{url: url}
}

func (c *Client) checkToken() error {
	if c.Token == "" {
		return errors.New("client is not authentified")
	}
	return nil
}

func (c *Client) send(request req.Request) ([]byte, error) {
	xmlHeader := req.AddHeader(c.Token)
	xmlBody := req.AddBody(request.ToXML())

	reader := strings.NewReader(req.AddEnvelope(xmlHeader, xmlBody))

	resp, err := http.Post(c.url, "application/xml", reader)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func handleError(body []byte) error {
	r := regexp.MustCompile(`<faultstring>(.*)</faultstring>`)

	if r.Match(body) {
		return errors.New(string(r.FindSubmatch(body)[1]))
	}

	return fmt.Errorf("body is unexpected: %v", string(body))
}

func parseAttributes(body []byte, attributes []string) (utils.AttrsMap, error) {
	attrsReg := regexp.MustCompile(`<a n="(.*?)">(.*?)</a>`)
	emptyReg := regexp.MustCompile(`<account name=".*" id=".*"/></GetAccountResponse>`)

	res := make(utils.AttrsMap)

	if attrsReg.Match(body) {
		for _, attr := range attrsReg.FindAllStringSubmatch(string(body), -1) {
			res[attr[1]] = attr[2]
		}
	} else if !emptyReg.Match(body) {
		return nil, handleError(body)
	}

	for _, attr := range attributes {
		if res[attr] == "" {
			res[attr] = "N/A"
		}
	}

	return res, nil
}
