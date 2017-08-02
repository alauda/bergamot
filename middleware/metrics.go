package middleware

import (
	"fmt"
	"time"

	"github.com/alauda/bergamot/errors"
	"github.com/alauda/bergamot/metrics"
)

// based on: https://gokit.io/examples/stringsvc.html

// BaseMetrics middleware for logging
type BaseMetrics struct {
	// component name
	component string
	// sample rate
	rate float64
	// map of metrics formats
	metMap map[string]string
	// real metrics client
	metrics metrics.Client
}

// NewMetrics constructor for base metrics middleware
func NewMetrics(component string, sampleRate float64, client metrics.Client) BaseMetrics {
	return BaseMetrics{
		component: component,
		rate:      sampleRate,
		metrics:   client,
		metMap: map[string]string{
			"latency":  "comp." + component + ".requests.latency",
			"request":  "comp." + component + ".requests.%d",
			"mLatency": "comp." + component + ".requests.%s.latency",
			"mRequest": "comp." + component + ".requests.%s.%d",
		},
	}
}

func (mw BaseMetrics) getStatus(err error) int {
	if err == nil {
		// OK
		return 200
	}
	if al, ok := err.(*errors.AlaudaError); ok {
		return al.StatusCode
	}
	// Unknown issue
	return 500
}

func (mw BaseMetrics) getTag(key, value string) string {
	return fmt.Sprintf("%s:%s", key, value)
}

// GenerateMetrics API to automatically generate metrics for a given endpoint
func (mw BaseMetrics) GenerateMetrics(begin time.Time, module, method, action string, err error) {
	// converting difference to milliseconds
	milliseconds := time.Since(begin).Seconds() * 1e3

	// general latency
	mw.metrics.Gauge(
		mw.metMap["latency"],
		// value
		milliseconds,
		// tags
		[]string{mw.getTag("module", module)},
		mw.rate,
	)
	// general request
	mw.metrics.Count(
		fmt.Sprintf(mw.metMap["request"], mw.getStatus(err)),
		// value
		1,
		// tags
		[]string{mw.getTag("module", module)},
		mw.rate,
	)
	// spacific latency
	mw.metrics.Gauge(
		fmt.Sprintf(mw.metMap["mLatency"], module),
		// value
		milliseconds,
		// tags
		[]string{mw.getTag("action", action)},
		mw.rate,
	)
	// specific request
	mw.metrics.Count(
		fmt.Sprintf(mw.metMap["mRequest"], module, mw.getStatus(err)),
		// value
		1,
		// tags
		[]string{mw.getTag("action", action)},
		mw.rate,
	)
}
