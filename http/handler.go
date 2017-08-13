package http

import (
	"context"

	"github.com/alauda/bergamot/errors"
	"github.com/alauda/bergamot/log"

	iris "gopkg.in/kataras/iris.v6"
)

const (
	// USER constant key for user in iris.Context
	USER = "USER"
)

// ContextUserKey user key
type ContextUserKey struct{}

// ContextParamsKey params key
type ContextParamsKey struct{}

// ContextPathKey path key
type ContextPathKey struct{}

var (
	// UserKey static key for UserObject in context
	UserKey ContextUserKey
	// ParamsKey static key for Query parameters in request
	ParamsKey ContextParamsKey
	// PathKey path string key for context
	PathKey ContextPathKey
)

// Handler base handler with common shared methods
// can be composed to provide default behaviour
type Handler struct {
}

// GetContext get context with predefined keys
func (Handler) GetContext(ctx *iris.Context) context.Context {
	// string
	c := context.WithValue(nil, PathKey, ctx.Path())
	// map[string]string
	c = context.WithValue(c, ParamsKey, ctx.URLParams())
	// user
	c = context.WithValue(c, UserKey, ctx.Get(USER))
	return c
}

// HandleError Function to handle errors and return a message
func (Handler) HandleError(err error, ctx *iris.Context, log log.Logger) {
	status := getErrorStatusCode(err)
	log.Debugf("Error: %v - returning status: %d", err, status)
	ctx.JSON(
		status,
		NewAlaudaError(err),
	)
}

// HandleErrors Function to handle errors and return a message
func (Handler) HandleErrors(errs []error, ctx *iris.Context, log log.Logger) bool {
	status := getErrorsStatusCode(errs)
	if status == 0 {
		return false
	}
	log.Debugf("Error: %v - returning status: %d", errs, status)
	ctx.JSON(
		status,
		NewAlaudaError(errs...),
	)
	return true
}

func getErrorsStatusCode(errs []error) int {
	if errs == nil || len(errs) == 0 {
		return iris.StatusInternalServerError
	}
	for _, e := range errs {
		if e != nil {
			return getErrorStatusCode(e)
		}
	}
	return iris.StatusInternalServerError
}

func getErrorStatusCode(err error) int {
	if alaudaErr, ok := err.(*errors.AlaudaError); ok {
		if alaudaErr.StatusCode != 0 {
			return alaudaErr.StatusCode
		}
	}
	return iris.StatusInternalServerError
}

// AlaudaError structure to represent alauda's standard error format
type AlaudaError struct {
	Errors []error `json:"errors"`
}

// NewAlaudaError constructor function for the errors
func NewAlaudaError(err ...error) *AlaudaError {
	return &AlaudaError{
		Errors: err,
	}
}

// some issue with this function for now
// func evalError(source string, err error) error {
// 	by, _ := json.Marshal(err)
// 	str := string(by)
// 	if str == "" || str == "{}" {
// 		return errors.NewCommon(source, err)
// 	}
// 	return err
// }
