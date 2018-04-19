package cache

import (
	"fmt"
	"time"

	"github.com/alauda/bergamot/diagnose"

	aredis "github.com/alauda/go-redis-client"
)

// NilReply represent Redis nil reply, .e.g. when key does not exist.
const NilReply = aredis.RedisNil

type (
	// RedisCache abstraction used to store a redis connection
	RedisCache struct {
		Read  *aredis.Client
		Write *aredis.Client
	}

	// Cache interface to describe a Cache client
	Cache interface {
		Reader() aredis.Commander
		Writer() aredis.Commander
	}

	// RedisOpts abstraction for connection settings
	RedisOpts struct {
		Host     string
		Port     int
		DB       int
		Password string
	}
)

// NewRedis constructor method for a new Redis cache
func NewRedis(opts RedisOpts, writeOpts RedisOpts) (*RedisCache, error) {
	readOpts := aredis.Options{
		Type:     aredis.ClientNormal,
		Hosts:    []string{opts.GetAddr()},
		Password: opts.Password,
		Database: opts.DB,
	}
	writerOpts := aredis.Options{
		Type:     aredis.ClientNormal,
		Password: writeOpts.Password,
		Database: writeOpts.DB,
	}
	if len(writeOpts.Host) > 0 {
		readOpts.ReadOnly = true
		writerOpts.Hosts = []string{writeOpts.GetAddr()}
	}
	return NewAlaudaRedis(readOpts, writerOpts)
}

// NewAlaudaRedis construtor based on alauda redis client
func NewAlaudaRedis(opts aredis.Options, writerOpts aredis.Options) (*RedisCache, error) {
	reader := aredis.NewClient(opts)
	writer := reader
	if len(writerOpts.Hosts) > 0 {
		writer = aredis.NewClient(writerOpts)
	}
	client := &RedisCache{
		Read:  reader,
		Write: writer,
	}
	if err := client.Read.Ping().Err(); err != nil {
		return nil, err
	}
	return client, nil
}

// Reader returns a read commander
func (r *RedisCache) Reader() aredis.Commander {
	return r.Read
}

// Writer returns a write commander
func (r *RedisCache) Writer() aredis.Commander {
	return r.Write
}

// Diagnose start diagnose check
// http://confluence.alaudatech.com/pages/viewpage.action?pageId=14123161
func (r *RedisCache) Diagnose() diagnose.ComponentReport {
	var (
		err   error
		start time.Time
	)

	report := diagnose.NewReport("redis")
	start = time.Now()
	err = r.Read.Ping().Err()
	report.AddLatency(start)
	report.Check(err, "Redis reader ping failed", "Check environment variables or redis health")
	start = time.Now()
	err = r.Write.Ping().Err()
	report.AddLatency(start)
	report.Check(err, "Redis writer ping failed", "Check environment variables or redis health")
	return *report
}

// GetAddr will return a string of the addres which is host:port
func (r RedisOpts) GetAddr() string {
	return fmt.Sprintf("%s:%d", r.Host, r.Port)
}

// IsCacheErr will ignore redis.nil for error handler
func IsCacheErr(err error) bool {
	return err != nil && err != NilReply
}
