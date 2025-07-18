package phudong

type Error byte

const (
	ErrNoFunctionSet Error = iota
)

func (e Error) Error() string {
	switch e {
	case ErrNoFunctionSet:
		return "no function set to execute"
	default:
		return "unknown error"
	}
}

func (e Error) Is(target error) bool {
	typedError, ok := target.(Error)
	if !ok {
		return false
	}

	return e == typedError
}
