// Package logger - custom zap.Logger.
package logger

import (
	"log"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// zap.Logger settings.
var (
	Logger       *zap.Logger
	logInitError error
)

func init() {
	Logger, logInitError = New("info")
	if logInitError != nil {
		log.Fatal(logInitError)
	}
}

// New logger constructor.
func New(level string) (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	// set log minimum level.
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	// set app log level.
	atom := zap.NewAtomicLevel()
	err := atom.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}
	cfg.Level = atom

	// set output.
	cfg.OutputPaths = []string{"stdout"}
	cfg.DisableStacktrace = true

	// set log time mapping.
	cfg.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(time.StampNano))
	}

	return cfg.Build()
}
