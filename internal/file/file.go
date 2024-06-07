package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dyne/tgcom/internal/comment"
	"github.com/dyne/tgcom/internal/language"
)

// ProcessFile processes a single file.
func ProcessFile(filename string, lineNum [2]int, commentChars string, modFunc func(string, string) string) error {
	inputFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Create a backup of the original file
	backupFilename := filename + ".bak"
	if err := createBackup(filename, backupFilename); err != nil {
		return err
	}

	// Create a temporary file
	tmpFilename := filename + ".tmp"
	tmpFile, err := os.Create(tmpFilename)
	if err != nil {
		restoreBackup(filename, backupFilename)
		return err
	}
	defer tmpFile.Close()

	if _, err := inputFile.Seek(0, io.SeekStart); err != nil {
		restoreBackup(filename, backupFilename)
		tmpFile.Close()
		os.Remove(tmpFilename)
		return err
	}

	err = writeChanges(inputFile, tmpFile, lineNum, commentChars, modFunc)
	if err != nil {
		restoreBackup(filename, backupFilename)
		tmpFile.Close()
		os.Remove(tmpFilename)
		return err
	}

	if err := inputFile.Close(); err != nil {
		restoreBackup(filename, backupFilename)
		tmpFile.Close()
		os.Remove(tmpFilename)
		return err
	}

	// Close the temporary file before renaming
	if err := tmpFile.Close(); err != nil {
		os.Remove(tmpFilename)
		return err
	}

	// Rename temporary file to original file
	if err := os.Rename(tmpFilename, filename); err != nil {
		restoreBackup(filename, backupFilename)
		return err
	}

	// Remove backup file after successful processing
	os.Remove(backupFilename)

	return nil
}

func writeChanges(inputFile *os.File, outputFile *os.File, lineNum [2]int, commentChars string, modFunc func(string, string) string) error {
	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	currentLine := 1

	for scanner.Scan() {
		lineContent := scanner.Text()
		if lineNum[0] <= currentLine && currentLine <= lineNum[1] {
			lineContent = modFunc(lineContent, commentChars)
		}

		if _, err := writer.WriteString(lineContent + "\n"); err != nil {
			return err
		}

		currentLine++
	}

	if lineNum[1] > currentLine {
		return errors.New("line number is out of range")
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return writer.Flush()
}
func createBackup(filename, backupFilename string) error {
	inputFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	backupFile, err := os.Create(backupFilename)
	if err != nil {
		return err
	}
	defer backupFile.Close()

	_, err = io.Copy(backupFile, inputFile)
	if err != nil {
		return err
	}

	return nil
}

func restoreBackup(filename, backupFilename string) {
	// Remove the potentially corrupted file
	os.Remove(filename)
	// Restore the backup file
	os.Rename(backupFilename, filename)
}

// ProcessSingleFile processes a single file specified by filename.
func ProcessSingleFile(filename string, lineStr string, modFunc func(string, string) string) error {
	startLine, endLine, err := extractLines(lineStr)
	if err != nil {
		return err
	}
	lineNum := [2]int{startLine, endLine}

	commentChars, err := selectCommentChars(filename)
	if err != nil {
		return err
	}

	return ProcessFile(filename, lineNum, commentChars, modFunc)
}

// ProcessMultipleFiles processes multiple files specified by comma-separated filenames.
func ProcessMultipleFiles(filename string) error {
	fileLine := strings.Split(filename, ",")
	for _, fileInfo := range fileLine {
		if err := processFileWithLines(fileInfo); err != nil {
			return err
		}
	}
	return nil
}

func processFileWithLines(fileInfo string) error {
	if !strings.Contains(fileInfo, ":") {
		return fmt.Errorf("invalid syntax format. Use '<filename>:<lines>'")
	}

	sub := strings.Split(fileInfo, ":")
	if len(sub) != 2 {
		return fmt.Errorf("invalid syntax format. Use '<filename>:<lines>'")
	}

	file, lineString := sub[0], sub[1]
	startLine, endLine, err := extractLines(lineString)
	if err != nil {
		return err
	}
	lineNum := [2]int{startLine, endLine}

	commentChars, err := selectCommentChars(file)
	if err != nil {
		return err
	}

	return ProcessFile(file, lineNum, commentChars, comment.ToggleComments)
}

func extractLines(lineStr string) (startLine, endLine int, err error) {
	if strings.Contains(lineStr, "-") {
		parts := strings.Split(lineStr, "-")
		if len(parts) != 2 {
			return 0, 0, fmt.Errorf("invalid range format. Use 'start-end'")
		}
		startLine, err = strconv.Atoi(parts[0])
		if err != nil || startLine <= 0 {
			return 0, 0, fmt.Errorf("invalid start line number")
		}
		endLine, err = strconv.Atoi(parts[1])
		if err != nil || endLine < startLine {
			return 0, 0, fmt.Errorf("invalid end line number")
		}
	} else {
		startLine, err = strconv.Atoi(lineStr)
		if err != nil || startLine <= 0 {
			return 0, 0, fmt.Errorf("please provide a valid positive integer for the line number")
		}
		endLine = startLine
	}
	return
}

func selectCommentChars(filename string) (string, error) {
	extension := filepath.Ext(filename)
	var commentChars string
	switch extension {
	case ".go":
		commentChars = language.CommentChars["GoLang"]
	case ".js":
		commentChars = language.CommentChars["JS"]
	case ".sh", ".bash":
		commentChars = language.CommentChars["Bash"]
	case ".cpp", ".cc", ".h", ".c":
		commentChars = language.CommentChars["C++/C"]
	case ".java":
		commentChars = language.CommentChars["Java"]
	case ".py":
		commentChars = language.CommentChars["Pyhton"]
	case ".rb":
		commentChars = language.CommentChars["Ruby"]
	case ".pl":
		commentChars = language.CommentChars["Perl"]
	case ".php":
		commentChars = language.CommentChars["PHP"]
	case ".swift":
		commentChars = language.CommentChars["swift"]
	case ".kt", ".kts":
		commentChars = language.CommentChars["Kotlin"]
	case ".R":
		commentChars = language.CommentChars["R"]
	case ".hs":
		commentChars = language.CommentChars["Haskell"]
	case ".sql":
		commentChars = language.CommentChars["SQL"]
	case ".rs":
		commentChars = language.CommentChars["Rust"]
	case ".scala":
		commentChars = language.CommentChars["Scala"]
	case ".dart":
		commentChars = language.CommentChars["Dart"]
	case ".mm":
		commentChars = language.CommentChars["Objective-C"]
	case ".m":
		commentChars = language.CommentChars["MATLAB"]
	case ".lua":
		commentChars = language.CommentChars["Lua"]
	case ".erl":
		commentChars = language.CommentChars["Erlang"]
	case ".ex", ".exs":
		commentChars = language.CommentChars["Elixir"]
	case ".ts":
		commentChars = language.CommentChars["TS"]
	case ".vhdl", ".vhd":
		commentChars = language.CommentChars["VHDL"]
	case ".v", ".sv":
		commentChars = language.CommentChars["Verilog"]
	default:
		return "", fmt.Errorf("unsupported file extension: %s", extension)
	}
	return commentChars, nil
}
