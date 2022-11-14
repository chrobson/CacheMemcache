package main

import (
	"bytes"
	"encoding/gob"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
)

type Client struct {
	client *memcache.Client
}

func NewMemcached() (*Client, error) {
	// using local docker memache at localhost:11211
	client := memcache.New("localhost:11211")

	if err := client.Ping(); err != nil {
		return nil, err
	}

	client.Timeout = 100 * time.Millisecond
	client.MaxIdleConns = 100

	return &Client{
		client: client,
	}, nil
}

func (c *Client) GetPerson(id string) (Person, error) {
	item, err := c.client.Get(id)
	if err != nil {
		return Person{}, err
	}

	b := bytes.NewReader(item.Value)

	var res Person

	if err := gob.NewDecoder(b).Decode(&res); err != nil {
		return Person{}, err
	}

	return res, nil
}

func (c *Client) SetPerson(n Person) error {
	var b bytes.Buffer

	if err := gob.NewEncoder(&b).Encode(n); err != nil {
		return err
	}

	return c.client.Set(&memcache.Item{
		Key:        n.Id,
		Value:      b.Bytes(),
		Expiration: int32(time.Now().Add(25 * time.Second).Unix()),
	})
}
