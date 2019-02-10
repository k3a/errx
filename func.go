package errx

// StackTrace attempts to return a stack trace for the error if possible or an empty string
func StackTrace(err error) string {
	if err != nil {
		if exErr, ok := err.(*Error); ok {
			return exErr.StackTrace()
		}
	}

	return ""
}

// FullError returns the complete error message, including private (non-user-facing) parts
func FullError(err error) string {
	if err != nil {
		if exErr, ok := err.(*Error); ok {
			return exErr.FullError()
		}
		return err.Error()
	}

	return ""
}

// GetAttrs returns error attributes or nil
func GetAttrs(err error) Attributes {
	if err != nil {
		if exErr, ok := err.(*Error); ok {
			return exErr.GetAttrs()
		}
	}

	return nil
}
