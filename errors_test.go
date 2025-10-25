package phudong

import (
	"errors"
	"testing"
)

func TestErrorString(t *testing.T) {
	if ErrNoFunctionSet.Error() != "no function set to execute" {
		t.Errorf("Expected 'no function set to execute', got %s", ErrNoFunctionSet.Error())
	}

	var unknownError Error = 255
	if unknownError.Error() != "unknown error" {
		t.Errorf("Expected 'unknown error', got %s", unknownError.Error())
	}
}

func TestErrorIs(t *testing.T) {
	if !ErrNoFunctionSet.Is(ErrNoFunctionSet) {
		t.Error("ErrNoFunctionSet should be equal to itself")
	}

	if ErrNoFunctionSet.Is(errors.New("different error")) {
		t.Error("ErrNoFunctionSet should not be equal to different error")
	}

	var unknownError Error = 255
	if unknownError.Is(ErrNoFunctionSet) {
		t.Error("Unknown error should not be equal to ErrNoFunctionSet")
	}
}

func TestErrorType(t *testing.T) {
	var err error = ErrNoFunctionSet
	if !errors.Is(err, ErrNoFunctionSet) {
		t.Error("Error should be identifiable with errors.Is")
	}
}
