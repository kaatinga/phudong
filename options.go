package phudong

import (
	"context"
	"time"

	"github.com/kaatinga/qlog"
)

type optionFunc func(*options)

type options struct {
	// Add any fields that you want to configure for the worker
	name               string
	instantRun         bool
	duration           time.Duration
	doThis             func(ctx context.Context)
	doThisOrThrowError func(ctx context.Context) error

	withErrorProcessor func(ctx context.Context, err error)

	logger qlog.Logger
}

func newOptions(opts ...optionFunc) options {
	optsObj := options{
		name:     "noName worker",
		duration: time.Hour,
		logger:   qlog.New(),
	}

	for _, opt := range opts {
		opt(&optsObj)
	}

	return optsObj
}

// WithLogger sets the logger for the worker.
func WithLogger(logger qlog.Logger) optionFunc {
	return func(o *options) {
		if logger == nil {
			return
		}
		o.logger = logger
	}
}

func WithName(name string) optionFunc {
	return func(o *options) {
		if name == "" {
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
		o.doThis = f
	}
}

// WithDoThisOrThrowError sets a function that the worker will execute periodically.
func WithDoThisOrThrowError(f func(ctx context.Context) error) optionFunc {
	return func(o *options) {
		if f == nil {
			return
		}

		o.doThisOrThrowError = f
	}
}

// WithDuration sets the duration for how often the worker should run.
func WithDuration(d time.Duration) optionFunc {
	return func(o *options) {
		o.duration = d
		if o.duration <= 0 {
			o.duration = time.Hour // Default to 1 hour if not set or invalid
		}
	}
}

// WithErrorProcessor sets a function that will be called to process errors.
func WithErrorProcessor(f func(ctx context.Context, err error)) optionFunc {
	return func(o *options) {
		if f == nil {
			return
		}
		o.withErrorProcessor = f
	}
}
