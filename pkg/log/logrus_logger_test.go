package logrus

import (
	"os"
	"testing"

	"moocss.com/tiga/pkg/log"
)

func TestLogrusLogger(t *testing.T) {
	logger := NewLogrusLogger(os.Stdout, WithLevel(4))

	log.Debug(logger).Print("log", "test debug")
	log.Info(logger).Print("log", "test info")
	log.Warn(logger).Print("log", "test warn")
	log.Error(logger).Print("log", "test error")
}
