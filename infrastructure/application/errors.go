package application

import (
	"net/http"
)

//type HTTPError struct {
//	StatusCode int    `json:"statusCode"`
//	Message    string `json:"message"`
//}
//
//func (e HTTPError) Error() string {
//	message := fmt.Sprintf("%d %s", e.StatusCode, e.Message)
//	return message
//}

type ErrHTTP interface {
	Code() int
	Error() string
}

type ErrCustom struct {
	code    int
	message string
}

func NewErrCustom(code int, msg string) ErrCustom {
	return ErrCustom{code: code, message: msg}
}

func (e ErrCustom) Error() string {
	return e.message
}

func (e ErrCustom) Code() int {
	return e.code
}

type ErrConflict struct {
	message string
}

func NewErrConflict(msg string) ErrConflict {
	return ErrConflict{message: msg}
}

func (e ErrConflict) Error() string {
	return e.message
}

func (e ErrConflict) Code() int {
	return http.StatusConflict
}

type ErrNotFound struct {
	message string
}

func NewErrNotFound(msg string) ErrNotFound {
	return ErrNotFound{message: msg}
}

func (e ErrNotFound) Error() string {
	return e.message
}

func (e ErrNotFound) Code() int {
	return http.StatusNotFound
}

type ErrValidation struct {
	message string
}

func NewErrValidation(msg string) ErrValidation {
	return ErrValidation{message: msg}
}

func (e ErrValidation) Error() string {
	return e.message
}

func (e ErrValidation) Code() int {
	return http.StatusBadRequest
}

type ErrUnauthorized struct{}

func NewErrUnauthorized() ErrUnauthorized { return ErrUnauthorized{} }

func (e ErrUnauthorized) Error() string {
	return UNAUTHORIZED
}

func (e ErrUnauthorized) Code() int {
	return http.StatusUnauthorized
}

type ErrForbidden struct{}

func NewErrForbidden() ErrForbidden { return ErrForbidden{} }

func (e ErrForbidden) Error() string {
	return FORBIDDEN
}

func (e ErrForbidden) Code() int {
	return http.StatusForbidden
}
