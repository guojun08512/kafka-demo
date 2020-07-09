package errors

import (
	"fmt"

	"github.com/pkg/errors"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}

// Sprintlnn => Sprint no newline. This is to get the behavior of how
// sprintlnnln where spaces are always added between operands, regardless of
// their type. Instead of vendoring the Sprintln implementation to spare a
// string allocation, we do the simplest thing.
func sprintlnn(args ...interface{}) string {
	msg := fmt.Sprintln(args...)
	return msg[:len(msg)-1]
}

func Is(err, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target interface{}) bool {
	return errors.As(err, target)
}

func Unwrap(err error) error {
	return errors.Unwrap(err)
}

func Wrap(err error, message ...interface{}) error {
	if len(message) > 0 {
		return WithStack(WithMessage(err, sprintlnn(message...)))
	}
	return WithStack(err)
}

func Wrapf(err error, format string, args ...interface{}) error {
	return WithStack(WithMessagef(err, format, args...))
}

func WithMessage(err error, message string) error {
	return errors.WithMessage(err, message)
}

func WithMessages(err error, message ...interface{}) error {
	return errors.WithMessage(err, sprintlnn(message...))
}

func WithMessagef(err error, format string, args ...interface{}) error {
	return errors.WithMessagef(err, format, args...)
}

func New(message string) error {
	return errors.New(message)
}

func Errorf(format string, args ...interface{}) error {
	return errors.Errorf(format, args...)
}

// WithStack annotates err with a stack trace at the point WithStack was called.
// If err is nil, WithStack returns nil.
func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(),
	}
}
