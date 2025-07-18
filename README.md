# phudong

A simple, idiomatic Go worker helper to get rid of boilerplate code for periodic or background jobs. 

## Features
- Minimal, flexible API
- Functional options for configuration
- Context-aware worker lifecycle
- Pluggable logger (qlog.Logger interface)
- Custom error handling
- Easy to test and extend

## Name origin

"phudong" (Phù Đổng) — Named after Thánh Gióng, a mythical Vietnamese hero who helped his country.

## Installation

```
go get github.com/kaatinga/phudong
```

## Usage

```go
package main

import (
	"context"
	"time"

	"github.com/kaatinga/phudong"
	"github.com/kaatinga/qlog"
)

func main() {
	logger := qlog.New()
	worker := phudong.NewWorker(
		phudong.WithDuration(2*time.Second),
		phudong.WithInstantRun(true),
		phudong.WithLogger(logger),
		phudong.WithDoThis(func(ctx context.Context) {
			logger.Printf("tick at %v", time.Now())
		}),
		phudong.WithErrorProcessor(func(ctx context.Context, err error) {
			logger.Errorf("worker error: %v", err)
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()
}
```

## Configuration Options
- `WithName(name string)` — Set the worker's name
- `WithDuration(d time.Duration)` — Set how often the worker runs (default: 1 hour)
- `WithInstantRun(enabled bool)` — Run immediately on start
- `WithDoThis(func(ctx context.Context))` — Set the function to execute periodically
- `WithDoThisOrThrowError(func(ctx context.Context) error)` — Set a function that may return an error
- `WithErrorProcessor(func(ctx context.Context, err error))` — Set a function to process errors
- `WithLogger(qlog.Logger)` — Inject a custom logger (see below)

## Logger Interface
phudong uses a minimal logger interface compatible with [qlog](https://github.com/kaatinga/qlog):

```go
type Logger interface {
    Printf(format string, args ...any)
	Errorf(format string, args ...any)
}
```

You can use your own logger by implementing this interface and passing it via `WithLogger`.

## License

This project is licensed under the [MIT License](LICENSE).
