package db

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

// NewPostgresHandler instantiates a PostgresHandler struct an return its pointer
func NewPostgresHandler(options DatabaseConnectionOpts) (*sql.DB, error) {
	psqlOpts := NewPostgresConnOptions(options)
	connectionString := psqlOpts.GetConnString()

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(psqlOpts.MaxConnections)
	db.SetMaxIdleConns(psqlOpts.MaxIdleConnections)
	db.SetConnMaxLifetime(time.Duration(psqlOpts.ConnMaxLifetime) * time.Second)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// PostgresConnOptions give optional configuration options for db connecton
type PostgresConnOptions struct {
	DatabaseConnectionOpts
	SSLEnabled bool
}

// NewPostgresConnOptions constructor function for postgres specific connection options
func NewPostgresConnOptions(
	dbConnectionOpts DatabaseConnectionOpts,
) *PostgresConnOptions {
	return &PostgresConnOptions{
		DatabaseConnectionOpts: dbConnectionOpts,
		SSLEnabled:             false,
	}
}

// GetConnString returns connection string postgres database
func (opts *PostgresConnOptions) GetConnString() string {
	if opts.Port == "" {
		opts.Port = "5432"
	}
	if opts.Timeout < 0 {
		opts.Timeout = 0
	}
	sslMode := "disable"
	if opts.SSLEnabled {
		sslMode = "require"
	}
	return fmt.Sprintf(
		"host=%v dbname=%v port=%v user=%v password='%v' connect_timeout=%d sslmode=%v",
		opts.Host,
		opts.Database,
		opts.Port,
		opts.User,
		opts.Password,
		opts.Timeout,
		sslMode,
	)
}
