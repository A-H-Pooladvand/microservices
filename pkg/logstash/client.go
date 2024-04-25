package logstash

import (
	"context"
	"github.com/fatih/color"
	"go.uber.org/fx"
	"net"
	"po/configs"
)

type Client struct {
	conn   net.Conn
	Config *configs.Logstash
}

func New(lc fx.Lifecycle, config *configs.Logstash) *Client {
	c := &Client{
		Config: config,
	}

	lc.Append(fx.Hook{
		OnStart: c.OnStart,
		OnStop:  c.Shutdown,
	})

	return c
}

func (c *Client) OnStart(_ context.Context) error {
	conn, err := net.Dial("tcp", c.Config.Address)

	// Todo:: Handle production state, make a mechanism to dynamically turn on/off services
	// Since the application should not go down due to log failure, we won't panic
	if err != nil {
		conn1, conn2 := net.Pipe()
		conn2.Close()

		c.conn = conn1
		color.Red("unable to connect to logstash: %v", err.Error())
	} else {
		c.conn = conn
	}

	return nil
}

func (c *Client) Shutdown(_ context.Context) error {
	return c.conn.Close()
}

func (c *Client) Connection() net.Conn {
	return c.conn
}
