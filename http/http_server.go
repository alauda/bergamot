package http

import (
	"fmt"
	"time"

	"github.com/alauda/bergamot/log"

	"gopkg.in/kataras/iris.v6"
	"gopkg.in/kataras/iris.v6/adaptors/cors"
	"gopkg.in/kataras/iris.v6/adaptors/httprouter"
)

// Router interface for http endpoints
type Router interface {
	AddRoutes(router *iris.Router)
}

// Middleware adding middleware
type Middleware interface {
	Serve(ctx *iris.Context)
}

// Config configuration for the HTTP server
type Config struct {
	Host           string
	Port           string
	AddLog         bool
	AddHealthCheck bool
	Component      string
}

// Server Full HTTP server
type Server struct {
	config      Config
	start       time.Time
	log         log.Logger
	iris        *iris.Framework
	versions    map[int]*iris.Router
	middlewares map[string][]Middleware
}

// NewServer constructor function for the HTTP server
func NewServer(config Config, log log.Logger) *Server {
	return &Server{
		config:      config,
		log:         log,
		iris:        iris.New(),
		versions:    map[int]*iris.Router{},
		middlewares: map[string][]Middleware{MiddlewareTypeAll: []Middleware{}},
	}
}

// Init will setup any necessary data
func (h *Server) Init() *Server {
	h.iris.Adapt(
		// Logging all errors
		iris.DevLogger(),
		// adding router
		httprouter.New(),

		// Cors wrapper to the entire application, allow all origins.
		cors.New(cors.Options{AllowedOrigins: []string{"*"}}),
	)

	if h.config.AddHealthCheck {
		// adding health check
		h.iris.Any("/", h.Healthcheck)
		h.iris.Any("/_ping", h.Healthcheck)
	}

	if h.config.AddLog {
		// Adding request logger middleware
		h.iris.Use(h)
		// default error when requesting unexistent route
		h.iris.OnError(iris.StatusNotFound, func(ctx *iris.Context) {
			// print method and stuff
			h.Serve(ctx)
		})
	}

	return h
}

// AddVersion Adds a version number to the API route
func (h *Server) AddVersion(version int) *Server {
	if _, ok := h.versions[version]; !ok {
		// adds /v1 or /v2 route
		h.versions[version] = h.iris.Party(fmt.Sprintf("/v%d", version))
	}
	return h
}

// AddEndpoint ands a handler for the given relative path
// should be executed before the Start method and after the Init method
func (h *Server) AddEndpoint(relativePath string, handler Router) *Server {
	router := h.iris.Party(relativePath)
	handler.AddRoutes(router)

	return h
}

// AddVersionEndpoint add a root endpoint to a version specific API
// Used like AddEndpoint but will add on a specific version instead.
// If the version was not created previously will then be created automatically
func (h *Server) AddVersionEndpoint(version int, relativePath string, handler Router) *Server {
	h.AddVersion(version)
	handler.AddRoutes(h.versions[version].Party(relativePath))
	return h
}

// Serve will log all the requests
func (h *Server) Serve(ctx *iris.Context) {
	// logging all requests
	h.log.Infof("---- [%s] %s  - args: %s ", ctx.Method(), ctx.Path(), ctx.ParamsSentence())
	ctx.Next()
}

// Healthcheck healthcheck endpoint
func (h *Server) Healthcheck(ctx *iris.Context) {
	ctx.WriteString(fmt.Sprintf("%s:%s", h.config.Component, time.Since(h.start)))
}

// GetApp returns the iris app, used for testing
func (h *Server) GetApp() *iris.Framework {
	return h.iris
}

// Start will start serving the http server
// this method will block while serving http
func (h *Server) Start() {
	h.start = time.Now()
	h.iris.Listen(h.config.Host + ":" + h.config.Port)
}

const (
	// MiddlewareTypeAll special type
	MiddlewareTypeAll = "*"
)

// AddMiddleware adds a middleware for the given types
func (h *Server) AddMiddleware(mw Middleware, kinds ...string) *Server {
	var (
		collection []Middleware
		ok         bool
	)
	kinds = append(kinds, MiddlewareTypeAll)
	for _, k := range kinds {
		if collection, ok = h.middlewares[k]; !ok {
			collection = make([]Middleware, 0, 2)
		}
		collection = append(collection, mw)
		h.middlewares[k] = collection
	}
	return h
}

// AddMiddlewares adds a middleware for the given types
func (h *Server) AddMiddlewares(mws []Middleware, kinds ...string) *Server {
	for _, mw := range mws {
		h.AddMiddleware(mw, kinds...)
	}
	return h
}

// GetMiddlewares get all midlewares of a kind
func (h *Server) GetMiddlewares(kind string) []Middleware {
	return h.middlewares[kind]
}

// GetMiddlewareHandlerFun returns all the handler functions of a middleware kind
func (h *Server) GetMiddlewareHandlerFun(kind string) []iris.HandlerFunc {
	var funcs []iris.HandlerFunc
	mws := h.GetMiddlewares(kind)
	funcs = make([]iris.HandlerFunc, len(mws), len(mws)+1)
	for i, mw := range mws {
		funcs[i] = mw.Serve
	}
	return funcs
}

// GetMiddlewaresDecorated gets all the handler functions of a collection of kinds and decorate the target function
func (h *Server) GetMiddlewaresDecorated(handlerFunc iris.HandlerFunc, kinds ...string) []iris.HandlerFunc {
	var mws []iris.HandlerFunc
	for _, k := range kinds {
		mws = append(mws, h.GetMiddlewareHandlerFun(k)...)
	}
	mws = append(mws, handlerFunc)
	return mws
}
