package phudong

import (
	"context"
	"fmt"
	"os"
	"time"
)

type Logger interface {
	Printf(format string, args ...any)
	Errorf(format string, args ...any)
}

type stdLogger struct{}

func NewStdLogger() *stdLogger {
	return &stdLogger{}
}

func (l *stdLogger) Printf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stdout, format, args...)
}

func (l *stdLogger) Errorf(format string, args ...any) {
	_, _ = fmt.Fprintf(os.Stderr, "ERROR: "+format, args...)
}

type optionFunc func(*options)

type options struct {
	name       string
	instantRun bool
	duration   time.Duration

	doThis             []func(ctx context.Context)
	doThisOrThrowError []func(ctx context.Context) error

	withErrorProcessor func(ctx context.Context, err error)

	logger Logger
}

func newOptions(opts ...optionFunc) options {
	optsObj := options{
		name:     "noName worker",
		duration: time.Hour,
		logger:   NewStdLogger(),
	}

	for _, opt := range opts {
		opt(&optsObj)
	}

	return optsObj
}

// WithLogger sets the logger for the worker.
func WithLogger(logger Logger) optionFunc {
	return func(o *options) {
		if logger == nil {
			o.logger.Errorf("with logger: logger is nil\n")
			return
		}
		o.logger = logger
	}
}

func WithName(name string) optionFunc {
	return func(o *options) {
		if name == "" {
			o.logger.Errorf("with name: name is empty\n")
			return
		}

		o.name = name
	}
}

// WithInstantRun sets whether the worker should run immediately upon starting.
func WithInstantRun(enabled bool) optionFunc {
	return func(o *options) {
		o.instantRun = enabled
	}
}

// WithDoThis sets the function that the worker will execute periodically.
func WithDoThis(f func(ctx context.Context)) optionFunc {
	return func(o *options) {
		if f == nil {
			return
		}
		o.doThis = append(o.doThis, f)
	}
}

// WithDoThisOrThrowError sets a function that the worker will execute periodically.
func WithDoThisOrThrowError(f func(ctx context.Context) error) optionFunc {
	return func(o *options) {
		if f == nil {
			return
		}

		o.doThisOrThrowError = append(o.doThisOrThrowError, f)
	}
}

// WithDuration sets the duration for how often the worker should run.
func WithDuration(d time.Duration) optionFunc {
	return func(o *options) {
		o.duration = d
		if o.duration <= 0 {
			o.logger.Errorf("with duration: duration value is less than 0, setting to 1 hour\n")
			o.duration = time.Hour
		}
	}
}

// WithErrorProcessor sets a function that will be called to process errors.
func WithErrorProcessor(f func(ctx context.Context, err error)) optionFunc {
	return func(o *options) {
		if f == nil {
			o.logger.Errorf("with error processor: error processor is nil\n")
			return
		}
		o.withErrorProcessor = f
	}
}
