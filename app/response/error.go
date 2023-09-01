package response

import (
	"fmt"
	"net/http"
)

type Error interface {
	StatusCode() int
	ErrorCode() int
	Msg() string

	error
}

func NewError(statusCode, errorCode int, msg string) Error {
	return &BasicResponseError{
		statusCode: statusCode,
		errorCode:  errorCode,
		msg:        msg,
	}
}

var _ Error = &BasicResponseError{}

type BasicResponseError struct {
	statusCode int
	errorCode  int
	msg        string
}

func (e *BasicResponseError) ErrorCode() int {
	return e.errorCode
}

func (e *BasicResponseError) Msg() string {
	return e.msg
}

func (e *BasicResponseError) StatusCode() int {
	return e.statusCode
}

func (e *BasicResponseError) Error() string {
	return fmt.Sprintf("%d: %s", e.ErrorCode(), e.Msg())
}

func Warp(err error) Error {
	return ErrorInternalError(50000, err.Error())
}

func ErrorBadRequest(errorCode int, msg string) Error {
	return &BasicResponseError{
		statusCode: http.StatusBadRequest,
		errorCode:  errorCode,
		msg:        msg,
	}
}

func ErrorNotFound(errorCode int, msg string) Error {
	return &BasicResponseError{
		statusCode: http.StatusNotFound,
		errorCode:  errorCode,
		msg:        msg,
	}
}

func ErrorInternalError(errorCode int, msg string) Error {
	return &BasicResponseError{
		statusCode: http.StatusInternalServerError,
		errorCode:  errorCode,
		msg:        msg,
	}
}

func ErrorUnauthorized(errorCode int, msg string) Error {
	return &BasicResponseError{
		statusCode: http.StatusUnauthorized,
		errorCode:  errorCode,
		msg:        msg,
	}
}
