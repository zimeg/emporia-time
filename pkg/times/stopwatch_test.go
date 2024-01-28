package times

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBufferWriter(t *testing.T) {
	tests := map[string]struct {
		bounds       string
		output       []string
		expectedBuff []string
		expectedStd  []string
	}{
		"outputs are split between buff and std": {
			bounds: "xoxoxox",
			output: []string{
				"something command related",
				"that spans multiple lines",
				"xoxoxox",
				"information from the time",
			},
			expectedStd: []string{
				"something command related",
				"that spans multiple lines",
			},
			expectedBuff: []string{
				"information from the time",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			std := &bytes.Buffer{}
			bw := bufferWriter{
				bounds: tt.bounds,
				buff:   &bytes.Buffer{},
				std:    std,
			}
			for _, line := range tt.output {
				n, err := bw.Write([]byte(line))
				assert.NoError(t, err)
				assert.Equal(t, len(line), n)
			}
			assert.Equal(t, strings.Join(tt.expectedStd, ""), std.String())
			assert.Equal(t, strings.Join(tt.expectedBuff, ""), bw.buff.String())
		})
	}
}

func TestMakeBounds(t *testing.T) {
	t.Run("bounds are different between runs", func(t *testing.T) {
		bounds1 := makeBounds()
		bounds2 := makeBounds()
		assert.NotEqual(t, bounds1, bounds2)
	})
	t.Run("bounds are somewhat long", func(t *testing.T) {
		bounds := makeBounds()
		assert.Greater(t, len(bounds), 40)
	})
	t.Run("bounds ends in a newline", func(t *testing.T) {
		bounds := makeBounds()
		assert.Equal(t, len(bounds)-1, len(strings.TrimRight(bounds, "\n")))
	})
}
