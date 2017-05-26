package diagnose

import (
	"fmt"
	"time"

	"encoding/json"

	"sync"
)

// Component interface for a component health check
type Component interface {
	Diagnose() ComponentReport
}

// HealthChecker main health diagnoser
type HealthChecker struct {
	Components []Component
}

// New constructor function for HealthChecker
// http://confluence.alaudatech.com/pages/viewpage.action?pageId=14123161
func New() (*HealthChecker, error) {
	return &HealthChecker{}, nil
}

// Add add component for check
func (h *HealthChecker) Add(com Component) *HealthChecker {
	if h.Components == nil {
		h.Components = make([]Component, 0, 2)
	}
	h.Components = append(h.Components, com)
	return h
}

// Check starts health check of components
func (h *HealthChecker) Check() HealthReport {
	report := &HealthReport{Status: StatusOK}
	if h.Components == nil {
		return *report
	}
	wait := sync.WaitGroup{}
	for _, c := range h.Components {
		wait.Add(1)
		go func(component Component) {
			diagnose := component.Diagnose()
			if diagnose.Status == StatusError && report.Status != StatusError {
				report.Status = StatusError
			}
			report.Add(diagnose)
			wait.Done()
		}(c)
	}
	wait.Wait()
	return *report
}

// HealthStatus type to create health status
type HealthStatus string

const (
	// StatusOK means the component is K
	StatusOK HealthStatus = "OK"
	// StatusError means there is an error with the component
	StatusError HealthStatus = "ERROR"
)

// HealthReport report struct format
type HealthReport struct {
	Status  HealthStatus      `json:"status"`
	Details []ComponentReport `json:"details"`
}

// Add adds a new component report
func (h *HealthReport) Add(report ComponentReport) {
	if h.Details == nil {
		h.Details = make([]ComponentReport, 0, 3)
	}
	h.Details = append(h.Details, report)
}

// ComponentReport each component report
type ComponentReport struct {
	Status     HealthStatus  `json:"status"`
	Name       string        `json:"name"`
	Message    string        `json:"message"`
	Suggestion string        `json:"suggestion"`
	Latency    time.Duration `json:"latency"`
}

// NewReport constructor
func NewReport(component string) *ComponentReport {
	return &ComponentReport{
		Status:  StatusOK,
		Name:    component,
		Message: "ok",
	}
}

// Check check for error and add custom message
func (c *ComponentReport) Check(err error, message, suggestion string) {
	if err != nil {
		c.Status = StatusError
		c.Message = fmt.Sprintf("%s: \"%s\"", message, err.Error())
		c.Suggestion = suggestion
	}
}

// AddLatency add a latency for the start time
func (c *ComponentReport) AddLatency(start time.Time) {
	duration := time.Since(start)
	if c.Latency == time.Duration(0) || c.Latency < duration {
		c.Latency = duration
	}
}

// MarshalJSON custom json formater
func (c *ComponentReport) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Status     HealthStatus `json:"status"`
		Name       string       `json:"name"`
		Message    string       `json:"message"`
		Suggestion string       `json:"suggestion"`
		Latency    string       `json:"latency"`
	}{
		c.Status,
		c.Name,
		c.Message,
		c.Suggestion,
		c.Latency.String(),
	})
}
