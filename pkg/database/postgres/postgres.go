package postgres

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "github.com/lib/pq"

	"github.com/clubcedille/calidum-rotae-backend/pkg/database"
)

const (
	driverName string = "postgres"
)

type PostgresClient struct {
	db *sql.DB
}

var _ database.Operations = &PostgresClient{}

type Config struct {
	User     string
	Password string
	SSLMode  string
	DbName   string
	Host     string
	Schema   string

	MaxIdleConn     int // defaults to 100
	MaxOpenConn     int // defaults to 100
	MaxConnIdleTime time.Duration
	MaxConnLifeTime time.Duration
}

func NewPostgresClient(cfg Config) (*PostgresClient, error) {
	connectionString := buildConnectionString(cfg)
	db, err := sql.Open(driverName, connectionString)
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	db.SetConnMaxIdleTime(cfg.MaxConnIdleTime)
	db.SetConnMaxLifetime(cfg.MaxConnLifeTime)

	return &PostgresClient{db}, nil
}

func buildConnectionString(cfg Config) string {
	dbInfo := strings.Builder{}
	dbInfo.WriteString(fmt.Sprintf("user=%s ", cfg.User))
	dbInfo.WriteString(fmt.Sprintf("password=%s ", cfg.Password))
	dbInfo.WriteString(fmt.Sprintf("dbname=%s ", cfg.DbName))
	dbInfo.WriteString(fmt.Sprintf("host=%s ", cfg.Host))
	dbInfo.WriteString(fmt.Sprintf("search_path=%s ", cfg.Schema))
	if cfg.SSLMode != "" {
		dbInfo.WriteString(fmt.Sprintf("sslmode=%s ", cfg.SSLMode))
	}
	return dbInfo.String()
}
