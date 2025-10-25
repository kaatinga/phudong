[![Go CI](https://github.com/kaatinga/phudong/actions/workflows/golang_ci.yml/badge.svg)](https://github.com/kaatinga/phudong/actions/workflows/golang_ci.yml)
[![CodeQL](https://github.com/kaatinga/phudong/actions/workflows/github-code-scanning/codeql/badge.svg)](https://github.com/kaatinga/phudong/actions/workflows/github-code-scanning/codeql)

# phudong

A simple, idiomatic Go worker helper to get rid of boilerplate code for periodic or background jobs.

[![Go Reference](https://pkg.go.dev/badge/github.com/kaatinga/phudong.svg)](https://pkg.go.dev/github.com/kaatinga/phudong)
[![Go Report Card](https://goreportcard.com/badge/github.com/kaatinga/phudong)](https://goreportcard.com/report/github.com/kaatinga/phudong)

## Name origin

"phudong" (Phù Đổng) — Named after Thánh Gióng, a mythical Vietnamese hero who helped his country.

## Installation

```bash
go get github.com/kaatinga/phudong
```

## Quick Start

```go
package main

import (
	"context"
	"time"

	"github.com/kaatinga/phudong"
)

func main() {
	worker := phudong.NewWorker(
		phudong.WithName("my-worker"),
		phudong.WithDuration(2*time.Second),
		phudong.WithInstantRun(true),
		phudong.WithDoThis(func(ctx context.Context) {
			// Your periodic task here
			println("Working...")
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()
}
```

## Advanced Usage

### Error Handling

```go
worker := phudong.NewWorker(
	phudong.WithDoThisOrThrowError(func(ctx context.Context) error {
		// This function can return an error
		if someCondition {
			return errors.New("something went wrong")
		}
		return nil
	}),
	phudong.WithErrorProcessor(func(ctx context.Context, err error) {
		// Handle errors here
		log.Printf("Worker error: %v", err)
	}),
)
```

### Custom Logger

```go
type MyLogger struct{}

func (l *MyLogger) Printf(format string, args ...interface{}) {
	log.Printf("[INFO] "+format, args...)
}

func (l *MyLogger) Errorf(format string, args ...interface{}) {
	log.Printf("[ERROR] "+format, args...)
}

worker := phudong.NewWorker(
	phudong.WithLogger(&MyLogger{}),
	// ... other options
)
```

### Multiple Functions

```go
worker := phudong.NewWorker(
	phudong.WithDoThis(func(ctx context.Context) {
		// First function
		updateMetrics()
	}),
	phudong.WithDoThis(func(ctx context.Context) {
		// Second function
		cleanupTempFiles()
	}),
	phudong.WithDoThisOrThrowError(func(ctx context.Context) error {
		// Function that can return errors
		return syncWithDatabase()
	}),
)
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithName(name string)` | Set the worker's name for logging | `"noName worker"` |
| `WithDuration(d time.Duration)` | Set execution interval | `1 hour` |
| `WithInstantRun(enabled bool)` | Run immediately on start | `false` |
| `WithDoThis(func(ctx context.Context))` | Add a function to execute | - |
| `WithDoThisOrThrowError(func(ctx context.Context) error)` | Add a function that may return an error | - |
| `WithErrorProcessor(func(ctx context.Context, err error))` | Set error handler | - |
| `WithLogger(Logger)` | Set custom logger | `NewStdLogger()` |

## Logger Interface

phudong uses a minimal, dependency-free logger interface:

```go
type Logger interface {
	Printf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}
```

### Built-in Logger

The package includes a default logger that writes to stdout for regular messages and stderr for errors:

```go
logger := phudong.NewStdLogger()
```

### Custom Logger Implementation

You can implement your own logger by implementing the `Logger` interface:

```go
type CustomLogger struct {
	level string
}

func (l *CustomLogger) Printf(format string, args ...interface{}) {
	fmt.Printf("[%s] %s\n", l.level, fmt.Sprintf(format, args...))
}

func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR] %s\n", fmt.Sprintf(format, args...))
}

// Use it
worker := phudong.NewWorker(
	phudong.WithLogger(&CustomLogger{level: "INFO"}),
)
```

## Real-world Examples

### Database Cleanup Worker

```go
func startCleanupWorker(db *sql.DB) {
	worker := phudong.NewWorker(
		phudong.WithName("db-cleanup"),
		phudong.WithDuration(24*time.Hour),
		phudong.WithInstantRun(true),
		phudong.WithDoThisOrThrowError(func(ctx context.Context) error {
			// Clean up old records
			_, err := db.ExecContext(ctx, "DELETE FROM logs WHERE created_at < ?", 
				time.Now().Add(-30*24*time.Hour))
			return err
		}),
		phudong.WithErrorProcessor(func(ctx context.Context, err error) {
			log.Printf("Database cleanup failed: %v", err)
		}),
	)

	ctx := context.Background()
	worker.Start(ctx)
	// Worker runs in background
}
```

### Health Check Worker

```go
func startHealthCheckWorker(services []string) {
	worker := phudong.NewWorker(
		phudong.WithName("health-check"),
		phudong.WithDuration(30*time.Second),
		phudong.WithDoThis(func(ctx context.Context) {
			for _, service := range services {
				if !checkServiceHealth(service) {
					alertServiceDown(service)
				}
			}
		}),
	)

	ctx := context.Background()
	worker.Start(ctx)
}
```

### Metrics Collection Worker

```go
func startMetricsWorker(metricsCollector *MetricsCollector) {
	worker := phudong.NewWorker(
		phudong.WithName("metrics-collector"),
		phudong.WithDuration(1*time.Minute),
		phudong.WithDoThis(func(ctx context.Context) {
			metrics := metricsCollector.Collect()
			sendToMonitoring(metrics)
		}),
		phudong.WithErrorProcessor(func(ctx context.Context, err error) {
			// Log metrics collection errors
			log.Printf("Metrics collection failed: %v", err)
		}),
	)

	ctx := context.Background()
	worker.Start(ctx)
}
```

## Testing

The package includes comprehensive tests. Run them with:

```bash
go test ./...
```

### Testing Your Workers

```go
func TestMyWorker(t *testing.T) {
	var executed bool
	worker := phudong.NewWorker(
		phudong.WithDuration(100*time.Millisecond),
		phudong.WithInstantRun(true),
		phudong.WithDoThis(func(ctx context.Context) {
			executed = true
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	if !executed {
		t.Error("Worker function was not executed")
	}
}
```

## Best Practices

1. **Always use context cancellation** - Set appropriate timeouts
2. **Handle errors gracefully** - Use `WithErrorProcessor` for error handling
3. **Use meaningful names** - Set worker names for better logging
4. **Test your workers** - Include tests for your worker functions
5. **Consider instant run** - Use `WithInstantRun(true)` for immediate execution

## License

This project is licensed under the [MIT License](LICENSE).
