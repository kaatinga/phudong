package phudong

import (
	"context"
	"testing"
	"time"
)

func TestWithLogger(t *testing.T) {
	opts := newOptions(WithLogger(nil))
	if opts.logger == nil {
		t.Error("Logger should not be nil when nil is passed")
	}
}

func TestWithName(t *testing.T) {
	opts := newOptions(WithName(""))
	if opts.name != "noName worker" {
		t.Errorf("Expected default name, got %s", opts.name)
	}

	opts = newOptions(WithName("test-worker"))
	if opts.name != "test-worker" {
		t.Errorf("Expected 'test-worker', got %s", opts.name)
	}
}

func TestWithInstantRun(t *testing.T) {
	opts := newOptions(WithInstantRun(true))
	if !opts.instantRun {
		t.Error("Instant run should be true")
	}

	opts = newOptions(WithInstantRun(false))
	if opts.instantRun {
		t.Error("Instant run should be false")
	}
}

func TestWithDuration(t *testing.T) {
	opts := newOptions(WithDuration(5 * time.Second))
	if opts.duration != 5*time.Second {
		t.Errorf("Expected 5s duration, got %v", opts.duration)
	}

	opts = newOptions(WithDuration(-1 * time.Second))
	if opts.duration != time.Hour {
		t.Errorf("Expected 1h duration for negative value, got %v", opts.duration)
	}
}

func TestWithDoThis(t *testing.T) {
	var executed bool
	f := func(ctx context.Context) {
		executed = true
	}

	opts := newOptions(WithDoThis(f))
	if len(opts.doThis) != 1 {
		t.Errorf("Expected 1 function, got %d", len(opts.doThis))
	}

	opts.doThis[0](context.Background())
	if !executed {
		t.Error("Function was not executed")
	}

	// Test nil function
	opts = newOptions(WithDoThis(nil))
	if len(opts.doThis) != 0 {
		t.Error("Nil function should not be added")
	}
}

func TestWithDoThisOrThrowError(t *testing.T) {
	var executed bool
	f := func(ctx context.Context) error {
		executed = true
		return nil
	}

	opts := newOptions(WithDoThisOrThrowError(f))
	if len(opts.doThisOrThrowError) != 1 {
		t.Errorf("Expected 1 function, got %d", len(opts.doThisOrThrowError))
	}

	err := opts.doThisOrThrowError[0](context.Background())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !executed {
		t.Error("Function was not executed")
	}

	// Test nil function
	opts = newOptions(WithDoThisOrThrowError(nil))
	if len(opts.doThisOrThrowError) != 0 {
		t.Error("Nil function should not be added")
	}
}

func TestWithErrorProcessor(t *testing.T) {
	var processed bool
	f := func(ctx context.Context, err error) {
		processed = true
	}

	opts := newOptions(WithErrorProcessor(f))
	if opts.withErrorProcessor == nil {
		t.Error("Error processor should be set")
	}

	opts.withErrorProcessor(context.Background(), nil)
	if !processed {
		t.Error("Error processor was not called")
	}

	// Test nil function
	opts = newOptions(WithErrorProcessor(nil))
	if opts.withErrorProcessor != nil {
		t.Error("Nil error processor should not be set")
	}
}

func TestMultipleOptions(t *testing.T) {
	var executed bool
	f := func(ctx context.Context) {
		executed = true
	}

	opts := newOptions(
		WithName("multi-test"),
		WithDuration(2*time.Second),
		WithInstantRun(true),
		WithDoThis(f),
	)

	if opts.name != "multi-test" {
		t.Errorf("Expected 'multi-test', got %s", opts.name)
	}
	if opts.duration != 2*time.Second {
		t.Errorf("Expected 2s duration, got %v", opts.duration)
	}
	if !opts.instantRun {
		t.Error("Instant run should be true")
	}
	if len(opts.doThis) != 1 {
		t.Error("Expected 1 function")
	}

	// Test that the function works
	opts.doThis[0](context.Background())
	if !executed {
		t.Error("Function was not executed")
	}
}

func TestDefaultOptions(t *testing.T) {
	opts := newOptions()

	if opts.name != "noName worker" {
		t.Errorf("Expected default name, got %s", opts.name)
	}
	if opts.duration != time.Hour {
		t.Errorf("Expected 1h duration, got %v", opts.duration)
	}
	if opts.instantRun {
		t.Error("Instant run should be false by default")
	}
	if opts.logger == nil {
		t.Error("Logger should not be nil")
	}
}
