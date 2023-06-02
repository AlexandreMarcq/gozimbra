package client

import (
	account "github.com/AlexandreMarcq/gozimbra/internal/account"
	"github.com/AlexandreMarcq/gozimbra/internal/utils"
)

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

func (c *Client) ModifyAccount(name string, attributes utils.AttrsMap) (utils.AttrsMap, error) {
	if err := c.checkToken(); err != nil {
		return nil, err
	}

	id, err := c.GetAccount(name, []string{"zimbraId"})
	if err != nil {
		return nil, err
	}

	body, err := c.send(account.NewModifyAccountRequest(id["zimbraId"], attributes))

	if err != nil {
		return nil, err
	}

	res, err := parseAttributes(body, attributes.Keys())
	if err != nil {
		return nil, err
	}

	return res, nil
}
