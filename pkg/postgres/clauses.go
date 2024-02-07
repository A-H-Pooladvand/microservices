package postgres

func (c *Client) Where(query any, args ...any) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connection = c.connection.Where(query, args...)

	return c
}

func (c *Client) Create(value any) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connection = c.connection.Create(value)

	return c
}

func (c *Client) Limit(limit int) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connection = c.connection.Limit(limit)

	return c
}

func (c *Client) Offset(offset int) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connection = c.connection.Offset(offset)

	return c
}

func (c *Client) Find(dest any, conds ...any) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connection = c.connection.Find(dest, conds...)

	return c
}

func (c *Client) First(dest any, conds ...any) *Client {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.connection = c.connection.First(dest, conds...)

	return c
}
