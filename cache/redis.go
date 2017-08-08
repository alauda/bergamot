package cache

import (
	"fmt"
	"time"

	"github.com/alauda/bergamot/diagnose"

	aredis "github.com/alauda/go-redis-client"
)

type (
	// RedisCache abstraction used to store a redis connection
	RedisCache struct {
		Read  *aredis.RedisClient
		Write *aredis.RedisClient
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
	readOpts := aredis.RedisClientOptions{
		Type:     aredis.ClientNormal,
		Hosts:    []string{opts.GetAddr()},
		Password: opts.Password,
		Database: opts.DB,
	}
	writerOpts := aredis.RedisClientOptions{
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
func NewAlaudaRedis(opts aredis.RedisClientOptions, writerOpts aredis.RedisClientOptions) (*RedisCache, error) {
	reader := aredis.NewRedisClient(opts)
	writer := reader
	if len(writerOpts.Hosts) > 0 {
		writer = aredis.NewRedisClient(writerOpts)
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
