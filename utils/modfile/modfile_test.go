package modfile

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestWriteChanges(t *testing.T) {
	// Create a temporary test file
	testFilename := "testfile.txt"
	defer os.Remove(testFilename)

	// Write some content to the test file
	file, err := os.Create(testFilename)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("start\nline 1\nline 2\nline 3\nend\n")
	if err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}

	// Define test cases
	tests := []struct {
		name         string
		lineNum      [2]int
		startLabel   string
		endLabel     string
		commentChars string
		modFunc      func(string, string) string
		expected     string
	}{
		{
			name:         "Test with lines",
			lineNum:      [2]int{2, 4},
			startLabel:   "",
			endLabel:     "",
			commentChars: "//",
			modFunc: func(line, commentChars string) string {
				return "// " + line
			},
			expected: "start\n// line 1\n// line 2\n// line 3\nend\n",
		},
		{
			name:         "Test with labels",
			lineNum:      [2]int{0, 0},
			startLabel:   "start",
			endLabel:     "end",
			commentChars: "//",
			modFunc: func(line, commentChars string) string {
				return "// " + line
			},
			expected: "start\n// line 1\n// line 2\n// line 3\nend\n",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Open test file for reading
			file, err := os.Open(testFilename)
			if err != nil {
				t.Fatalf("failed to open test file: %v", err)
			}
			defer file.Close()

			// Create a temporary output file
			outputFilename := "output.txt"
			outputFile, err := os.Create(outputFilename)
			if err != nil {
				t.Fatalf("failed to create output file: %v", err)
			}
			defer func() {
				outputFile.Close()
				os.Remove(outputFilename)
			}()

			// Call writeChanges function
			err = writeChanges(file, outputFile, tt.lineNum, tt.startLabel, tt.endLabel, tt.commentChars, tt.modFunc)
			if err != nil {
				t.Fatalf("writeChanges returned an error: %v", err)
			}

			// Read the output file
			outputFile.Close()
			outputFile, err = os.Open(outputFilename)
			if err != nil {
				t.Fatalf("failed to open output file: %v", err)
			}
			defer outputFile.Close()
			assertFileContent(t, outputFile.Name(), tt.expected)
		})
	}
}

func TestPrintChanges(t *testing.T) {
	// Create a temporary test file
	testFilename := "testfile.txt"
	defer os.Remove(testFilename)

	// Write some content to the test file
	file, err := os.Create(testFilename)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString("start\nline 1\nline 2\nline 3\nend\n")
	if err != nil {
		t.Fatalf("failed to write to test file: %v", err)
	}

	// Define test cases
	tests := []struct {
		name         string
		lineNum      [2]int
		startLabel   string
		endLabel     string
		commentChars string
		modFunc      func(string, string) string
		expected     string
	}{
		{
			name:         "Test with lines",
			lineNum:      [2]int{2, 4},
			startLabel:   "",
			endLabel:     "",
			commentChars: "//",
			modFunc: func(line, commentChars string) string {
				return "// " + line
			},
			expected: "2: line 1 -> // line 1\n3: line 2 -> // line 2\n4: line 3 -> // line 3\n",
		},
		{
			name:         "Test with labels",
			lineNum:      [2]int{0, 0},
			startLabel:   "start",
			endLabel:     "end",
			commentChars: "//",
			modFunc: func(line, commentChars string) string {
				return "// " + line
			},
			expected: "2: line 1 -> // line 1\n3: line 2 -> // line 2\n4: line 3 -> // line 3\n",
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Open test file for reading
			file, err := os.Open(testFilename)
			if err != nil {
				t.Fatalf("failed to open test file: %v", err)
			}
			defer file.Close()
			old := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			// Redirect stdout to buffer

			// Call printChanges function
			err = printChanges(file, tt.lineNum, tt.startLabel, tt.endLabel, tt.commentChars, tt.modFunc)
			if err != nil {
				t.Fatalf("printChanges returned an error: %v", err)
			}
			w.Close()
			var buf bytes.Buffer
			io.Copy(&buf, r)
			r.Close()
			os.Stdout = old

			// Compare output with expected
			output := buf.String()
			if output != tt.expected {
				t.Errorf("Unexpected output:\nGot:\n%s\nExpected:\n%s", output, tt.expected)
			}
		})
	}
}

var dryRun = false

