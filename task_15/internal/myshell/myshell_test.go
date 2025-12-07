package myshell

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitByLogicOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []LogicPart
	}{
		{
			name:  "Simple command without operators",
			input: "echo hello",
			expected: []LogicPart{
				{cmd: "echo hello", operator: ""},
			},
		},
		{
			name:  "Command with && operator",
			input: "cd /tmp && pwd",
			expected: []LogicPart{
				{cmd: "cd /tmp", operator: "&&"},
				{cmd: "pwd", operator: ""},
			},
		},
		{
			name:  "Command with || operator",
			input: "false || echo success",
			expected: []LogicPart{
				{cmd: "false", operator: "||"},
				{cmd: "echo success", operator: ""},
			},
		},
		{
			name:  "Multiple && operators",
			input: "cmd1 && cmd2 && cmd3",
			expected: []LogicPart{
				{cmd: "cmd1", operator: "&&"},
				{cmd: "cmd2", operator: "&&"},
				{cmd: "cmd3", operator: ""},
			},
		},
		{
			name:  "Mixed operators",
			input: "cmd1 && cmd2 || cmd3",
			expected: []LogicPart{
				{cmd: "cmd1", operator: "&&"},
				{cmd: "cmd2", operator: "||"},
				{cmd: "cmd3", operator: ""},
			},
		},
		{
			name:  "Extra spaces around operators",
			input: "cmd1  &&  cmd2  ||  cmd3",
			expected: []LogicPart{
				{cmd: "cmd1", operator: "&&"},
				{cmd: "cmd2", operator: "||"},
				{cmd: "cmd3", operator: ""},
			},
		},
	}

	sh := NewMyShell()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sh.splitByLogicOperators(tt.input)
			require.Equal(t, len(tt.expected), len(result))

			for i, exp := range tt.expected {
				require.Equal(t, exp.cmd, result[i].cmd)
				require.Equal(t, exp.operator, result[i].operator)
			}
		})
	}
}

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Single argument",
			input:    "echo hello",
			expected: []string{"echo", "hello"},
		},
		{
			name:     "Multiple arguments",
			input:    "echo hello world foo",
			expected: []string{"echo", "hello", "world", "foo"},
		},
		{
			name:     "Extra spaces",
			input:    "echo   hello    world",
			expected: []string{"echo", "hello", "world"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Only spaces",
			input:    "   ",
			expected: []string{},
		},
	}

	sh := NewMyShell()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sh.parseArgs(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestExpandEnvVars(t *testing.T) {
	err := os.Setenv("TEST_VAR1", "test_value")
	require.NoError(t, err)
	err = os.Setenv("TEST_VAR2", "/home/user")
	require.NoError(t, err)

	defer func() {
		_ = os.Unsetenv("TEST_VAR1")
		_ = os.Unsetenv("TEST_VAR2")
	}()

	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "No environment variables",
			input:    []string{"echo", "hello"},
			expected: []string{"echo", "hello"},
		},
		{
			name:     "Single environment variable",
			input:    []string{"echo", "$TEST_VAR1"},
			expected: []string{"echo", "test_value"},
		},
		{
			name:     "Multiple environment variables",
			input:    []string{"$TEST_VAR1", "and", "$TEST_VAR2"},
			expected: []string{"test_value", "and", "/home/user"},
		},
		{
			name:     "Non-existent environment variable",
			input:    []string{"echo", "$NON_EXISTENT_VAR"},
			expected: []string{"echo", ""},
		},
	}

	sh := NewMyShell()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sh.expandEnvVars(tt.input)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestRunBuiltinEcho(t *testing.T) {
	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name: "echo without arguments",
			args: []string{"echo"},
		},
		{
			name: "echo with single argument",
			args: []string{"echo", "hello"},
		},
		{
			name: "echo with multiple arguments",
			args: []string{"echo", "hello", "world"},
		},
	}

	sh := NewMyShell()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.runBuiltin(tt.args)
			require.NoError(t, err)
		})
	}
}

func TestRunBuiltinCd(t *testing.T) {
	sh := NewMyShell()

	originalWd, err := os.Getwd()
	require.NoError(t, err)
	defer func() {
		_ = os.Chdir(originalWd)
	}()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "cd to temp directory",
			args:    []string{"cd", "/tmp"},
			wantErr: false,
		},
		{
			name:    "cd with default argument",
			args:    []string{"cd"},
			wantErr: false,
		},
		{
			name:    "cd to non-existent directory",
			args:    []string{"cd", "/non/existent/path"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.runBuiltin(tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRunBuiltinKill(t *testing.T) {
	sh := NewMyShell()

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "kill without arguments",
			args:    []string{"kill"},
			wantErr: true,
		},
		{
			name:    "kill with invalid pid",
			args:    []string{"kill", "invalid"},
			wantErr: true,
		},
		{
			name:    "kill with non-existent process",
			args:    []string{"kill", "99999999"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.runBuiltin(tt.args)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestRunBuiltinPs(t *testing.T) {
	sh := NewMyShell()

	err := sh.runBuiltin([]string{"ps"})
	require.NoError(t, err)
}

func TestRunBuiltinUnknownCommand(t *testing.T) {
	sh := NewMyShell()

	err := sh.runBuiltin([]string{"unknown_command"})
	require.Error(t, err)
}

func TestProcessRedirects(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	outputFile := filepath.Join(tmpDir, "output.txt")

	err := os.WriteFile(testFile, []byte("test data"), 0644)
	require.NoError(t, err)

	sh := NewMyShell()

	tests := []struct {
		name          string
		args          []string
		shouldHaveCmd bool
		cmdCount      int
	}{
		{
			name:          "No redirects",
			args:          []string{"arg1", "arg2"},
			shouldHaveCmd: false,
			cmdCount:      2,
		},
		{
			name:          "Redirect output without file",
			args:          []string{"arg1", ">"},
			shouldHaveCmd: false,
			cmdCount:      1,
		},
		{
			name:          "Redirect input",
			args:          []string{"<", testFile},
			shouldHaveCmd: false,
			cmdCount:      0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := tt.args
			if tt.name == "Redirect output without file" {
				args = []string{"arg1", ">", outputFile}
			}

			cmd := exec.Command("echo")
			result := sh.processRedirects(args, cmd)
			require.Equal(t, tt.cmdCount, len(result))
		})
	}
}

func TestProcessLine(t *testing.T) {
	sh := NewMyShell()

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Simple echo command",
			input:   "echo hello",
			wantErr: false,
		},
		{
			name:    "Empty command",
			input:   "",
			wantErr: false,
		},
		{
			name:    "Non-existent command",
			input:   "nonexistent_cmd_12345",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := sh.processLine(tt.input)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
