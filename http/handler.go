package http

import (
	"context"

	"github.com/alauda/bergamot/contexts"
	"github.com/alauda/bergamot/errors"
	"github.com/alauda/bergamot/log"

	iris "gopkg.in/kataras/iris.v6"
)

// Handler base handler with common shared methods
// can be composed to provide default behaviour
type Handler struct {
}

const (
	// USER constant key for user in iris.Context
	USER = "USER"
)

// GetContext get context with predefined keys
func (Handler) GetContext(ctx *iris.Context, attach bool) context.Context {
	var parent context.Context
	if attach {
		parent = ctx
	}
	// string
	c := contexts.SetPath(parent, ctx.Path())

	// map[string]string
	c = contexts.SetParams(c, ctx.URLParams())

	// user
	c = contexts.SetUser(c, ctx.Get(USER))

	// URL arguments
	arguments := map[string]string{}
	ctx.VisitValues(func(key string, value interface{}) {
		if val, ok := value.(string); ok {
			arguments[key] = val
		}
	})
	c = contexts.SetArgs(c, arguments)
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
