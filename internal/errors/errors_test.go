package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestErrorsNew(t *testing.T) {
	tests := map[string]struct {
		code     errorCode
		expected Err
	}{
		"error code missing is unexpected": {
			expected: Err{
				Code: "err_unexpected_problem",
			},
		},
		"error code unknown is messaged": {
			code: "something strange",
			expected: Err{
				Code:    "err_unexpected_problem",
				Message: "something strange",
			},
		},
		"error code found has details": {
			code: ErrEmporiaMaintenance,
			expected: Err{
				Code:    "err_emporia_maintenance",
				Message: "cannot measure during maintenance",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := New(tt.code)
			assert.Equal(t, tt.expected, err)
		})
	}
}
