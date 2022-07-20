package postgres

func (c PostgresClient) Ping() error {
	return c.db.Ping()
}
