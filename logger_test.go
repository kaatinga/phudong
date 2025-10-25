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

func TestStdLoggerPrintf(t *testing.T) {
	// Create a temporary file to capture output
	tmpfile, err := os.CreateTemp("", "test-logger")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		_ = os.Remove(tmpfile.Name())
	}()

	logger := &stdLogger{output: tmpfile}
	logger.Printf("test message %d\n", 42)

	// Read the content
	content, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := "test message 42\n"
	if string(content) != expected {
		t.Errorf("Expected %q, got %q", expected, string(content))
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