func TestChangeFile(t *testing.T) {

	t.Run("SingleLine", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		conf := Config{
			Filename:   tmpFile.Name(),
			LineNum:    "2",
			StartLabel: "",
			EndLabel:   "",
			Lang:       "GoLang",
			Action:     "comment",
			DryRun:     dryRun,
		}
		ChangeFile(conf)

		expected := "Line 1\n// Line 2\nLine 3\nLine 4\n"
		assertFileContent(t, tmpFile.Name(), expected)
	})

	t.Run("MultipleLines", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		conf := Config{
			Filename:   tmpFile.Name(),
			LineNum:    "2-3",
			StartLabel: "",
			EndLabel:   "",
			Lang:       "GoLang",
			Action:     "comment",
			DryRun:     dryRun,
		}
		ChangeFile(conf)

		expected := "Line 1\n// Line 2\n// Line 3\nLine 4\n"
		assertFileContent(t, tmpFile.Name(), expected)
	})

	t.Run("Labels", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Start Label\nLine 1\nLine 2\nLine 3\nEnd Label")
		defer cleanup()

		conf := Config{
			Filename:   tmpFile.Name(),
			LineNum:    "",
			StartLabel: "Start Label",
			EndLabel:   "End Label",
			Lang:       "GoLang",
			Action:     "comment",
			DryRun:     dryRun,
		}
		ChangeFile(conf)

		expected := "Start Label\n// Line 1\n// Line 2\n// Line 3\nEnd Label\n"
		assertFileContent(t, tmpFile.Name(), expected)
	})

	t.Run("DryRun", func(t *testing.T) {
		tmpFile, cleanup := createTempFile(t, "Line 1\nLine 2\nLine 3\nLine 4")
		defer cleanup()

		conf := Config{
			Filename:   tmpFile.Name(),
			LineNum:    "2",
			StartLabel: "",
			EndLabel:   "",
			Lang:       "GoLang",
			Action:     "comment",
			DryRun:     true,
		}

		// Redirect stdout temporarily to capture dry run output
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		ChangeFile(conf)

		// Capture the output
		w.Close()
		var buf bytes.Buffer
		io.Copy(&buf, r)
		r.Close()
		os.Stdout = old

		got := buf.String()
		expected := "2: Line 2 -> // Line 2\n"
		if got != expected {
			t.Errorf("Dry run log does not match.\nExpected: %s\nGot: %s", expected, got)
		}
	})
	t.Run("Stdin", func(t *testing.T) {
		input := "line 1\nline 2\nline 3\nline 4\n"
		conf := Config{
			Filename:   "",
			LineNum:    "1-3",
			StartLabel: "",
			EndLabel:   "",
			Lang:       "GoLang",
			Action:     "comment",
			DryRun:     false,
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

		ChangeFile(conf)

		wStdout.Close()
		var buf bytes.Buffer
		io.Copy(&buf, rStdout)
		rStdout.Close()

		// Check the output
		got := buf.String()
		expected := "// line 1\n// line 2\n// line 3\nline 4\n"
		if got != expected {
			t.Errorf("expected %q, got %q", expected, got)
		}
	})

}

func TestCreateBackup(t *testing.T) {
	content := "Line 1\nLine 2\n"
	tmpFile, cleanup := createTempFile(t, content)
	defer cleanup()

	backupFile, cleanupBackup := createTempFile(t, "")
	defer cleanupBackup()

	if err := createBackup(tmpFile.Name(), backupFile.Name()); err != nil {
		t.Fatalf("createBackup() error = %v", err)
	}

	assertFileContent(t, backupFile.Name(), content)
}

func TestRestoreBackup(t *testing.T) {
	content := "Line 1\nLine 2\n"
	tmpFile, cleanup := createTempFile(t, content)
	defer cleanup()

	backupFile, cleanupBackup := createTempFile(t, content)
	defer cleanupBackup()

	// Corrupt the original file
	if err := os.WriteFile(tmpFile.Name(), []byte("corrupted content"), 0644); err != nil {
		t.Fatalf("Failed to corrupt file: %v", err)
	}

	restoreBackup(tmpFile.Name(), backupFile.Name())

	assertFileContent(t, tmpFile.Name(), content)
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
		if !tt.shouldErr && commentChars != tt.expectedChars {
			t.Errorf("selectCommentChars(%s) = %v, want %v", tt.filename, commentChars, tt.expectedChars)
		}
	}
}
func TestFindLines(t *testing.T) {
	tests := []struct {
		lineStr   string
		expected  [2]int
		shouldErr bool
	}{
		{"1-3", [2]int{1, 3}, false},
		{"2", [2]int{2, 2}, false},
		{"-1", [2]int{0, 0}, true},
		{"2-1", [2]int{0, 0}, true},
		{"invalid", [2]int{0, 0}, true},
	}

	for _, tt := range tests {
		result, err := findLines(tt.lineStr)
		if (err != nil) != tt.shouldErr {
			t.Fatalf("findLines(%s) error = %v, wantErr %v", tt.lineStr, err, tt.shouldErr)
		}
		if !tt.shouldErr && result != tt.expected {
			t.Errorf("findLines(%s) = %v, want %v", tt.lineStr, result, tt.expected)
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
