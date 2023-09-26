// myapp/myerrors/myerrors.go
package myerrors

import "errors"

// Define custom errors as variables or constants.
var (
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidInput   = errors.New("invalid input")
	ErrAuthentication = errors.New("authentication failed")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrInternalServer = errors.New("internal server error")
	// Add more custom errors as needed
)
