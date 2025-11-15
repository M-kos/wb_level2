package unpacker

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUnpack(t *testing.T) {
	tests := []struct {
		input       string
		result      string
		expectedErr error
	}{
		{
			input:       "a4bc2d5e",
			result:      "aaaabccddddde",
			expectedErr: nil,
		},
		{
			input:       "abcd",
			result:      "abcd",
			expectedErr: nil,
		},
		{
			input:       "45",
			result:      "",
			expectedErr: ErrStringMustNotContainOnlyNumbers,
		},
		{
			input:       "",
			result:      "",
			expectedErr: nil,
		},
		{
			input:       "qwe\\4\\5",
			result:      "qwe45",
			expectedErr: nil,
		},
		{
			input:       "qwe\\45",
			result:      "qwe44444",
			expectedErr: nil,
		},
		{
			input:       "ab45cd",
			result:      "",
			expectedErr: ErrInvalidStringFormat,
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			t.Parallel()
			res, err := Unpack(tt.input)
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.result, res)
		})
	}
}
