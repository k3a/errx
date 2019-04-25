[![GoDoc](https://godoc.org/github.com/k3a/errx?status.svg)](https://godoc.org/github.com/k3a/errx)
[![Build Status](https://travis-ci.org/k3a/errx.svg?branch=master)](https://travis-ci.org/k3a/errx)
[![Coverage Status](https://coveralls.io/repos/k3a/errx/badge.svg?branch=master&service=github)](https://coveralls.io/github/k3a/errx?branch=master)
[![Report Card](https://goreportcard.com/badge/github.com/k3a/errx)](https://goreportcard.com/report/github.com/k3a/errx)

# ErrX
Error object for Go with private/public part, attributes and stack traces.

It implements common `error` interface to be compatible with any other library
but additionally implements a concept of *extended private part*.
That part can be accessed after a successful type-casting of common `error` pointer
to this errx.Error or using a global functions like `errx.FullError()`.

Function `errx.FullError()` will attempt to do the typecast and return a full message
(with public and private parts) or, if called with a generic `error` type argument, 
it returns just the common `.Error()` result string.

## Usage

```go
// a private function deep in a proprietary library.
// It allocates a new error, describes it with a note and assigns an attribute to it.
func databaseError() error {
    return errx.Err(sql.ErrNoRows, "database error").Attr("db", "mydb")
}

// a function closer to the user, transforming to a public error with
// more generic public message "unable to save data".
//
// From now on, calling .Error() on the returned error returns just
// the public part "unable to save data".
//
// But smarter code can type-cast it back to errx.Error to get access
// to the private parts, or use top-level helpers like errx.FullError().
func requestErrorPrivateWithMessage() error {
    if err := databaseError(); err != nil {
        return errx.Err(err, "error processing request").Public("unable to save data")
    }

    return nil
}

// ....
func main() {
  errx.RecordStackTrace = true // optionally enable stack tracking

  err := requestErrorPrivateWithMessage()
  
  // public code using .Error() gets public part only:
  
  commonErrStr := err.Error() // just "unable to save data"


  // errx-aware code can request complete error state:

  fullErrStr := errx.FullError(err) // "unable to save data: error processing request: database error: sql: no rows in result set"

  errx.StackTrace(err) // functionName (shortfile.go:123) ... lines
  
  errx.GetAttrs(err) // map[db:"mydb"]
}
```

## License

MIT
