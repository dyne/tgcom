package file

import (
	"os"
	"testing"

	"github.com/dyne/tgcom/internal/comment"
)

func TestProcessFile(t *testing.T) {

	commentChars := map[string]string{"singleLine": "+++"}
	t.Run("Temporary file", func(t *testing.T) {
		// Create a temporary file
		tmpFile, err := os.CreateTemp("", "testfile")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		// Write some test data to the temporary file
		initialContent := `Line 1
Line 2
Line 3
Line 4`
		if _, err := tmpFile.WriteString(initialContent); err != nil {
			t.Fatalf("Failed to write to temporary file: %v", err)
		}

		if err := tmpFile.Close(); err != nil {
			t.Fatalf("Failed to close temporary file: %v", err)
		}

		lineNum := 2
		if err := ProcessFile(tmpFile.Name(), lineNum, commentChars, comment.Comment); err != nil {
			t.Fatalf("ProcessFile failed: %v", err)
		}

		// Read the modified file
		modified, err := os.ReadFile(tmpFile.Name())
		if err != nil {
			t.Fatalf("Failed to read modified file: %v", err)
		}

		// Check the content of the modified file
		expected := `Line 1
+++ Line 2
Line 3
Line 4
`
		if string(modified) != expected {
			t.Errorf("Unexpected file content:\nGot:\n%s\nExpected:\n%s", string(modified), expected)
		}
	})

	// Test for non-existing file
	t.Run("NonExistingFile", func(t *testing.T) {
		nonExistingFile := "non_existing_file.txt"
		lineNum := 2
		err := ProcessFile(nonExistingFile, lineNum, commentChars, comment.Comment)
		if err == nil {
			t.Fatalf("ProcessFile did not return an error for non-existing file: %v", nonExistingFile)
		}

	})

	// Test for empty file
	t.Run("EmptyFile", func(t *testing.T) {
		emptyFile, err := os.CreateTemp("", "emptyfile")
		if err != nil {
			t.Fatalf("Failed to create temporary file: %v", err)
		}
		defer os.Remove(emptyFile.Name())

		lineNum := 2
		err = ProcessFile(emptyFile.Name(), lineNum, commentChars, comment.Comment)
		if err == nil {
			t.Fatalf("ProcessFile did not return an error for empty file: %v", emptyFile.Name())
		}
	})
}
