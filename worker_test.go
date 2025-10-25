package phudong

import (
	"context"
	"errors"
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
