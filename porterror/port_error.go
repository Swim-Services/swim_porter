package porterror

import "runtime/debug"

type PortError struct {
	Err   string
	Trace string
}

func Wrap(err error) *PortError {
	return &PortError{Err: err.Error(), Trace: string(debug.Stack())}
}

func New(err string) *PortError {
	return &PortError{Err: err, Trace: string(debug.Stack())}
}

func (e *PortError) Error() string {
	return e.Err
}

func (e *PortError) StackTrace() string {
	return e.Trace
}
