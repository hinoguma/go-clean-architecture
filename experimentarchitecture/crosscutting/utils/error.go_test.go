package utils

import (
	"errors"
	"testing"
)

func TestError_Is(t *testing.T) {
	type testCase struct {
		label    string
		err      error
		target   error
		expected bool
	}

	timeout := NewTimeoutError("5s")
	dataNotFound := NewDataNotFoundError()
	dataNotFound.SetAttrId("123")
	stErr := errors.New("some standard error")
	normalErr := NewError("some normal error")
	sameStdErr := errors.New("same standard error")

	layerd_Nrm_Nrm := ErrWrap(
		NewError("layer2"),
		NewError("layer1"),
	)
	layerd_Nrm_Timeout := ErrWrap(
		NewTimeoutError("10s"),
		NewError("layer1"),
	)
	layerd_Timeout_Nrm := ErrWrap(
		NewError("layer1"),
		NewTimeoutError("10s"),
	)
	layerd_DataNotFound_Std := ErrWrap(
		sameStdErr,
		NewDataNotFoundError().ToPointer().SetAttr("name", "abc").ToValue(),
	)
	layerd_Std_DataNotFound := ErrWrap(
		NewDataNotFoundError().ToPointer().SetAttr("name", "abc").ToValue(),
		sameStdErr,
	)
	cases := []testCase{
		{
			label:    "same error type timeout",
			err:      timeout,
			target:   NewTimeoutError("1s"),
			expected: true,
		},
		{
			label:    "same error type data not found",
			err:      dataNotFound,
			target:   NewDataNotFoundError().ToPointer().SetAttr("name", "abc").ToValue(),
			expected: true,
		},
		{
			label:    "different error type. timeout - data not found",
			err:      timeout,
			target:   dataNotFound,
			expected: false,
		},
		{
			label:    "different error type. standard error - data not found",
			err:      stErr,
			target:   dataNotFound,
			expected: false,
		},
		{
			label:    "different error type. normal error - data not found",
			err:      normalErr,
			target:   dataNotFound,
			expected: false,
		},
		{
			label:    "same error type on first layer. normal",
			err:      layerd_Timeout_Nrm,
			target:   NewError("another normal error"),
			expected: true,
		},
		{
			label:    "same error type on second layer. normal",
			err:      layerd_Nrm_Timeout,
			target:   NewError("another normal error"),
			expected: true,
		},
		{
			label:    "different error type in layers. normal",
			err:      layerd_DataNotFound_Std,
			target:   NewError("another normal error"),
			expected: false,
		},
		{
			label:    "same error type on first layer. timeout",
			err:      layerd_Timeout_Nrm,
			target:   NewTimeoutError("1s"),
			expected: true,
		},
		{
			label:    "same error type on second layer. timeout",
			err:      layerd_Nrm_Timeout,
			target:   NewTimeoutError("1s"),
			expected: true,
		},
		{
			label:    "different error type in layers. timeout",
			err:      layerd_Nrm_Nrm,
			target:   NewTimeoutError("1s"),
			expected: false,
		},
		{
			label:    "same error type on first layer. std",
			err:      layerd_Std_DataNotFound,
			target:   sameStdErr,
			expected: true,
		},
		{
			label:    "same error type on second layer. std",
			err:      layerd_DataNotFound_Std,
			target:   sameStdErr,
			expected: true,
		},
		{
			label:    "different error type in layers. std",
			err:      layerd_DataNotFound_Std,
			target:   errors.New("different standard error"),
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.label, func(t *testing.T) {
			res := errors.Is(c.err, c.target)
			if res != c.expected {
				t.Errorf("label: %s, expected: %v, result: %v", c.label, c.expected, res)
			}
		})
	}

}
