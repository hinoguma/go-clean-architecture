package crosscuttinglayer

import (
	"context"
	"errors"
	"fmt"
	"runtime"
)

func ErrWrap(wrapped, wrapper error) error {
	if !IsBaseError(wrapped) {
		wrapped = NewError(wrapped, context.TODO())
	}
	if wrapper == nil {
		return wrapped
	}
	return errors.Join(wrapped, wrapper)
}

func ErrLift(err error, ctx context.Context) error {
	if !IsBaseError(err) {
		err = NewError(err, ctx)
	}
	return err
}

func NewError(err error, ctx context.Context) BaseError {
	e := BaseError{
		err:        err,
		stacktrace: NewStankTrace(4),
		attributes: make(map[string]interface{}),
		requestID:  "",
	}
	e.WithContext(ctx)
	return e
}

func IsBaseError(err error) bool {
	e := BaseError{}
	return errors.As(err, &e)
}

type BaseError struct {
	err        error
	category   string
	stacktrace StackTrace
	requestID  string
	attributes map[string]interface{}
}

func (e BaseError) Error() string {
	return fmt.Sprintf(
		`{}"`,
	)
}

func (e BaseError) Unwrap() error {
	return e.err
}

func (e *BaseError) Attr(key string, value interface{}) *BaseError {
	if e.attributes == nil {
		e.attributes = make(map[string]interface{})
	}
	e.attributes[key] = value
	return e
}

func (e *BaseError) WithContext(ctx context.Context) *BaseError {
	rid, ok := ctx.Value("requestId").(string)
	if ok {
		e.requestID = rid
	}
	return e
}

type StackTraceItem struct {
	Function string
	File     string
	Line     int
}

func (item StackTraceItem) JsonString() string {
	return fmt.Sprintf(`{"function":"%s","file":"%s","line":%d}`, item.Function, item.File, item.Line)
}

const MAX_STACK_DEPTH = 30

type StackTrace []StackTraceItem

func (trace StackTrace) JsonString() string {
	traceStr := ""
	for i, item := range trace {
		if i > 0 {
			traceStr = traceStr + ","
		}
		traceStr = traceStr + item.JsonString()
	}
	return fmt.Sprintf(`[%s]`, traceStr)
}

// if skip is 3, you can get traces from your function calling this.
func NewStankTrace(skip int) StackTrace {
	pc := make([]uintptr, MAX_STACK_DEPTH)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	traces := make([]StackTraceItem, 0, n)
	for {
		frame, more := frames.Next()
		traces = append(traces, StackTraceItem{
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
	return traces
}

const (
	CategoryAppUtilError string = "application util error"
	CategoryLibraryError string = "library error"
	CategoryDataNotFound string = "data not found"
	CategoryForbidden    string = "forbidden"
)

// error in golang standard library or imported library
func NewLibraryError(err error, ctx context.Context) LibraryError {
	be := NewError(err, ctx)
	be.category = CategoryLibraryError
	return LibraryError{BaseError: be}
}

type LibraryError struct {
	BaseError
}

// error in our codes
func NewAppUtilError(message string, ctx context.Context) AppUtilError {
	be := NewError(errors.New(message), ctx)
	be.category = CategoryAppUtilError
	return AppUtilError{BaseError: be}
}

type AppUtilError struct {
	BaseError
}

// data not found error
type DataNotFound struct {
	BaseError
}

func NewDataNotFound(err error, ctx context.Context) DataNotFound {
	be := NewError(err, ctx)
	be.category = CategoryDataNotFound
	return DataNotFound{BaseError: be}
}

func NewDataNotFoundWithID(err error, id interface{}, ctx context.Context) DataNotFound {
	e := NewDataNotFound(err, ctx)
	e.ID(id)
	return e
}

func (e *DataNotFound) ID(id interface{}) *DataNotFound {
	return e.Key("id", id)
}

func (e *DataNotFound) Key(field string, value interface{}) *DataNotFound {
	e.Attr(field, value)
	return e
}

func IsDataNotFound(err error) bool {
	if err != nil {
		return false
	}
	return errors.As(err, &DataNotFound{})
}

// forbidden
type Forbidden struct {
	BaseError
}

type EncodeFailed struct {
	BaseError
}

type DecodeFailed struct {
	BaseError
}
