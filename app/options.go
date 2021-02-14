package app

import (
	"context"
	"moocss.com/tiga/pkg/server"
	"os"
	"syscall"

	"moocss.com/tiga/pkg/log"
	"moocss.com/tiga/pkg/log/stdlog"
)

// Option is an application option.
type Option func(o *options)

// options is an application options.
type options struct {
	version string

	ctx  context.Context
	sigs []os.Signal

	logger log.Logger
	server *server.Server
}

// DefaultOptions .
func DefaultOptions() *options {
	return &options{
		logger: stdlog.NewLogger(),
		ctx:    context.Background(),
		sigs:   []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
	}
}

// Version with app version.
func Version(version string) Option {
	return func(o *options) {
		o.version = version
	}
}

// Context with app context.
func Context(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

// Signal with exit signals.
func Signal(sigs ...os.Signal) Option {
	return func(o *options) { o.sigs = sigs }
}

// Logger with app logger.
func Logger(logger log.Logger) Option {
	return func(o *options) { o.logger = logger }
}

// Server with server.
func Server(srv *server.Server) Option {
	return func(o *options) {
		o.server = srv
	}
}
