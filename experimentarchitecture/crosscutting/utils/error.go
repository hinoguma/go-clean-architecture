package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
)

// error type constants in error_custom.go
type ErrorType string

type Error struct {
	Type   ErrorType
	Err    error
	traces []StackTrace
	attrs  map[string]interface{}
}

type HasError struct {
	Err Error
}

func (e Error) Error() string {
	msg := e.Err.Error()
	f := ErrorFormat{
		Message: msg,
		Traces:  e.traces,
	}
	b, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf(`{"message":"%s","traces":[]}`, msg)
	}
	return string(b)
}

func (e Error) ErrorWithAttrs(attrs map[string]interface{}) string {
	msg := e.Err.Error()
	f := ErrorFormat{
		Message:    msg,
		Traces:     e.traces,
		Attributes: attrs,
	}
	b, err := json.Marshal(f)
	if err != nil {
		return fmt.Sprintf(`{"message":"%s","traces":[]}`, msg)
	}
	return string(b)
}

func (e Error) UnWrap() error {
	return e.Err
}

func (e Error) Wrap(err error) error {
	if err == nil {
		return e
	}
	e.Err = fmt.Errorf("%w: %w", e.Err, err)
	return e
}

func (e Error) Is(target error) bool {
	if target == nil {
		return false
	}
	err, ok := target.(Error)
	if ok && err.Type == e.Type {
		return true
	}
	return errors.Is(e.Err, target)
}

func (e Error) As(target any) bool {
	if target == nil {
		return false
	}
	err, ok := target.(*Error)
	if ok && err != nil && err.Type == e.Type {
		*err = e
		return true
	}
	return errors.As(e.Err, target)
}

func (e *Error) SetAttr(key string, value interface{}) *Error {
	if e.attrs == nil {
		e.attrs = make(map[string]interface{})
	}
	e.attrs[key] = value
	return e
}

func (e *Error) ToValue() Error {
	return *e
}
func (e Error) ToPointer() *Error {
	return &e
}

func (e *Error) SetAttrId(value interface{}) *Error {
	e.SetAttr("id", value)
	return e
}

func NewError(message string) Error {
	// stacktrace
	pc := make([]uintptr, 10)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	traces := make([]StackTrace, 0, n)
	for {
		frame, more := frames.Next()
		traces = append(traces, StackTrace{
			Function: frame.Function,
			File:     frame.File,
			Line:     frame.Line,
		})
		if len(traces) >= n {
			break
		}
		if !more {
			break
		}
	}
	return Error{
		Type:   ErrorTypeGeneral,
		Err:    errors.New(message),
		traces: traces,
	}
}

type StackTrace struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type ErrorFormat struct {
	Message    string                 `json:"message"`
	Traces     []StackTrace           `json:"traces"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

func ErrWrap(wrapped error, wrapper error) error {
	if wrapped == nil && wrapper == nil {
		return wrapper
	}
	if wrapped == nil {
		return wrapper
	}
	if wrapper == nil {
		return wrapped
	}
	wrapperEntity, ok := wrapper.(interface {
		Wrap(err error) error
	})
	if ok {
		fmt.Println("ErrWrap: use Wrap method")
		return wrapperEntity.Wrap(wrapped)
	}
	fmt.Println("ErrWrap: use fmt.Errorf")
	return fmt.Errorf("%s: %w", wrapped, wrapper)
}

type ValidateErrorDetail struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
