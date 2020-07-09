package errors

// Compare with github.com/hashicorp/go-multierror
//
// go-multierror's implementation of `Is(err) bool` only considers the first item of error list
//   func (e chain) Is(target error) bool {
//     return errors.Is(e[0], target)
//   }
// while the implementation of "go.uber.org/multierr" will go through the error list, try every item until no match found
import (
	"go.uber.org/multierr"
)

// Append merges multi errors into one
func Append(lsh error, rhs ...error) error {
	if e := Unwrap(lsh); e != nil {
		lsh = e
	}
	for _, e := range rhs {
		lsh = multierr.Append(lsh, e)
	}
	return Wrap(lsh)
}

// All retrieve error array if err is a multi-error
// All will NOT gather errors recursively
func All(err error) []error {
	var all interface{ Errors() []error }
	if ok := As(err, &all); ok {
		return all.Errors()
	}
	return nil
}
