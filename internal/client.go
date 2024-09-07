package internal

import (
	cmap "github.com/orcaman/concurrent-map/v2"
	"time"
)

type Client struct {
	lastActivity time.Time
	Data         cmap.ConcurrentMap[string, string]
}

func NewClient() *Client {
	return &Client{
		Data: cmap.New[string](),
	}
}

func (c *Client) SetValue(key string, value string) {
	c.Data.Set(key, value)
}

func (c *Client) GetValue(key string) (string, bool) {
	return c.Data.Get(key)
}
