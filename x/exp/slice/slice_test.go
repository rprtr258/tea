package slice

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Take(t *testing.T) {
	for _, tc := range []struct {
		input    []int
		take     int
		expected []int
	}{
		{
			input:    []int{1, 2, 3, 4, 5},
			take:     3,
			expected: []int{1, 2, 3},
		},
		{
			input:    []int{1, 2, 3},
			take:     5,
			expected: []int{1, 2, 3},
		},
		{
			input:    []int{},
			take:     2,
			expected: []int{},
		},
		{
			input:    []int{1, 2, 3},
			take:     0,
			expected: []int{},
		},
		{
			input:    nil,
			take:     2,
			expected: []int{},
		},
	} {
		assert.Equal(t, tc.expected, Take(tc.input, tc.take))
	}
}
