package phudong

import (
	"context"
	"errors"
	"reflect"
	"sync"
	"testing"
	"time"
)

func TestNewWorker(t *testing.T) {
	worker := NewWorker()
	if worker == nil {
		t.Fatal("NewWorker() returned nil")
	}
}

type lifecycleLogger struct {
	formats []string
}

func (l *lifecycleLogger) Printf(format string, args ...any) {
	l.formats = append(l.formats, format)
}

func (l *lifecycleLogger) Errorf(format string, args ...any) {}

func TestWorkerLifecycleLogFormats(t *testing.T) {
	logger := &lifecycleLogger{}
	worker := NewWorker(
		WithLogger(logger),
		WithDuration(time.Hour),
	)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	worker.Start(ctx)
	worker.Wait()

	expected := []string{"noName worker started", "noName worker stopped"}
	if !reflect.DeepEqual(logger.formats, expected) {
		t.Errorf("Expected lifecycle log formats %q, got %q", expected, logger.formats)
	}
}

func TestWorkerWithOptions(t *testing.T) {
	var executed bool
	var mu sync.Mutex

	worker := NewWorker(
		WithName("test-worker"),
		WithDuration(100*time.Millisecond),
		WithInstantRun(true),
		WithDoThis(func(ctx context.Context) {
			mu.Lock()
			executed = true
			mu.Unlock()
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	if !executed {
		t.Error("Function was not executed")
	}
	mu.Unlock()
}

func TestWorkerInstantRun(t *testing.T) {
	var executed bool
	var mu sync.Mutex

	worker := NewWorker(
		WithInstantRun(true),
		WithDuration(time.Hour), // Long duration to ensure instant run is tested
		WithDoThis(func(ctx context.Context) {
			mu.Lock()
			executed = true
			mu.Unlock()
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	if !executed {
		t.Error("Instant run did not execute function")
	}
	mu.Unlock()
}

func TestWorkerNoInstantRun(t *testing.T) {
	var executed bool
	var mu sync.Mutex

	worker := NewWorker(
		WithInstantRun(false),
		WithDuration(time.Hour), // Long duration to ensure no instant execution
		WithDoThis(func(ctx context.Context) {
			mu.Lock()
			executed = true
			mu.Unlock()
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	if executed {
		t.Error("Function was executed despite instant run being false")
	}
	mu.Unlock()
}

func TestWorkerWithError(t *testing.T) {
	var errorProcessed bool
	var mu sync.Mutex

	worker := NewWorker(
		WithDuration(50*time.Millisecond),
		WithInstantRun(true),
		WithDoThisOrThrowError(func(ctx context.Context) error {
			return errors.New("test error")
		}),
		WithErrorProcessor(func(ctx context.Context, err error) {
			mu.Lock()
			errorProcessed = true
			mu.Unlock()
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	if !errorProcessed {
		t.Error("Error was not processed")
	}
	mu.Unlock()
}

func TestWorkerNoFunctionSet(t *testing.T) {
	var errorProcessed bool
	var mu sync.Mutex

	worker := NewWorker(
		WithDuration(50*time.Millisecond),
		WithInstantRun(true),
		WithErrorProcessor(func(ctx context.Context, err error) {
			mu.Lock()
			errorProcessed = true
			mu.Unlock()
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	if !errorProcessed {
		t.Error("No function set error was not processed")
	}
	mu.Unlock()
}

func TestWorkerContextCancellation(t *testing.T) {
	var executionCount int
	var mu sync.Mutex

	worker := NewWorker(
		WithDuration(50*time.Millisecond),
		WithDoThis(func(ctx context.Context) {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	if executionCount < 2 {
		t.Errorf("Expected at least 2 executions, got %d", executionCount)
	}
	mu.Unlock()
}

func TestWorkerMultipleFunctions(t *testing.T) {
	var executionCount int
	var mu sync.Mutex

	worker := NewWorker(
		WithDuration(100*time.Millisecond),
		WithInstantRun(true),
		WithDoThis(func(ctx context.Context) {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}),
		WithDoThis(func(ctx context.Context) {
			mu.Lock()
			executionCount++
			mu.Unlock()
		}),
		WithDoThisOrThrowError(func(ctx context.Context) error {
			mu.Lock()
			executionCount++
			mu.Unlock()
			return nil
		}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	if executionCount != 3 {
		t.Errorf("Expected 3 executions (instant run), got %d", executionCount)
	}
	mu.Unlock()
}

func TestNextDailyAtToday(t *testing.T) {
	now := time.Date(2026, 8, 27, 0, 30, 0, 0, time.UTC)
	got := nextDailyAt(now, 1, 0)
	want := time.Date(2026, 8, 27, 1, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestNextDailyAtTomorrowWhenAlreadyPassed(t *testing.T) {
	now := time.Date(2026, 8, 27, 1, 0, 0, 0, time.UTC)
	got := nextDailyAt(now, 1, 0)
	want := time.Date(2026, 8, 28, 1, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

func TestNextDailyAtConvertsToUTC(t *testing.T) {
	// 03:30 UTC+2 is 01:30 UTC, so 01:00 UTC has already passed.
	now := time.Date(2026, 8, 27, 3, 30, 0, 0, time.FixedZone("UTC+2", 2*60*60))
	got := nextDailyAt(now, 1, 0)
	want := time.Date(2026, 8, 28, 1, 0, 0, 0, time.UTC)
	if !got.Equal(want) {
		t.Errorf("Expected %v, got %v", want, got)
	}
}

type settableClock struct {
	mu sync.Mutex
	t  time.Time
}

func (c *settableClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.t
}

func (c *settableClock) Set(t time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.t = t
}

func TestWorkerDailyAtRunsOnceThenWaits(t *testing.T) {
	clock := &settableClock{t: time.Date(2026, 8, 27, 0, 59, 59, 900_000_000, time.UTC)}

	var mu sync.Mutex
	var n int
	worker := NewWorker(
		WithDailyAt(1, 0),
		WithDoThis(func(ctx context.Context) {
			mu.Lock()
			n++
			mu.Unlock()
			clock.Set(clock.Now().Add(time.Second))
		}),
	)
	worker.now = clock.Now

	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	defer mu.Unlock()
	if n != 1 {
		t.Errorf("Expected 1 run at 1:00 UTC, got %d", n)
	}
}

func TestWorkerDailyAtDoesNotRunOnStart(t *testing.T) {
	now := time.Date(2026, 8, 27, 0, 59, 59, 900_000_000, time.UTC)

	var executed bool
	worker := NewWorker(
		WithDailyAt(1, 0),
		WithDoThis(func(ctx context.Context) {
			executed = true
		}),
	)
	worker.now = func() time.Time { return now }

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	if executed {
		t.Error("Daily-at worker should not run before the UTC time")
	}
}

func TestWorkerDailyAtIgnoresDuration(t *testing.T) {
	now := time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)

	var executed bool
	worker := NewWorker(
		WithDuration(10*time.Millisecond),
		WithDailyAt(1, 0),
		WithDoThis(func(ctx context.Context) {
			executed = true
		}),
	)
	worker.now = func() time.Time { return now }

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	if executed {
		t.Error("Duration ticker should be ignored when DailyAt is set")
	}
}

func TestWorkerDailyAtWithInstantRun(t *testing.T) {
	now := time.Date(2026, 8, 27, 12, 0, 0, 0, time.UTC)

	var mu sync.Mutex
	var n int
	worker := NewWorker(
		WithInstantRun(true),
		WithDailyAt(1, 0),
		WithDoThis(func(ctx context.Context) {
			mu.Lock()
			n++
			mu.Unlock()
		}),
	)
	worker.now = func() time.Time { return now }

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	worker.Start(ctx)
	worker.Wait()

	mu.Lock()
	defer mu.Unlock()
	if n != 1 {
		t.Errorf("Expected only the instant run, got %d", n)
	}
}

func TestWorkerWaitGroup(t *testing.T) {
	worker := NewWorker(
		WithDuration(time.Hour), // Long duration
		WithDoThis(func(ctx context.Context) {}),
	)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	start := time.Now()
	worker.Start(ctx)
	worker.Wait()
	duration := time.Since(start)

	// Should wait for context cancellation
	if duration < 5*time.Millisecond {
		t.Error("Worker did not wait properly")
	}
}
