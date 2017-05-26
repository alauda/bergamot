package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"time"

	"github.com/alauda/bergamot/diagnose"

	goqu "gopkg.in/doug-martin/goqu.v4"
)

// Engine database engine
type Engine string

const (
	// MySQL mysql database
	MySQL Engine = "mysql"
	// Postgres postgres database
	Postgres Engine = "postgres"
)

func (en Engine) String() string {
	return string(en)
}

// New initiates and returns a goqu database
func New(engine Engine, options DatabaseConnectionOpts) (*goqu.Database, error) {
	var (
		db  *sql.DB
		err error
	)
	switch engine {
	case MySQL:
		db, err = NewMySQLHandler(options)
	case Postgres:
		fallthrough
	default:
		engine = Postgres
		db, err = NewPostgresHandler(options)
	}
	if err != nil {
		return nil, err
	}
	return goqu.New(engine.String(), db), nil
}

// DatabaseConnectionOpts common connection options
type DatabaseConnectionOpts struct {
	Host               string
	Database           string
	User               string
	Password           string
	Port               string
	Timeout            int
	MaxConnections     int
	MaxIdleConnections int
	Params             map[string]string
}

// NewDatabaseConnectionOpts constructor function for DatabaseConnectionOpts
func NewDatabaseConnectionOpts(host, database, user, password, port string, timeout, maxConn, maxIdleConn int) *DatabaseConnectionOpts {
	return &DatabaseConnectionOpts{
		Host:               host,
		Database:           database,
		User:               user,
		Password:           password,
		Port:               port,
		Timeout:            timeout,
		MaxConnections:     maxConn,
		MaxIdleConnections: maxIdleConn,
	}
}

// GetConnString returns a standard connection string
func (opts *DatabaseConnectionOpts) GetConnString() string {
	return fmt.Sprintf(
		"%s:%s@%s:%s/%s?%s",
		opts.User,
		opts.Password,
		opts.Host,
		opts.Port,
		opts.Database,
		opts.GetParams(),
	)
}

// GetParams returns all the parameters with their values
// to be used in a connection string
func (opts *DatabaseConnectionOpts) GetParams() string {
	if opts.Params == nil || len(opts.Params) == 0 {
		return ""
	}
	var buffer bytes.Buffer
	for k, v := range opts.Params {
		buffer.WriteString(k)
		buffer.WriteString("=")
		buffer.WriteString(v)
		buffer.WriteString("&")
	}
	str := buffer.String()
	return str[:len(str)-1]
}

// DatabaseChecker simple database checker for application
type DatabaseChecker struct {
	db *goqu.Database
}

// NewChecker constructor
func NewChecker(db *goqu.Database) *DatabaseChecker {
	return &DatabaseChecker{
		db: db,
	}
}

// Diagnose start diagnose check
// http://confluence.alaudatech.com/pages/viewpage.action?pageId=14123161
func (d *DatabaseChecker) Diagnose() diagnose.ComponentReport {
	var (
		err   error
		start time.Time
	)

	report := diagnose.NewReport("database")
	start = time.Now()

	// does not work
	// err = d.db.Db.Ping()

	// this works
	tx, err := d.db.Db.Begin()
	if tx != nil {
		tx.Rollback()
	}
	report.Check(err, "Database ping failed", "Check environment variables or database health")
	report.AddLatency(start)
	return *report
}
