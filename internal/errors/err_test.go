package errors

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsError(t *testing.T) {
	tests := map[string]struct {
		err      Err
		expected string
	}{
		"formats an individual error": {
			err: Err{
				Code:    "err_example_problem",
				Message: "a reason or cause",
			},
			expected: "code=err_example_problem\x00message=a reason or cause",
		},
		"formats multiple errors with sources": {
			err: Err{
				Code:    "err_example_problem",
				Message: "a reason or cause",
				Source: Err{
					Code:    "err_example_issue",
					Message: "a breaking thing",
					Source: Err{
						Code:    "err_example_permissions",
						Message: "the error in question",
					},
				},
			},
			expected: "code=err_example_problem\x00message=a reason or cause\x00source=code=err_example_issue\x00message=a breaking thing\x00source=code=err_example_permissions\x00message=the error in question",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := tt.err.Error()
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestErrorsAs(t *testing.T) {
	tests := map[string]struct {
		tree     error
		target   any
		expected bool
	}{
		"matching errors are same": {
			tree:     New("err_emporia_maintenance"),
			target:   New("err_emporia_maintenance"),
			expected: true,
		},
		"nested errors are same": {
			tree: Err{
				Code: "err_emporia_maintenance",
				Source: Err{
					Code: "err_emporia_status",
				},
			},
			target:   New("err_emporia_status"),
			expected: true,
		},
		"different errors are different": {
			tree:     New("err_emporia_maintenance"),
			target:   New("err_emporia_status"),
			expected: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := As(tt.tree, &tt.target)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestErrorsIs(t *testing.T) {
	tests := map[string]struct {
		tree     error
		target   error
		expected bool
	}{
		"matching errors are same": {
			tree:     New("err_emporia_unplugged"),
			target:   New("err_emporia_unplugged"),
			expected: true,
		},
		"nested errors are same": {
			tree: Err{
				Code: "err_emporia_unplugged",
				Source: Err{
					Code: "err_emporia_status",
				},
			},
			target:   New("err_emporia_status"),
			expected: true,
		},
		"different errors are different": {
			tree:     New("err_emporia_unplugged"),
			target:   New("err_emporia_status"),
			expected: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := Is(tt.tree, tt.target)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestErrorsUnwrap(t *testing.T) {
	tests := map[string]struct {
		err      error
		expected Err
	}{
		"formats an individual error": {
			err: fmt.Errorf("code=err_emporia_unplugged\x00message=a reason or cause"),
			expected: Err{
				Code:    "err_emporia_unplugged",
				Message: "a reason or cause",
			},
		},
		"formats an unknown error": {
			err: fmt.Errorf("oops"),
			expected: Err{
				Code:    "err_unexpected_problem",
				Message: "oops",
			},
		},
		"formats multiple errors with sources": {
			err: fmt.Errorf("code=err_emporia_checkup\x00message=a reason or cause\x00source=code=err_emporia_status\x00message=a breaking thing\x00source=code=err_emporia_maintenance\x00message=the error in question\x00source=some downtimes"),
			expected: Err{
				Code:    "err_emporia_checkup",
				Message: "a reason or cause",
				Source: Err{
					Code:    "err_emporia_status",
					Message: "a breaking thing",
					Source: Err{
						Code:    "err_emporia_maintenance",
						Message: "the error in question",
						Source:  fmt.Errorf("some downtimes"),
					},
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := Unwrap(tt.err)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestErrorsWrap(t *testing.T) {
	tests := map[string]struct {
		code     errorCode
		source   error
		expected Err
	}{
		"extends an unknown error": {
			code:   "err_emporia_checkup",
			source: fmt.Errorf("oops"),
			expected: Err{
				Code:    "err_emporia_checkup",
				Message: "failed to check uptime",
				Source:  fmt.Errorf("oops"),
			},
		},
		"extends a known error": {
			code: "err_emporia_status",
			source: Err{
				Code:    "err_emporia_maintenance",
				Message: "the error in question",
				Source:  fmt.Errorf("some downtimes"),
			},
			expected: Err{
				Code:    "err_emporia_status",
				Message: "failed to get status",
				Source: Err{
					Code:    "err_emporia_maintenance",
					Message: "the error in question",
					Source:  fmt.Errorf("some downtimes"),
				},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			actual := Wrap(tt.code, tt.source)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
