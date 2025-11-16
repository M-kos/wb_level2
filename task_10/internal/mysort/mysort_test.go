package mysort

import (
	"github.com/stretchr/testify/require"
	"os"
	"strings"
	"testing"
)

func TestRunSort(t *testing.T) {
	tests := []struct {
		name   string
		flags  []string
		input  string
		result []string
	}{
		{
			name:   "numeric sort",
			flags:  []string{"-n"},
			input:  "10\n2\n5\n",
			result: []string{"2", "5", "10"},
		},
		{
			name:   "reverse numeric sort",
			flags:  []string{"-nr"},
			input:  "10\n2\n5\n",
			result: []string{"10", "5", "2"},
		},
		{
			name:   "unique numeric",
			flags:  []string{"-nu"},
			input:  "1\n2\n1\n3\n2\n",
			result: []string{"1", "2", "3"},
		},
		{
			name:   "month sort",
			flags:  []string{"-M"},
			input:  "Mar\nJan\nDec\nFeb\n",
			result: []string{"Jan", "Feb", "Mar", "Dec"},
		},
		{
			name:   "human numeric sort",
			flags:  []string{"-h"},
			input:  "3K\n200K\n1G\n150K\n",
			result: []string{"3K", "150K", "200K", "1G"},
		},
		{
			name:   "reverse month",
			flags:  []string{"-Mr"},
			input:  "Mar\nJan\nDec\nFeb\n",
			result: []string{"Dec", "Mar", "Feb", "Jan"},
		},
		{
			name:   "ignore leading blanks",
			flags:  []string{"-b"},
			input:  " apple\n  kiwi\nbanana\n",
			result: []string{" apple", "banana", "  kiwi"},
		},
		{
			name:   "column sort",
			flags:  []string{"-k2", "-n"},
			input:  "apple\t10\nbanana\t2\ncherry\t5\n",
			result: []string{"banana\t2", "cherry\t5", "apple\t10"},
		},
		{
			name:   "numeric + unique + reverse",
			flags:  []string{"-nru"},
			input:  "5\n1\n3\n1\n2\n",
			result: []string{"5", "3", "2", "1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Args = append([]string{"program_name"}, tt.flags...)

			result, err := runSort(strings.NewReader(tt.input))

			require.NoError(t, err)
			require.Equal(t, tt.result, result)
		})
	}
}
