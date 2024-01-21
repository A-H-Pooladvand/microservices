package logstash

import (
	"fmt"
	"net"
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

	once.Do(func() {
		con, e := net.Dial("tcp", "127.0.0.1:5000")

		c.Conn = con
		err = e
		client = c
	})

	if err != nil {
		panic(err)
	}

	return client, err
}

func Get() *Client {
	if client == nil {
		_, _ = NewSingleton()
	}

	return client
}
