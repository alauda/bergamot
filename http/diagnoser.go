package http

import (
	"github.com/alauda/bergamot/diagnose"
	iris "gopkg.in/kataras/iris.v6"
)

// DiagnoseRouter struct to add a route
type DiagnoseRouter struct {
	*diagnose.HealthChecker
}

// NewDiagnoser constructor for a diagnose checker
func NewDiagnoser(checker *diagnose.HealthChecker) *DiagnoseRouter {
	return &DiagnoseRouter{
		HealthChecker: checker,
	}
}

// AddRoutes will add a route for diagnose endpoint
func (h *DiagnoseRouter) AddRoutes(router *iris.Router, server *Server) {
	router.Any("", func(ctx *iris.Context) {
		ctx.JSON(iris.StatusOK, h.Check())
	})
}
