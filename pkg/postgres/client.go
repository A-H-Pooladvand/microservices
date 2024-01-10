package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"sync"
	"time"
)

var (
	once = sync.Once{}
	c    *Client
)

type Client struct {
	config     Config
	connection *gorm.DB
}

func (c *Client) Shutdown() error {
	fmt.Println("Shutting down logstash client")

	db, err := c.connection.DB()

	if err != nil {
		return err
	}

	return db.Close()
}

func New(config Config) (*Client, error) {
	var err error

	once.Do(func() {
		ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
		defer cancel()

		c = &Client{
			config: config,
		}

		db, err := gorm.Open(c.dialect(), &gorm.Config{})

		if err != nil {
			return
		}

		c.connection = db

		connection, err := db.DB()

		if err != nil {
			return
		}

		for {
			if err = connection.Ping(); err == nil {
				break
			}

			select {
			case <-time.After(500 * time.Millisecond):
				continue
			case <-ctx.Done():
				err = errors.New("unable to connect to postgres client, context deadline exceeded")
				return
			}
		}
	})

	// There is no need to check if err != nil
	//since we return both client and err

	return c, err
}

func (c *Client) dsn() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Tehran",
		c.config.Host,
		c.config.Username,
		c.config.Password,
		c.config.DB,
		c.config.Port,
	)
}

func (c *Client) dialect() gorm.Dialector {
	return postgres.Open(c.dsn())
}

func (c *Client) DB() *sql.DB {
	db, _ := c.connection.DB()

	return db
}

func (c *Client) Ping() error {
	return c.DB().Ping()
}
