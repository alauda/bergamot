package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQLHandler instantiates a MySQLHandler struct an return its pointer
func NewMySQLHandler(options DatabaseConnectionOpts) (*sql.DB, error) {
	mySQLOpts := NewMySQLConnOptions(options)
	connectionString := mySQLOpts.GetConnString()
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(options.MaxConnections)
	db.SetMaxIdleConns(options.MaxIdleConnections)
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

// MySQLConnectionOpts connection options for MySQL
type MySQLConnectionOpts struct {
	DatabaseConnectionOpts
	Protocol string
}

// NewMySQLConnOptions constructor function for mysql specific connection options
func NewMySQLConnOptions(
	dbConnOpts DatabaseConnectionOpts,
) *MySQLConnectionOpts {
	return &MySQLConnectionOpts{
		DatabaseConnectionOpts: dbConnOpts,
		Protocol:               "tcp",
	}
}

// GetConnString returns a connection string used to connect to a mysql database
func (opts *MySQLConnectionOpts) GetConnString() string {
	if opts.Params == nil {
		opts.Params = make(map[string]string)
	}
	if opts.Port == "" {
		opts.Port = "3306"
	}
	if opts.Protocol == "" {
		opts.Protocol = "tcp"
	}

	// Parsing DATE TIME to time.Time
	//https://github.com/Go-SQL-Driver/MySQL/
	opts.Params["parseTime"] = "true"
	// opts.Params["timeout"] = (time.Duration(opts.Timeout) * time.Second).String()

	return fmt.Sprintf(
		//username:password@host/database?params
		"%s:%s@%s/%s?%s",
		opts.User,
		opts.Password,
		opts.getFullHost(),
		opts.Database,
		opts.GetParams(),
	)
}

func (opts *MySQLConnectionOpts) getFullHost() string {
	return fmt.Sprintf(
		// protocol(host:port)
		"%s(%s:%s)",
		opts.Protocol,
		opts.Host,
		opts.Port,
	)
}
