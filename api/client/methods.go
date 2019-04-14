package client

import (
	"net/http"

	"github.com/alexxeis/keyval/api"
)

func (c *Client) Keys() ([]string, error) {
	req, err := c.newRequest(http.MethodGet, "/keys", nil)
	if err != nil {
		return nil, err
	}

	var keys []string
	err = c.process(req, &keys)
	return keys, err
}

func (c *Client) Get(key string) (string, error) {
	req, err := c.newRequest(http.MethodGet, "/get/"+key, nil)
	if err != nil {
		return "", err
	}

	val := &api.Value{}
	err = c.process(req, val)
	return val.Value, err
}

func (c *Client) Set(key string, params *api.SetParams) error {
	req, err := c.newRequest(http.MethodPost, "/set/"+key, params)
	if err != nil {
		return err
	}
	return c.process(req, nil)
}

func (c *Client) Expire(key string, ttl *api.Ttl) error {
	req, err := c.newRequest(http.MethodPost, "/expire/"+key, ttl)
	if err != nil {
		return err
	}
	return c.process(req, nil)
}

func (c *Client) Remove(key string) error {
	req, err := c.newRequest(http.MethodPost, "/remove/"+key, nil)
	if err != nil {
		return err
	}
	return c.process(req, nil)
}

func (c *Client) Hget(key, field string) (string, error) {
	req, err := c.newRequest(http.MethodGet, "/hget/"+key+"/"+field, nil)
	if err != nil {
		return "", err
	}

	val := &api.Value{}
	err = c.process(req, val)
	return val.Value, err
}

func (c *Client) Hset(key, field, value string) error {
	req, err := c.newRequest(http.MethodPost, "/hset/"+key+"/"+field, api.Value{Value: value})
	if err != nil {
		return err
	}
	return c.process(req, nil)
}

func (c *Client) Hdel(key, field string) error {
	req, err := c.newRequest(http.MethodPost, "/hdel/"+key+"/"+field, nil)
	if err != nil {
		return err
	}
	return c.process(req, nil)
}
