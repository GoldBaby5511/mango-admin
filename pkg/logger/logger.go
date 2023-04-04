package logger

import (
	"mango-admin/pkg/debug/writer"

	"io"
	"mango-admin/pkg"

	"os"
)

var (
	// DefaultLogger logger
	DefaultLogger Logger
)

// Logger is a generic logging interface
type Logger interface {
	// Init initialises options
	Init(options ...Option) error
	// Options The Logger options
	Options() Options
	// Fields set fields to always be logged
	Fields(fields map[string]interface{}) Logger
	// Log writes a log entry
	Log(level Level, v ...interface{})
	// Logf writes a formatted log entry
	Logf(level Level, format string, v ...interface{})
	// String returns the name of logger
	String() string
}

// SetupLogger 日志 cap 单位为kb
func SetupLogger(opts ...Option) Logger {
	op := setDefault()
	for _, o := range opts {
		o(&op)
	}
	if !pkg.PathExist(op.path) {
		err := pkg.PathCreate(op.path)
		if err != nil {
			Fatalf("create dir error: %s", err.Error())
		}
	}
	var err error
	var output io.Writer
	switch op.stdout {
	case "file":
		output, err = writer.NewFileWriter(
			writer.WithPath(op.path),
			writer.WithCap(op.cap<<10),
		)
		if err != nil {
			Fatal("logger setup error: %s", err.Error())
		}
	default:
		output = os.Stdout
	}
	var level Level
	level, err = GetLevel(op.level)
	if err != nil {
		Fatalf("get logger level error, %s", err.Error())
	}

	switch op.driver {
	case "zap":
		DefaultLogger, err = NewLoggerZap(WithLevel(level), WithOutput(output), WithCallerSkip(2))
		if err != nil {
			Fatalf("new zap logger error, %s", err.Error())
		}
	//case "logrus":
	//	setLogger = logrus.NewLogger(logger.WithLevel(level), logger.WithOutput(output), logrus.ReportCaller())
	default:
		DefaultLogger = NewLogger(WithLevel(level), WithOutput(output))
	}
	return DefaultLogger
}

func Init(opts ...Option) error {
	return DefaultLogger.Init(opts...)
}

func Fields(fields map[string]interface{}) Logger {
	return DefaultLogger.Fields(fields)
}

func Log(level Level, v ...interface{}) {
	DefaultLogger.Log(level, v...)
}

func Logf(level Level, format string, v ...interface{}) {
	DefaultLogger.Logf(level, format, v...)
}

func String() string {
	return DefaultLogger.String()
}
