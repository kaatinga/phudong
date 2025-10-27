package phudong

import (
	"bytes"
	"os"
	"testing"
)

func TestStdLogger(t *testing.T) {
	logger := NewStdLogger()
	if logger == nil {
		t.Fatal("NewStdLogger() returned nil")
	}
}

func TestStdLoggerErrorf(t *testing.T) {
	// Capture stderr
	oldStderr := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stderr = w

	logger := &stdLogger{}
	logger.Errorf("error message %d\n", 42)

	_ = w.Close()
	os.Stderr = oldStderr

	// Read the content
	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	content := buf.String()

	expected := "ERROR: error message 42\n"
	if content != expected {
		t.Errorf("Expected %q, got %q", expected, content)
	}
}

func TestLoggerInterface(t *testing.T) {
	// Test that stdLogger implements Logger interface
	var _ Logger = &stdLogger{}
}
