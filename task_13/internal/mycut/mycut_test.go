package mycut

import (
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestRunCut(t *testing.T) {
	tests := []struct {
		name    string
		options *CutOptions
		input   string
		want    []string
	}{
		{
			name: "single field",
			options: &CutOptions{
				Fields:    []int{0},
				Delimiter: ",",
				Separate:  false,
			},
			input: "a,b,c\nx,y,z\n",
			want:  []string{"a", "x"},
		},
		{
			name: "multiple fields",
			options: &CutOptions{
				Fields:    []int{0, 2},
				Delimiter: ",",
				Separate:  false,
			},
			input: "a,b,c\nx,y,z\n",
			want:  []string{"a,c", "x,z"},
		},
		{
			name: "separate=true ignores lines without delimiter",
			options: &CutOptions{
				Fields:    []int{0},
				Delimiter: ",",
				Separate:  true,
			},
			input: "abc\n1,2,3\nxyz\n",
			want:  []string{"1"},
		},
		{
			name: "custom delimiter",
			options: &CutOptions{
				Fields:    []int{1},
				Delimiter: ";",
				Separate:  false,
			},
			input: "a;b;c\nx;y;z\n",
			want:  []string{"b", "y"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			reader := strings.NewReader(tt.input)

			got, err := runCut(reader, tt.options)

			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
