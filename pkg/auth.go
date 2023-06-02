package client

import (
	"regexp"

	"github.com/AlexandreMarcq/gozimbra/internal/auth"
)

func (c *Client) Auth(user, password string) error {
	body, err := c.send(auth.NewAuthRequest(user, password))
	if err != nil {
		return err
	}

	r := regexp.MustCompile(`<authToken>(.*)</authToken>`)

	if r.Match(body) {
		c.Token = string(r.FindSubmatch(body)[1])
		return nil
	} else {
		return handleError(body)
	}
}
