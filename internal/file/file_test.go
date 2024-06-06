package file

import (
	"bufio"
	"os"
	"testing"

	"github.com/dyne/tgcom/internal/comment"
)

func TestProcessFile(t *testing.T) {
	commentChars := map[string]string{"singleLine": "+++"}

	t.Run("SingleLine", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		lineNum := [2]int{2, 2}
		if err := ProcessFile(tmpFile.Name(), lineNum, commentChars, comment.Comment); err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		expected := "Line 1\n+++ Line 2\nLine 3\nLine 4\n"
		assertFileContent(t, tmpFile.Name(), expected)
	})

	t.Run("NonExistingFile", func(t *testing.T) {
		nonExistingFile := "non_existing_file.txt"
		lineNum := [2]int{2, 2}
		err := ProcessFile(nonExistingFile, lineNum, commentChars, comment.Comment)
		if err == nil {
			t.Fatalf("Expected error for non-existing file, got nil")
		}
	})

	t.Run("EmptyFile", func(t *testing.T) {
		emptyFile, cleanup := createTempFile(t, "")
		defer cleanup()

		lineNum := [2]int{2, 2}
		err := ProcessFile(emptyFile.Name(), lineNum, commentChars, comment.Comment)
		if err == nil {
			t.Fatalf("Expected error for empty file, got nil")
		}
	})

	// Test case: Multiple lines modification
	t.Run("MultipleLines", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		lineNum := [2]int{1, 3}
		if err := ProcessFile(tmpFile.Name(), lineNum, commentChars, comment.Comment); err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		expected := "+++ Line 1\n+++ Line 2\n+++ Line 3\nLine 4\n"
		assertFileContent(t, tmpFile.Name(), expected)
	})
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

		modFunc := func(line string, commentChars map[string]string) string {
			return comment.Comment(line, commentChars)
		}

		err := ProcessSingleFile(tt.filename, tt.lines, modFunc)
		if (err != nil) != tt.shouldErr {
			t.Errorf("ProcessSingleFile(%s, %s) error = %v, wantErr %v", tt.filename, tt.lines, err, tt.shouldErr)
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

		err := ProcessMultipleFiles(tt.fileInfo)
		if (err != nil) != tt.shouldErr {
			t.Errorf("ProcessMultipleFiles(%s) error = %v, wantErr %v", tt.fileInfo, err, tt.shouldErr)
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
			t.Errorf("extractLines(%s) error = %v, wantErr %v", tt.lineStr, err, tt.shouldErr)
		}
		if start != tt.startLine || end != tt.endLine {
			t.Errorf("extractLines(%s) = (%d, %d), want (%d, %d)", tt.lineStr, start, end, tt.startLine, tt.endLine)
		}
	}
}

func TestSelectCommentChars(t *testing.T) {
	tests := []struct {
		filename      string
		expectedChars map[string]string
		shouldErr     bool
	}{
		{"testfile.go", map[string]string{"singleLine": "//", "multiLineStart": "/*", "multiLineEnd": "*/"}, false},
		{"testfile.py", nil, true},
	}

	for _, tt := range tests {
		commentChars, err := selectCommentChars(tt.filename)
		if (err != nil) != tt.shouldErr {
			t.Errorf("selectCommentChars(%s) error = %v, wantErr %v", tt.filename, err, tt.shouldErr)
		}
		if !tt.shouldErr && !equalMaps(commentChars, tt.expectedChars) {
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

func equalMaps(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if b[k] != v {
			return false
		}
	}
	return true
}
