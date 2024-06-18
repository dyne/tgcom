package file

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/dyne/tgcom/internal/comment"
)

var dryRun = false

func TestProcessFile(t *testing.T) {
	commentChars := "+++"
	modFunc := func(line string, commentChars string) string {
		return commentChars + " " + line
	}
	t.Run("SingleLine", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		lineNum := [2]int{2, 2}
		if err := ProcessFile(tmpFile.Name(), lineNum, "", "", commentChars, modFunc, dryRun); err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		expected := "Line 1\n+++ Line 2\nLine 3\nLine 4\n"
		assertFileContent(t, tmpFile.Name(), expected)
	})

	t.Run("NonExistingFile", func(t *testing.T) {
		nonExistingFile := "non_existing_file.txt"
		lineNum := [2]int{2, 2}
		err := ProcessFile(nonExistingFile, lineNum, "", "", commentChars, modFunc, dryRun)
		if err == nil {
			t.Fatalf("Expected error for non-existing file, got nil")
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		emptyFile, cleanup := createTempFile(t, "")
		defer cleanup()

		lineNum := [2]int{2, 2}
		err := ProcessFile(emptyFile.Name(), lineNum, "", "", commentChars, modFunc, dryRun)
		if err == nil {
			t.Fatalf("Expected error for empty file, got nil")
		}
	})

	// Test case: Multiple lines modification
	t.Run("MultipleLines", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		lineNum := [2]int{1, 3}
		if err := ProcessFile(tmpFile.Name(), lineNum, "", "", commentChars, comment.Comment, dryRun); err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		expected := "+++ Line 1\n+++ Line 2\n+++ Line 3\nLine 4\n"
		assertFileContent(t, tmpFile.Name(), expected)
	})
	t.Run("WithBackup", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()
		modFunc := func(line string, commentChars string) string {
			return ""
		}

		lineNum := [2]int{3, 10}
		if err := ProcessFile(tmpFile.Name(), lineNum, "", "", commentChars, modFunc, dryRun); err == nil {
			t.Fatal("Expected error, got nil")
		}

		fileContent, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			t.Fatalf("Error reading backup file: %v", err)
		}

		expectedContent := "Line 1\nLine 2\nLine 3\nLine 4"
		if string(fileContent) != expectedContent {
			t.Errorf("Backup file content does not match:\nGot:\n%s\nExpected:\n%s", string(fileContent), expectedContent)
		}
	})
	t.Run("DryRun", func(t *testing.T) {
		// Test case for dry run
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		lineNum := [2]int{2, 2}
		dryRun = true
		defer func() { dryRun = false }()

		// Redirect stdout temporarily to capture dry run output
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		if err := ProcessFile(tmpFile.Name(), lineNum, "", "", commentChars, modFunc, dryRun); err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		// Capture the output
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		r.Close()
		os.Stdout = old

		got := buf.String()
		expected := "2: Line 2 -> +++ Line 2\n"
		if got != expected {
			t.Errorf("Dry run log does not match.\nExpected: %s\nGot: %s", expected, got)
		}
	})
	t.Run("Labels", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 0\nStart Label\nLine 1\nLine 2\nLine 3\nEnd Label\nLine 5")
		defer cleanup()

		startLabel := "Start Label"
		endLabel := "End Label"
		if err := ProcessFile(tmpFile.Name(), [2]int{}, startLabel, endLabel, commentChars, modFunc, dryRun); err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		expected := "Line 0\nStart Label\n+++ Line 1\n+++ Line 2\n+++ Line 3\nEnd Label\nLine 5"
		assertFileContent(t, tmpFile.Name(), expected)
	})
}

func TestProcessStdin(t *testing.T) {
	input := "line 1\nline 2\nline 3\nline 4\n"
	modFunc := func(line string, commentChars string) string {
		return commentChars + " " + line
	}

	// Create pipes for stdin and stdout redirection
	rStdin, wStdin, _ := os.Pipe()
	rStdout, wStdout, _ := os.Pipe()

	// Write the mock input to the writer end of the stdin pipe
	go func() {
		defer wStdin.Close()
		_, _ = wStdin.Write([]byte(input))
	}()

	// Save original stdin and stdout
	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()
	oldStdout := os.Stdout
	defer func() { os.Stdout = oldStdout }()

	// Redirect stdin and stdout
	os.Stdin = rStdin
	os.Stdout = wStdout

	err := ProcessStdin("1-3", "", "", "go", modFunc, false)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	wStdout.Close()

	// Read the captured output
	var buf bytes.Buffer
	io.Copy(&buf, rStdout)
	rStdout.Close()

	// Check the output
	got := buf.String()
	expected := "// line 1\n// line 2\n// line 3\nline 4\n"
	if got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}

}

