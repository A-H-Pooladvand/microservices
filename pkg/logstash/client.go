package logstash

import (
	"fmt"
	"net"
	"po/configs"
	"sync"
)

var (
	client *Client
	once   sync.Once
)

type Client struct {
	Conn net.Conn
}

func (c *Client) Shutdown() error {
	fmt.Println("Shutting down logstash client")
	return c.Conn.Close()
}

func NewSingleton() (*Client, error) {
	var err error

	c := &Client{}

	config, err := configs.NewLogstash()

	if err != nil {
		return nil, err
	}

	once.Do(func() {
		con, e := net.Dial("tcp", config.Address)

		c.Conn = con
		err = e
		client = c
	})

	return client, err
}
