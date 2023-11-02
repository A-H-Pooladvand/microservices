package mysql

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

var (
	once = sync.Once{}
	db   *gorm.DB
)

type Client struct {
	config Config
	db     *gorm.DB
}

func New(config Config) *Client {
	c := &Client{
		config: config,
	}

	client, err := c.connect()

	if err != nil {
		panic(err)
	}

	c.db = client

	return c
}

func (c *Client) connect() (*gorm.DB, error) {
	var err error
	once.Do(func() {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       c.dsn(),
			DefaultStringSize:         256,
			DisableDatetimePrecision:  true,
			DontSupportRenameIndex:    true,
			DontSupportRenameColumn:   true,
			SkipInitializeWithVersion: false,
		}), &gorm.Config{})
	})

	return db, err
}

func (c *Client) dsn() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		c.config.Username,
		c.config.Password,
		c.config.Host,
		c.config.Port,
		c.config.Database,
	)
}

func (c *Client) Close() error {
	sql, err := c.db.DB()

	if err != nil {
		return err
	}

	return sql.Close()
}
