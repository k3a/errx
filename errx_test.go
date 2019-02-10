package errx_test

import (
	"database/sql"
	"strings"
	"testing"

	"ndemiccreations.com/gameservices/pkg/errx"
)

func TestPrivate(t *testing.T) {
	pubMsg := "unexpected database error"
	err := errx.Err(sql.ErrNoRows).Private(pubMsg)

	if err.Error() != pubMsg {
		t.Fatalf("short error is not returning public message: %s != %s", err.Error(), pubMsg)
	}

	if len(err.FullError()) <= len(pubMsg) {
		t.Fatalf("full error is too short: %s", err.FullError())
	}
}

func TestNilErr(t *testing.T) {
	var err *errx.Error
	if len(err.Error()) < 0 {
		t.Fatal("negative length of nil error message?!")
	}

	if len(err.FullError()) < 0 {
		t.Fatal("negative length of nil error message?!")
	}
}

func databaseError() error {
	return errx.Err(sql.ErrNoRows, "database error").Attr("db", "mydb")
}

func requestErrorPrivate() error {
	if err := databaseError(); err != nil {
		return errx.Err(err).Private("unable to save data")
	}

	return nil
}

func requestErrorPrivateWithMessage() error {
	if err := databaseError(); err != nil {
		return errx.Err(err, "error processing request").Private("unable to save data")
	}

	return nil
}

func requestErrorPublicCallingPrivate() error {
	if err := requestErrorPrivate(); err != nil {
		return errx.Err(err, "public err").Attr("server", "west-12")
	}
	return nil
}

func TestHierarchyPrivate(t *testing.T) {
	// public part
	err := requestErrorPrivate()
	if err.Error() != `unable to save data` {
		t.Fatalf("unexpected public error format: %v", err)
	}

	// private part
	fullErrStr := errx.FullError(err)
	if fullErrStr != `unable to save data: database error: sql: no rows in result set` {
		t.Fatalf("unexpected full error format: %s", fullErrStr)
	}

	// private part with additional message
	err = requestErrorPrivateWithMessage()
	fullErrStr = errx.FullError(err)
	if fullErrStr != `(unable to save data) error processing request: database error: sql: no rows in result set` {
		t.Fatalf("unexpected full error format with message: %s", fullErrStr)
	}

	// public part with 2 levels, then private
	err = requestErrorPublicCallingPrivate()
	if err.Error() != `public err: unable to save data` {
		t.Fatalf("unexpected 2-levels of public error format: %v", err)
	}
}

func TestAttrs(t *testing.T) {
	err := requestErrorPublicCallingPrivate()

	attrs := errx.GetAttrs(err)
	if attrs == nil {
		t.Fatal("no attributes at all!")
	}
	if attrs["server"] != "west-12" {
		t.Fatal("missing or wrong 'server' attribute")
	}
	if attrs["db"] != "mydb" {
		t.Fatal("missing or wrong 'db' attribute from the deeper error")
	}
}

func requestError() error {
	if err := databaseError(); err != nil {
		return errx.Err(err)
	}

	return nil
}

func requestErrorWithMessage() error {
	if err := databaseError(); err != nil {
		return errx.Err(err, "error processing request")
	}

	return nil
}

func TestHierarchy(t *testing.T) {
	err := requestError()
	if err.Error() != `database error: sql: no rows in result set` {
		t.Fatalf("unexpected public error format: %v", err)
	}

	// --
	err = requestErrorWithMessage()
	if err.Error() != `error processing request: database error: sql: no rows in result set` {
		t.Fatalf("unexpected public error format with message: %v", err)
	}
}

func TestNoParent(t *testing.T) {
	err := errx.Err("some error")
	if err.Error() != "some error" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestFormat(t *testing.T) {
	err := errx.Err("str %s, float %.2f, int %d", "ok", 1.23, -45)
	if err.Error() != "str ok, float 1.23, int -45" {
		t.Fatalf("unexpected error message: %s", err.Error())
	}
}

func TestStackTrace(t *testing.T) {
	errx.RecordStackTrace = true

	err := requestErrorPublicCallingPrivate()

	if !strings.Contains(errx.StackTrace(err), ".go") {
		t.Fatal("stack tracing doesn't work")
	}

	// -- when disabled
	errx.RecordStackTrace = false

	err = requestErrorPublicCallingPrivate()

	if errx.StackTrace(err) != "" {
		t.Fatal("stack trace is not empty on an error without recorded pc")
	}

}

func TestTopLevelFuncsWithNils(t *testing.T) {
	var err error

	if errx.FullError(err) != "" {
		t.Fatal("FullError on nil returns something")
	}

	if errx.StackTrace(err) != "" {
		t.Fatal("StackTrace on nil returns something")
	}

	if errx.GetAttrs(err) != nil {
		t.Fatal("GetAttrs on nil returns something")
	}
}

func TestTopLevelFuncsOnOrdinaryErrors(t *testing.T) {
	if errx.FullError(sql.ErrNoRows) != sql.ErrNoRows.Error() {
		t.Fatal("FullError doesn't work for standard error objects")
	}
}
