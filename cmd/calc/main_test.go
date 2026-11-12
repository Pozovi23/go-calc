package main_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var binPath string

func TestMain(m *testing.M) {
	tmpDir, err := os.MkdirTemp("", "calc-test-*")
	if err != nil {
		os.Exit(1)
	}
	defer func() {
		_ = os.RemoveAll(tmpDir)
	}()

	binPath = filepath.Join(tmpDir, "calc")
	cmd := exec.Command("go", "build", "-o", binPath, "./")
	if err := cmd.Run(); err != nil {
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func runCmd(input string, args ...string) (string, string, int) {
	cmd := exec.Command(binPath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if input != "" {
		cmd.Stdin = strings.NewReader(input)
	}
	err := cmd.Run()
	code := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			code = exitErr.ExitCode()
		} else {
			code = 1
		}
	}
	return stdout.String(), stderr.String(), code
}

func TestRun_ExpressionArgument(t *testing.T) {
	stdout, stderr, code := runCmd("", "2 + 2 * 2")
	if code != 0 {
		t.Fatalf("expected code 0, got %d, stderr: %s", code, stderr)
	}
	if strings.TrimSpace(stdout) != "6" {
		t.Errorf("expected stdout '6', got %q", stdout)
	}
}

func TestRun_PrecisionFlag(t *testing.T) {
	tests := []struct {
		name     string
		args     []string
		expected string
	}{
		{
			name:     "short precision flag -p 2",
			args:     []string{"-p", "2", "1 / 3"},
			expected: "0.33",
		},
		{
			name:     "default precision 10",
			args:     []string{"1 / 3"},
			expected: "0.3333333333",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, code := runCmd("", tt.args...)
			if code != 0 {
				t.Fatalf("expected code 0, got %d, stderr: %s", code, stderr)
			}
			if strings.TrimSpace(stdout) != tt.expected {
				t.Errorf("expected stdout %q, got %q", tt.expected, stdout)
			}
		})
	}
}

func TestRun_Errors(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedCode  int
		expectedError string
	}{
		{
			name:          "division by zero",
			args:          []string{"1 / 0"},
			expectedCode:  1,
			expectedError: "division by zero",
		},
		{
			name:          "syntax error",
			args:          []string{"2 +"},
			expectedCode:  1,
			expectedError: "syntax error",
		},
		{
			name:          "negative precision",
			args:          []string{"-p", "-1", "2 + 2"},
			expectedCode:  1,
			expectedError: "precision cannot be negative",
		},
		{
			name:          "invalid flag",
			args:          []string{"--invalid-flag"},
			expectedCode:  2,
			expectedError: "flag provided but not defined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, stderr, code := runCmd("", tt.args...)
			if code != tt.expectedCode {
				t.Fatalf("expected code %d, got %d", tt.expectedCode, code)
			}
			if !strings.Contains(stderr, tt.expectedError) {
				t.Errorf("expected stderr to contain %q, got %q", tt.expectedError, stderr)
			}
		})
	}
}

func TestRun_REPL(t *testing.T) {
	tests := []struct {
		name           string
		input          string
		expectedStdout string
		expectedStderr string
	}{
		{
			name:           "exit command",
			input:          "2 + 2\nexit\n",
			expectedStdout: "4",
		},
		{
			name:           "q command",
			input:          "3 * 4\nq\n",
			expectedStdout: "12",
		},
		{
			name:           "empty lines and error continuation",
			input:          "\n  \n1 / 0\n5 + 5\nq\n",
			expectedStdout: "10",
			expectedStderr: "division by zero",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stdout, stderr, code := runCmd(tt.input)
			if code != 0 {
				t.Fatalf("expected code 0, got %d, stderr: %s", code, stderr)
			}
			if !strings.Contains(stdout, tt.expectedStdout) {
				t.Errorf("expected stdout to contain %q, got %q", tt.expectedStdout, stdout)
			}
			if tt.expectedStderr != "" && !strings.Contains(stderr, tt.expectedStderr) {
				t.Errorf("expected stderr to contain %q, got %q", tt.expectedStderr, stderr)
			}
		})
	}
}
