# ErrX
Error object for Go with private/public part, attributes and stack traces.

## Usage
```go
func databaseError() error {
    return errx.Err(sql.ErrNoRows, "database error").Attr("db", "mydb")
}

func requestErrorPrivateWithMessage() error {
    if err := databaseError(); err != nil {
        return errx.Err(err, "error processing request").Private("unable to save data")
    }

    return nil
}

// ....
func main() {
  errx.RecordStackTrace = true // enable stack tracking

  err := requestErrorPrivateWithMessage()
  
  commonStr := err.Error() // "unable to save data"
  fullErrStr := errx.FullError(err) // "(unable to save data) error processing request: database error: sql: no rows in result set"

  errx.StackTrace(err) // functionName (shortfile.go:123) ... lines
  
  errx.GetAttrs(err) // map[db:"mydb"]
}
```

## License
MIT
