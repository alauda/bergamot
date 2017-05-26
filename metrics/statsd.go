package metrics

import (
	"fmt"
	"time"

	"github.com/DataDog/datadog-go/statsd"
)

// Client interface for metrics client
// currently using: https://github.com/DataDog/datadog-go/statsd
type Client interface {
	Gauge(name string, value float64, tags []string, rate float64) error
	Count(name string, value int64, tags []string, rate float64) error
	Histogram(name string, value float64, tags []string, rate float64) error
	Decr(name string, tags []string, rate float64) error
	Incr(name string, tags []string, rate float64) error
	Set(name string, value string, tags []string, rate float64) error
	Timing(name string, value time.Duration, tags []string, rate float64) error
	TimeInMilliseconds(name string, value float64, tags []string, rate float64) error
}

// New initiates a metrics client using statsd
func New(addr string, port int, open bool, bufferSize int) (Client, error) {
	if len(addr) > 0 && open {
		return statsd.NewBuffered(fmt.Sprintf("%s:%d", addr, port), bufferSize)
	}
	// if the address is not provided or it is closed will provide a mock
	return ClosedClient{}, nil
}

// ClosedClient an empty client for using a mocked simple client as a closed
type ClosedClient struct {
}

// Gauge does nothing
func (c ClosedClient) Gauge(name string, value float64, tags []string, rate float64) error {
	return nil
}

// Count does nothing
func (c ClosedClient) Count(name string, value int64, tags []string, rate float64) error {
	return nil
}

// Histogram does nothing
func (c ClosedClient) Histogram(name string, value float64, tags []string, rate float64) error {
	return nil
}

// Decr does nothing
func (c ClosedClient) Decr(name string, tags []string, rate float64) error {
	return nil
}

// Incr does nothing
func (c ClosedClient) Incr(name string, tags []string, rate float64) error {
	return nil
}

// Set does nothing
func (c ClosedClient) Set(name string, value string, tags []string, rate float64) error {
	return nil
}

// Timing does nothing
func (c ClosedClient) Timing(name string, value time.Duration, tags []string, rate float64) error {
	return nil
}

// TimeInMilliseconds does nothing
func (c ClosedClient) TimeInMilliseconds(name string, value float64, tags []string, rate float64) error {
	return nil
}
