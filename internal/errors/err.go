package errors

import (
	"errors"
	"fmt"
	"strings"
)

// errorCode gives meaning
type errorCode string

// Err explains the problem
type Err struct {
	Code    errorCode
	Message string
	Source  error
}

// Error formats a sequence of error values
func (err Err) Error() string {
	response := []string{
		fmt.Sprintf("code=%s", string(err.Code)),
	}
	if err.Message != "" {
		response = append(response, fmt.Sprintf("message=%s", err.Message))
	}
	if err.Source != nil {
		response = append(response, fmt.Sprintf("source=%s", err.Source.Error()))
	}
	return strings.Join(response, "\x00")
}

// Is checks if target is found in this error
func (err Err) Is(target error) bool {
	detail := Unwrap(target)
	if detail.Code == err.Code {
		return true
	}
	if err.Source != nil {
		return Is(err.Source, target)
	}
	return false
}

// As check if tree errors match target
func As(tree error, target any) bool {
	return errors.As(tree, target)
}

// Is check if tree errors match target
func Is(tree error, target error) bool {
	return errors.Is(tree, target)
}

// Unwrap converts an error into stacks of details
func Unwrap(in error) (err Err) {
	splits := strings.Split(in.Error(), "\x00")
	for ii, detail := range splits {
		switch {
		case strings.HasPrefix(detail, "code="):
			code := errorCode(strings.TrimPrefix(detail, "code="))
			err = New(code)
		case strings.HasPrefix(detail, "message="):
			message := strings.TrimPrefix(detail, "message=")
			err.Message = message
		case strings.HasPrefix(detail, "source="):
			rest := strings.Join(splits[ii:], "\x00")
			source := strings.TrimPrefix(rest, "source=")
			detail := fmt.Errorf("%s", source)
			if strings.HasPrefix(source, "code=") {
				err.Source = Unwrap(detail)
			} else {
				err.Source = detail
			}
			return err
		default:
			err = New(errorCode(detail))
		}
	}
	return err
}

// Wrap sets the error source of an error
func Wrap(code errorCode, source error) Err {
	err := New(code)
	err.Source = source
	return err
}
