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
func ProcessFile(filename string, lineNum [2]int, commentChars map[string]string, modFunc func(string, map[string]string) string) error {
	inputFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	// Create a temporary file
	tmpfilename := filename + ".tmp"
	tmpFile, err := os.Create(tmpfilename)
	if err != nil {
		return err
	}
	defer tmpFile.Close()

	if _, err := inputFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	err = writeChanges(inputFile, tmpFile, lineNum, commentChars, modFunc)
	if err != nil {
		return err
	}

	// Close the original file
	if err := inputFile.Close(); err != nil {
		return err
	}

	// Rename temporary file to original file
	if err := os.Rename(tmpfilename, filename); err != nil {
		return err
	}

	return nil
}

func writeChanges(inputFile *os.File, outputFile *os.File, lineNum [2]int, commentChars map[string]string, modFunc func(string, map[string]string) string) error {
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

// ProcessSingleFile processes a single file specified by filename.
func ProcessSingleFile(filename string, lineStr string, modFunc func(string, map[string]string) string) error {
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

func selectCommentChars(filename string) (map[string]string, error) {
	extension := filepath.Ext(filename)
	var commentChars map[string]string
	switch extension {
	case ".go":
		commentChars = language.GoCommentChars
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", extension)
	}
	return commentChars, nil
}
