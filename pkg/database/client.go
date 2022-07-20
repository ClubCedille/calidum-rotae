package database

type Operations interface {
	Ping() error
}
