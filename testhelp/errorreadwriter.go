package testhelp

import "errors"

var (
	READ_FAIL_ERROR  error = errors.New("copy failed on read")
	WRITE_FAIL_ERROR error = errors.New("copy failed on write")
)

type ErrorReadWriter struct{}

func (r *ErrorReadWriter) Read(p []byte) (n int, err error) {
	err = READ_FAIL_ERROR
	return
}

func (r *ErrorReadWriter) Write(p []byte) (n int, err error) {
	err = WRITE_FAIL_ERROR
	return
}
