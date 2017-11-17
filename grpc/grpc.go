package grpc

import (
	"errors"
	"net"
	"runtime/debug"
	"time"

	"github.com/alauda/bergamot/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"fmt"
	"net/http"

	"github.com/cockroachdb/cmux"
)

// Registration used to register a grpc Service
type Registration func(server *grpc.Server)

// Server is a multiplexed server that adds a default HTTP1.1 healthcheck
type Server struct {
	config     Config
	registrars []Registration
	log        log.StandardLogger
}

// Config configuration for GRPC server
type Config struct {
	Port           string
	Component      string
	PeriodicMemory time.Duration
}

// New constructor function for the gRPC server
func New(config Config, log log.StandardLogger) *Server {
	return &Server{
		config:     config,
		log:        log,
		registrars: make([]Registration, 0, 1),
	}
}

// Add add a GRPC registration method
// GRPC register will be generated automatically using a proto file
// use this method to add different server registrars to a grpc.Server
func (g *Server) Add(registration Registration) {
	g.registrars = append(g.registrars, registration)
}

// Start will start serving on the GRPC server and block further execution
// should prefebly run inside a goroutine
func (g *Server) Start() error {
	if len(g.registrars) == 0 {
		return errors.New("No registration method added. impossible to boot")
	}
	listen, err := net.Listen("tcp", ":"+g.config.Port)
	if err != nil {
		g.log.Errorf("error listening to port", "port", g.config.Port, "err", err)
		return err
	}

	// The code below is to bootstrap a multiplexed server
	// this is necessary to create healthcheck endpoint from regular LoadBalancers
	// as they generate only use HTTP or TCP

	// creating multiplexed server
	mux := cmux.New(listen)

	// Matching connections by priority order
	grpcListener := mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
	// used for health checks
	httpListener := mux.Match(cmux.Any())

	// initiating grpc server
	grpcServer := grpc.NewServer()
	// registering handlers
	for _, r := range g.registrars {
		r(grpcServer)
	}
	reflection.Register(grpcServer)

	// creating http server
	httpServer := http.NewServeMux()
	httpServer.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})

	httpS := &http.Server{
		Handler: httpServer,
	}

	// starting it all
	go grpcServer.Serve(grpcListener)
	go httpS.Serve(httpListener)
	// will periodically free memory if set
	if g.config.PeriodicMemory > 0 {
		go g.PeriodicFree(g.config.PeriodicMemory)
	}

	g.log.Debugf("serving grpc", "port", g.config.Port)
	// Start serving...
	if err = mux.Serve(); err != nil {
		g.log.Errorf("error service grpc", "err", err)
	}
	return err
}

// PeriodicFree returns memory to OS given a span of time
func (g *Server) PeriodicFree(d time.Duration) {
	tick := time.Tick(d)
	for _ = range tick {
		debug.FreeOSMemory()
	}
}
