package errx

import (
	"fmt"
	"path/filepath"
	"runtime"
)

const nilString = "(nil)"
const genericPublicMessage = "unexpected error"

// RecordStackTrace enables or disables stack tracking for errors (default disabled)
var RecordStackTrace = false

// Attributes is a map of key-value attributes
type Attributes map[string]interface{}

// Error contains an extended error info
type Error struct {
	// parent error (can be *Error, something implementing "error" interface or nil)
	parent error
	// is the error marked private?
	private bool

	// the actual error message of this error
	message string
	// used if public is true
	publicMessage string
	// program counter (0 if not known)
	pc uintptr

	// optional additional attributes
	attrs Attributes
}

func (e *Error) error(publicOnly bool) string {
	publicOnly = publicOnly || e.private

	msg := e.message
	if publicOnly {
		msg = e.publicMessage
	}

	// print parent errors
	err := e.parent
	if err != nil {
		if exErr, ok := err.(*Error); ok {
			parentMsg := exErr.error(publicOnly)

			// call ex Error directly
			if len(msg) > 0 && len(parentMsg) > 0 {
				msg += ": "
			}
			msg += parentMsg
		} else if !publicOnly {
			// if not private and not ex Error, print raw parent error
			if len(msg) > 0 {
				msg += ": "
			}
			msg += err.Error()
		}
	}

	return msg
}

// Error returns a public error
func (e *Error) Error() string {
	if e == nil {
		return nilString
	}

	return e.error(e.private)
}

// FullError returns a full error with additional information which shouldn't be returned to users.
// This is suitable for placing to a logfile.
func (e *Error) FullError() string {
	if e == nil {
		return nilString
	}

	msg := e.publicMessage

	if len(e.publicMessage) > 0 && len(e.message) > 0 {
		msg = "(" + e.publicMessage + ") "
	}
	msg += e.message

	// print parent errors
	err := e.parent
	if err != nil {
		if len(msg) > 0 {
			msg += ": "
		}

		if exErr, ok := err.(*Error); ok {
			// extra error
			msg += exErr.FullError()
		} else {
			// common error
			msg += err.Error()
		}
	}

	return msg
}

// StackTrace attempts to return a stack trace of this error and parents.
// Returns an empty string if unable to get the stack trace.
func (e *Error) StackTrace() string {
	if e.pc == 0 {
		return ""
	}

	str := ""

	fn := runtime.FuncForPC(e.pc)
	if fn == nil {
		str += " unknown:0 unknown\n"
	} else {
		file, line := fn.FileLine(e.pc)
		str += fmt.Sprintf(" %s (%s:%d)\n", filepath.Base(fn.Name()), filepath.Base(file), line)
	}

	// parents next
	err := e.parent
	if err != nil {
		if exErr, ok := err.(*Error); ok {
			// extra error
			str += exErr.StackTrace()
		}
	}

	return str
}

// Private marks this error and parents as private.
// Private errors return public messages only.
//
// Argument is an optional format string with formating arguments.
// Resulting sting is used as a public message.
//
// Example: Private("this is public user-facing msg: %s", additionalString)
func (e *Error) Private(fmtArgs ...interface{}) *Error {
	msg := genericPublicMessage

	for i, a := range fmtArgs {
		if str, ok := a.(string); ok {
			msg = fmt.Sprintf(str, fmtArgs[i+1:]...)
			break
		}
	}

	e.private = true
	e.publicMessage = msg

	return e
}

// Attr sets a value for a key attribute
func (e *Error) Attr(key string, val interface{}) *Error {
	if e.attrs == nil {
		e.attrs = make(Attributes)
	}
	e.attrs[key] = val

	return e
}

func (e *Error) getAttrs(inout Attributes) {
	for k, v := range e.attrs {
		inout[k] = v
	}

	if exErr, ok := e.parent.(*Error); ok {
		exErr.getAttrs(inout)
	}
}

// GetAttrs returns assigned attributes of the whole error chain or nil
func (e *Error) GetAttrs() Attributes {
	ret := make(Attributes)
	e.getAttrs(ret)
	return ret
}

// Err creates a new error.
// Argument can be a parent error or format string with formating arguments.
// Example: Err(err, "error reading file %s", path)
func Err(args ...interface{}) *Error {
	err := new(Error)

	if RecordStackTrace {
		err.pc, _, _, _ = runtime.Caller(1)
	}

	for i, a := range args {
		if e, ok := a.(error); ok {
			err.parent = e
		} else if str, ok := a.(string); ok {
			err.message = fmt.Sprintf(str, args[i+1:]...)
			break
		}
	}

	return err
}
