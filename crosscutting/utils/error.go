package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"runtime"
)

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

type Error struct {
	Err    error
	Traces []StackTrace
}

func (e Error) Error() string {
	msg := e.Err.Error()
	f := ErrorFormat{
		Message: msg,
		Traces:  e.Traces,
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
		Traces:     e.Traces,
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
	fmt.Println("Error Is method called")
	_, ok := target.(Error)
	if ok {
		return true
	}
	return errors.Is(e.Err, target)
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
		Err:    errors.New(message),
		Traces: traces,
	}
}

type HasError struct {
	Err Error
}

type DataNotFoundError struct {
	HasError
	KeyName  string
	KeyValue interface{}
}

func NewDataNotFoundError(keyName string, keyValue interface{}) DataNotFoundError {
	return DataNotFoundError{
		HasError: HasError{
			Err: NewError("data not found"),
		},
		KeyName:  keyName,
		KeyValue: keyValue,
	}
}

func (e DataNotFoundError) Error() string {
	return e.Err.ErrorWithAttrs(
		map[string]interface{}{
			"type":     "data not found",
			"keyName":  e.KeyName,
			"keyValue": e.KeyValue,
		},
	)
}

func (e DataNotFoundError) Wrap(err error) error {
	_ = e.Err.Wrap(err)
	return e
}

func (e DataNotFoundError) Unwrap() error {
	return e.Err.UnWrap()
}

func (e DataNotFoundError) Is(target error) bool {
	fmt.Println("DataNotFoundError Is method called")
	var dataNotFoundError DataNotFoundError
	ok := errors.As(target, &dataNotFoundError)
	if ok {
		return true
	}
	return e.Err.Is(target)
}

type ForbiddenError struct {
	HasError
}

func (e ForbiddenError) Error() string {
	return e.Err.ErrorWithAttrs(
		map[string]interface{}{
			"type": "forbidden",
		},
	)
}

func (e *ForbiddenError) Wrap(err error) error {
	_ = e.Err.Wrap(err)
	return e
}

func (e ForbiddenError) Unwrap() error {
	return e.Err.UnWrap()
}

func (e ForbiddenError) Is(target error) bool {
	_, ok := target.(ForbiddenError)
	if ok {
		return true
	}
	return e.Err.Is(target)
}

func ErrWrap(wrapped error, wrap error) error {
	if wrapped == nil {
		if wrap == nil {
			return nil
		}
		return wrap
	}
	if wrap == nil {
		return wrapped
	}
	u, ok := wrapped.(interface {
		Wrap(err error) error
	})
	if ok && u != nil {
		fmt.Println("ErrWrap: use Wrap method")
		return u.Wrap(wrap)
	}
	fmt.Println("ErrWrap: use fmt.Errorf")
	return fmt.Errorf("%s: %w", wrapped.Error(), wrap)
}

func func1(cnt int) error {
	if cnt <= 5 {
		return NewError("cnt is less than or equal to 5")
	}
	err := func2(cnt)
	if err != nil {
		return fmt.Errorf("func1: %w", err)
	}
	return nil
}

func func2(cnt int) error {
	if cnt <= 10 {
		return NewError("cnt is less than or equal to 10")
	}
	return nil
}

// BadRequestError

// ValidationError

// UnexpectedFormatError

// InvalidFormatError

// DataLockedError

// UNAUTHORIZEDError

// FORBIDDENError

// TIMEOUTError

// ConnectionFailedError

// UNKNOWNError
