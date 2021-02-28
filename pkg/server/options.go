package server

import (
	"net/http"
	"time"

	"moocss.com/tiga/pkg/log"
	"moocss.com/tiga/pkg/server/middleware"
)

// Option is http server option.
type Option func(o *options)

// options is an application options.
type options struct {
	// run mode 可选 dev/prod/test
	mode string

	// TCP address to listen on, ":http" if empty
	network string
	address string

	// app 超时
	timeout time.Duration

	// https://medium.com/@simonfrey/go-as-in-golang-standard-net-http-config-will-break-your-production-environment-1360871cb72b
	readHeaderTimeout time.Duration
	readTimeout       time.Duration
	writeTimeout      time.Duration
	idleTimeout       time.Duration

	middleware middleware.Middleware
	logger     log.Logger

	// 注入第三方 http server
	handler http.Handler
}

// DefaultOptions .
func DefaultOptions() *options {
	return &options{
		mode:              "dev",
		network:           "tcp",
		address:           ":8000",
		readHeaderTimeout: 20 * time.Second,
		readTimeout:       60 * time.Second,
		writeTimeout:      120 * time.Second,
		idleTimeout:       90 * time.Second,
		logger:            log.DefaultLogger,
	}
}

// Network with server network.
func Network(network string) Option {
	return func(o *options) {
		o.network = network
	}
}

// Address with server address.
func Address(addr string) Option {
	return func(o *options) {
		o.address = addr
	}
}

// Mode with server mode.
func Mode(a string) Option {
	return func(s *options) {
		s.mode = a
	}
}

// ReadHeaderTimeout with server readHeaderTimeout.
func ReadHeaderTimeout(a time.Duration) Option {
	return func(s *options) {
		s.readHeaderTimeout = a
	}
}

// ReadTimeout with server readTimeout.
func ReadTimeout(a time.Duration) Option {
	return func(s *options) {
		s.readTimeout = a
	}
}

// WriteTimeout with server writeTimeout.
func WriteTimeout(a time.Duration) Option {
	return func(s *options) {
		s.writeTimeout = a
	}
}

// IdleTimeout with server idleTimeout.
func IdleTimeout(a time.Duration) Option {
	return func(s *options) {
		s.idleTimeout = a
	}
}

// Logger with app logger.
func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// Middleware with server middleware option.
func Middleware(m middleware.Middleware) Option {
	return func(s *options) {
		s.middleware = m
	}
}

// WithApp with options
func HttpHandler(h http.Handler) Option {
	return func(s *options) {
		s.handler = h
	}
}
