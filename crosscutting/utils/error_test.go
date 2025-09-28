package utils

import (
	"errors"
	"testing"
)

func TestAssert(t *testing.T) {

	initialErr := NewError("initial error")
	norelevantErr := errors.New("no relevant error")

	layer1 := NewDataNotFoundError("id", "123")
	layer2 := ForbiddenError{
		HasError: HasError{
			Err: NewError("forbidden error"),
		},
	}
	wraped1 := ErrWrap(initialErr, layer1)
	test1 := errors.Is(wraped1, initialErr)
	t.Logf(
		"test1: %v, wraped1: %v", test1, wraped1,
	)

	wraped2 := ErrWrap(wraped1, layer2)
	test2 := errors.Is(wraped2, initialErr)
	t.Logf(
		"test2: %v,", test2,
	)
	test3 := errors.Is(wraped2, norelevantErr)
	t.Logf(
		"test3: %v,", test3,
	)

	// log only test result
	t.Logf("test1: %v, test2: %v, test3:%v", test1, test2, test3)
}
