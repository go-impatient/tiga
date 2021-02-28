package log

import (
	"bytes"
	"fmt"
	"io"
	"sync"

	"github.com/sirupsen/logrus"
)

// logrusLogger .
type logrusLogger struct {
	entry *logrus.Entry
	pool  *sync.Pool
	level logrus.Level
}

type Option func(*logrusLogger)

// NewLogrusLogger returns a log.Logger that sends log events to a logrus.Logger.
func NewLogrusLogger(w io.Writer, opts ...Option) Logger {
	logger := logrus.New()
	logger.SetOutput(w)
	logger.SetReportCaller(true)
	e := logrus.NewEntry(logger)

	l := &logrusLogger{
		entry: e,
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		level: logrus.InfoLevel,
	}

	for _, fn := range opts {
		fn(l)
	}

	return l
}

// WithLevel configures a logrus logger to log at level for all events.
func WithLevel(level logrus.Level) Option {
	return func(c *logrusLogger) {
		c.level = level
	}
}

func (l *logrusLogger) Print(pairs ...interface{}) {
	if len(pairs) == 0 {
		return
	}
	if len(pairs)%2 != 0 {
		pairs = append(pairs, "")
	}
	buf := l.pool.Get().(*bytes.Buffer)
	for i := 0; i < len(pairs); i += 2 {
		fmt.Fprintf(buf, "%s=%v ", pairs[i], pairs[i+1])
	}

	switch l.level {
	case logrus.InfoLevel:
		l.entry.Info(buf.String())
	case logrus.ErrorLevel:
		l.entry.Error(buf.String())
	case logrus.DebugLevel:
		l.entry.Debug(buf.String())
	case logrus.WarnLevel:
		l.entry.Warn(buf.String())
	default:
		l.entry.Print(buf.String())
	}

	buf.Reset()
	l.pool.Put(buf)
}
