package internal

import "sync"

var instance *Coordinator
var once sync.Once

type Coordinator struct {
	Clients map[string]*Client
}

func GetCoordinator() *Coordinator {
	once.Do(func() {
		instance = &Coordinator{
			make(map[string]*Client),
		}
	})

	return instance
}

func (c *Coordinator) GetInstance(key string) *Client {

	client, exists := c.Clients[key]

	if !exists {
		client = NewClient()
		c.Clients[key] = client
	}

	return client
}
