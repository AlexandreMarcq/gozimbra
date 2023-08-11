package client

import (
	"errors"
	"regexp"
	"strings"

	account "github.com/AlexandreMarcq/gozimbra/internal/account"
	"github.com/AlexandreMarcq/gozimbra/internal/utils"
)

func (c *Client) GetId(name string) (string, error) {
	id, err := c.GetAccount(name, []string{"zimbraId"})
	if err != nil {
		return "", err
	}
	return id["zimbraId"], nil
}

func (c *Client) GetAccount(name string, attributes []string) (utils.AttrsMap, error) {
	if err := c.checkToken(); err != nil {
		return nil, err
	}

	body, err := c.send(account.NewGetAccountRequest(name, attributes))
	if err != nil {
		return nil, err
	}

	res, err := parseAttributes(body, attributes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) ModifyAccount(id string, attributes utils.AttrsMap) (utils.AttrsMap, error) {
	if err := c.checkToken(); err != nil {
		return nil, err
	}

	body, err := c.send(account.NewModifyAccountRequest(id, attributes))
	if err != nil {
		return nil, err
	}

	res, err := parseAttributes(body, attributes.Keys())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *Client) SetPassword(id, newPassword string) (string, error) {
	if err := c.checkToken(); err != nil {
		return "", err
	}

	body, err := c.send(account.NewSetPasswordRequest(id, newPassword))
	if err != nil {
		return "", err
	}

	r := regexp.MustCompile(`<message>(.*)</message>`)

	if r.Match(body) {
		return "", errors.New(string(r.FindSubmatch(body)[1]))
	} else if strings.Contains(string(body), "faultstring") {
		return "", handleError(body)
	} else {
		return "Password successfully changed", nil
	}
}
