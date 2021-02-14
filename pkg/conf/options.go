package conf

import (
	"moocss.com/tiga/pkg/conf/file"
	"moocss.com/tiga/pkg/log"
	"moocss.com/tiga/pkg/log/stdlog"
)

// Option is config option
type Option func(*options)

type options struct {
	sources []*file.File
	logger  log.Logger
}

// DefaultOptions .
func DefaultOptions() *options {
	return &options{
		logger: stdlog.NewLogger(),
	}
}

// WithSource with config source.
func WithSource(s ...*file.File) Option {
	return func(o *options) {
		o.sources = s
	}
}

// WithLogger with config loogger.
func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}
