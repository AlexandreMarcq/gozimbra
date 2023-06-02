package client

import (
	"fmt"
	"regexp"
	"strconv"

	dom "github.com/AlexandreMarcq/gozimbra/internal/domain"
)

func (c *Client) GetQuotaUsage(domain string) ([]dom.Quota, error) {
	if err := c.checkToken(); err != nil {
		return nil, err
	}

	body, err := c.send(dom.NewGetQuotaUsageRequest(domain))
	if err != nil {
		return nil, err
	}

	res := []dom.Quota{}

	r := regexp.MustCompile(`<account name="(.*?)" limit="(\d*?)" id=".*?" used="(\d*?)"/>`)
	if r.Match(body) {
		for _, quota := range r.FindAllStringSubmatch(string(body), -1) {
			used, err := strconv.ParseUint(quota[3], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("could not convert usage to uint, got %v", quota[3])
			}

			total, err := strconv.ParseUint(quota[2], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("could not convert total to uint, got %v", quota[2])
			}

			res = append(res, dom.Quota{
				Account: quota[1],
				Used:    uint(used),
				Total:   uint(total),
			})
		}
	} else {
		return nil, handleError(body)
	}

	return res, nil
}
