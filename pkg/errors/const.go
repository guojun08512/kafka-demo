package errors

import (
	stderr "errors"
)

var (
	errBadData          = stderr.New("unexpected data")
	errBadService       = stderr.New("bad service request")
	errConflict         = stderr.New("already exists")
	errClosed           = stderr.New("already closed")
	errForbidden        = stderr.New("forbidden")
	errInternalError    = stderr.New("internal error")
	errInvalidArguments = stderr.New("invalid arguments")
	errInvalidType      = stderr.New("invalid type")
	errMarshal          = stderr.New("failed to marshal")
	errNilObject        = stderr.New("accessing nil object")
	errNotFound         = stderr.New("not found")
	errNotSupport       = stderr.New("not support or not implemented")
	errTimeout          = stderr.New("timeout")
	errUnavailable      = stderr.New("unavailable")
)

// template of error
//func %s(message ...interface{}) error {
//	return Wrap(err%s, message...)
//}
//
//func Is%s(err error) bool {
//	return Is(err, err%s)
//}

func BadService(message ...interface{}) error {
	return Wrap(errBadService, message...)
}

func IsBadService(err error) bool {
	return Is(err, errBadService)
}

func Forbidden(message ...interface{}) error {
	return Wrap(errForbidden, message...)
}

func IsForbidden(err error) bool {
	return Is(err, errForbidden)
}

// BadData is designed to handle corrupted/malformed data that generated/retrieved during data processing
// if the input arguments are in bad shape, or don't fulfill the requirement, InvalidArg should be used
func BadData(message ...interface{}) error {
	return Wrap(errBadData, message...)
}

func IsBadData(err error) bool {
	return Is(err, errBadData)
}

func NotSupport(message ...interface{}) error {
	return Wrap(errNotSupport, message...)
}

func IsNotSupport(err error) bool {
	return Is(err, errNotSupport)
}

func Timeout(message ...interface{}) error {
	return Wrap(errTimeout, message...)
}

func IsTimeout(err error) bool {
	return Is(err, errTimeout)
}

func Unavailable(message ...interface{}) error {
	return Wrap(errUnavailable, message...)
}

func IsUnavailable(err error) bool {
	return Is(err, errUnavailable)
}

func NilObject(message ...interface{}) error {
	return Wrap(errNilObject, message...)
}

func IsNilObject(err error) bool {
	return Is(err, errNilObject)
}

func InvalidArg(message ...interface{}) error {
	return Wrap(errInvalidArguments, message...)
}

func IsInvalidArg(err error) bool {
	return Is(err, errInvalidArguments)
}

func InternalError(message ...interface{}) error {
	return Wrap(errInternalError, message...)
}

func InternalErrorf(format string, args ...interface{}) error {
	return Wrapf(errInternalError, format, args...)
}

func IsInternalError(err error) bool {
	return Is(err, errInternalError)
}

func Conflict(message ...interface{}) error {
	return Wrap(errConflict, message...)
}

func IsConflict(err error) bool {
	return Is(err, errConflict)
}

func NotFound(message ...interface{}) error {
	return Wrap(errNotFound, message...)
}

func IsNotFound(err error) bool {
	return Is(err, errNotFound)
}

func Closed(message ...interface{}) error {
	return Wrap(errClosed, message...)
}

func IsClosed(err error) bool {
	return Is(err, errClosed)
}

func Marshal(message ...interface{}) error {
	return Wrap(errMarshal, message...)
}

func IsMarshal(err error) bool {
	return Is(err, errMarshal)
}

func Unmarshal(message ...interface{}) error {
	return WithMessage(InvalidArg(message...), "failed to unmarshal")
}

func InvalidType(message ...interface{}) error {
	return Wrap(errInvalidType, message...)
}

func IsInvalidType(err error) bool {
	return Is(err, errInvalidType)
}
