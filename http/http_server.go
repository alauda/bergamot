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
	AddRoutes(router *iris.Router, server *Server)
}

// Middleware adding middleware
type Middleware interface {
	Serve(ctx *iris.Context)
}

// Config configuration for the HTTP server
type Config struct {
	Host               string
	Port               string
	AddLog             bool
	AddHealthCheck     bool
	TreatNotFoundError bool
	Component          string
	MaxReadBufferSize  int
	AllowedOrigins     []string
}

// SaneDefaults verifies the options and sets some sane defaults if
// the set values are not setup or not valid
func (c Config) SaneDefaults() Config {
	if c.MaxReadBufferSize <= 0 {
		c.MaxReadBufferSize = 1024 * 1024
	}
	if len(c.AllowedOrigins) == 0 {
		c.AllowedOrigins = []string{"*"}
	}
	return c
}

// GetIrisOptions returns the options for the iris http server
func (c Config) GetIrisOptions() []iris.OptionSetter {
	return []iris.OptionSetter{
		iris.OptionMaxHeaderBytes(c.MaxReadBufferSize),
	}
}

// GetCorsOptions return cors options for iris http router
func (c Config) GetCorsOptions() cors.Options {
	return cors.Options{AllowedOrigins: c.AllowedOrigins}
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
	config = config.SaneDefaults()
	return &Server{
		config:      config,
		log:         log,
		iris:        iris.New(config.GetIrisOptions()...),
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
		cors.New(h.config.GetCorsOptions()),
	)

	if h.config.AddHealthCheck {
		// adding health check
		h.iris.Any("/", h.Healthcheck)
		h.iris.Any("/_ping", h.Healthcheck)
	}

	if h.config.AddLog {
		// Adding request logger middleware
		h.iris.Use(h)
	}
	if h.config.TreatNotFoundError {
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
	handler.AddRoutes(router, h)

	return h
}

// AddVersionEndpoint add a root endpoint to a version specific API
// Used like AddEndpoint but will add on a specific version instead.
// If the version was not created previously will then be created automatically
func (h *Server) AddVersionEndpoint(version int, relativePath string, handler Router) *Server {
	h.AddVersion(version)
	handler.AddRoutes(h.versions[version].Party(relativePath), h)
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
	return GetMiddlewareHandlerFunc(h.GetMiddlewares(kind)...)
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

// GetMiddlewareHandlerFunc get only the functions from middlewares
func GetMiddlewareHandlerFunc(mws ...Middleware) []iris.HandlerFunc {
	var funcs []iris.HandlerFunc
	funcs = make([]iris.HandlerFunc, len(mws), len(mws)+1)
	for i, mw := range mws {
		funcs[i] = mw.Serve
	}
	return funcs
}

// DecorateHandlerFunc prepend all the given middlewares
func DecorateHandlerFunc(handlerFunc iris.HandlerFunc, mws ...Middleware) []iris.HandlerFunc {
	return append(GetMiddlewareHandlerFunc(mws...), handlerFunc)
}
