package database

import (
	"time"

	"moocss.com/tiga/pkg/log"
	"moocss.com/tiga/pkg/log/stdlog"
)

// Option is database option
type Option func(*options)

type options struct {
	dialect         string        // mysql|sqlite3|postgres, default:"mysql"
	dsn             string        // "foo:bar@tcp(127.0.0.1:3306)/baz?charset=utf8&parseTime=True&loc=Local"
	maxOpenConns    int           // default: 100
	maxIdleConns    int           // default: 10
	connMaxLifetime time.Duration // default: 1h
	logging         bool          // default: "false"

	logger log.Logger
}

// DefaultOptions .
func DefaultOptions() *options {
	return &options{
		dialect:         "mysql",
		dsn:             "",
		maxOpenConns:    100,
		maxIdleConns:    10,
		connMaxLifetime: 10 * time.Minute,
		logging:         false,

		logger: stdlog.NewLogger(),
	}
}

// Dialect with database Dialect.
func Dialect(dialect string) Option {
	return func(o *options) {
		o.dialect = dialect
	}
}

// DSN with database DSN.
func DSN(dsn string) Option {
	return func(o *options) {
		o.dsn = dsn
	}
}

// MaxOpenConns with database MaxOpenConns.
func MaxOpenConns(moc int) Option {
	return func(o *options) {
		o.maxOpenConns = moc
	}
}

// Dsn with database MaxIdleConns.
func MaxIdleConns(mic int) Option {
	return func(o *options) {
		o.maxIdleConns = mic
	}
}

// Dsn with database ConnMaxLifetime.
func ConnMaxLifetime(cml time.Duration) Option {
	return func(o *options) {
		o.connMaxLifetime = cml
	}
}

// Dsn with database Logging.
func Logging(logging bool) Option {
	return func(o *options) {
		o.logging = logging
	}
}

// WithLogger with database loogger.
func WithLogger(l log.Logger) Option {
	return func(o *options) {
		o.logger = l
	}
}