func TestProcessSingleFile(t *testing.T) {
	// Setup test files with content
	tests := []struct {
		filename  string
		lines     string
		expected  string
		shouldErr bool
	}{
		{
			"temp_testfile1.go",
			"5-10",
			`line 1
line 2
line 3
line 4
// line 5
// line 6
// line 7
// line 8
// line 9
// line 10
line 11
line 12
`,
			false,
		},
		{
			"temp_testfile2.go",
			"15",
			"",
			true,
		},
	}

	for _, tt := range tests {
		setupTestFile(t, tt.filename, []string{
			"line 1", "line 2", "line 3", "line 4",
			"line 5", "line 6", "line 7", "line 8",
			"line 9", "line 10", "line 11", "line 12",
		})
		defer os.Remove(tt.filename)

		modFunc := func(line string, commentChars string) string {
			return commentChars + " " + line
		}

		err := ProcessSingleFile(tt.filename, tt.lines, "", "", modFunc, dryRun)
		if (err != nil) != tt.shouldErr {
			t.Errorf("ProcessSingleFile(%s, %s) error = %v", tt.filename, tt.lines, err)
		}
		if !tt.shouldErr {
			assertFileContent(t, tt.filename, tt.expected)
		}
	}
}

func TestProcessMultipleFiles(t *testing.T) {
	tests := []struct {
		fileInfo  string
		expected  map[string]string
		shouldErr bool
	}{
		{
			"temp_testfile3.go:1-3,temp_testfile4.go:2-4",
			map[string]string{
				"temp_testfile3.go": `// line 1
// line 2
// line 3
line 4
line 5
`,
				"temp_testfile4.go": `line 1
// line 2
// line 3
// line 4
line 5
line 6
`,
			},
			false,
		},
		{
			"temp_testfile5.go:1-3,non_existing_file.txt:2-4",
			nil,
			true,
		},
	}

	for _, tt := range tests {
		setupTestFile(t, "temp_testfile3.go", []string{"line 1", "line 2", "line 3", "line 4", "line 5"})
		setupTestFile(t, "temp_testfile4.go", []string{"line 1", "line 2", "line 3", "line 4", "line 5", "line 6"})
		defer os.Remove("temp_testfile3.go")
		defer os.Remove("temp_testfile4.go")

		err := ProcessMultipleFiles(tt.fileInfo, dryRun)
		if (err != nil) != tt.shouldErr {
			t.Errorf("ProcessMultipleFiles(%s) error = %v", tt.fileInfo, err)
		}
		if !tt.shouldErr {
			for filename, expectedContent := range tt.expected {
				assertFileContent(t, filename, expectedContent)
			}
		}
	}
}

func TestExtractLines(t *testing.T) {
	tests := []struct {
		lineStr   string
		startLine int
		endLine   int
		shouldErr bool
	}{
		{"5-10", 5, 10, false},
		{"15", 15, 15, false},
		{"invalid", 0, 0, true},
	}

	for _, tt := range tests {
		start, end, err := extractLines(tt.lineStr)
		if (err != nil) != tt.shouldErr {
			t.Errorf("extractLines(%s) error = %v", tt.lineStr, err)
		}
		if start != tt.startLine || end != tt.endLine {
			t.Errorf("extractLines(%s) = (%d, %d), want (%d, %d)", tt.lineStr, start, end, tt.startLine, tt.endLine)
		}
	}
}

func TestSelectCommentChars(t *testing.T) {
	tests := []struct {
		filename      string
		expectedChars string
		shouldErr     bool
	}{
		{"testfile.go", "//", false},
		{"testfile.false", "", true},
	}

	for _, tt := range tests {
		commentChars, err := selectCommentChars(tt.filename, "")
		if (err != nil) != tt.shouldErr {
			t.Errorf("selectCommentChars(%s) error = %v", tt.filename, err)
		}
		if !tt.shouldErr && !(commentChars == tt.expectedChars) {
			t.Errorf("selectCommentChars(%s) = %v, want %v", tt.filename, commentChars, tt.expectedChars)
		}
	}
}

// Utility functions

func createTempFile(t testing.TB, content string) (*os.File, func()) {
	t.Helper()
	tmpFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}
	return tmpFile, func() { os.Remove(tmpFile.Name()) }
}

func assertFileContent(t testing.TB, filename string, expected string) {
	t.Helper()
	modified, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read modified file: %v", err)
	}
	if string(modified) != expected {
		t.Errorf("Unexpected file content:\nGot:\n%s\nExpected:\n%s", string(modified), expected)
	}
}

func setupTestFile(t testing.TB, filename string, lines []string) {
	t.Helper()
	if err := createTestFile(filename); err != nil {
		t.Fatalf("Error creating test file: %v", err)
	}
	if err := writeTestContent(filename, lines); err != nil {
		t.Fatalf("Error writing test content: %v", err)
	}
}

func createTestFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	return f.Close()
}

func writeTestContent(filename string, lines []string) error {
	f, err := os.OpenFile(filename, os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	for _, line := range lines {
		if _, err := writer.WriteString(line + "\n"); err != nil {
			return err
		}
	}
	return writer.Flush()
}
