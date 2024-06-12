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
func ProcessFile(filename string, lineNum [2]int, startLabel, endLabel string, commentChars string, modFunc func(string, string) string, dryRun bool) error {
	inputFile, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	if dryRun {
		// Perform a dry run: print the changes instead of writing them
		return printChanges(inputFile, lineNum, startLabel, endLabel, commentChars, modFunc)
	}

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

	err = writeChanges(inputFile, tmpFile, lineNum, startLabel, endLabel, commentChars, modFunc)
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

// processes input from stdin.
func ProcessStdin(lineStr, startLabel, endLabel, lang string, modFunc func(string, string) string, dryRun bool) error {
	var lineNum [2]int
	if startLabel == "" && endLabel == "" {
		startLine, endLine, err := extractLines(lineStr)
		if err != nil {
			return err
		}
		lineNum = [2]int{startLine, endLine}
	}
	commentChars, err := selectCommentChars("", lang)
	if err != nil {
		return err
	}

	input := os.Stdin

	if dryRun {
		return printChanges(input, lineNum, startLabel, endLabel, commentChars, modFunc)
	}

	// Process input from stdin directly
	scanner := bufio.NewScanner(input)
	currentLine := 1
	inSection := false
	for scanner.Scan() {
		lineContent := scanner.Text()

		// Determine if we are processing based on line numbers or labels
		if startLabel != "" && endLabel != "" {
			if strings.Contains(lineContent, startLabel) {
				inSection = true
			}
			if inSection {
				lineContent = modFunc(lineContent, commentChars)
			}
			if strings.Contains(lineContent, endLabel) {
				inSection = false
			}
		} else {
			if lineNum[0] <= currentLine && currentLine <= lineNum[1] {
				lineContent = modFunc(lineContent, commentChars)
			}
		}

		// Print the modified line to stdout
		fmt.Println(lineContent)
		currentLine++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func writeChanges(inputFile *os.File, outputFile *os.File, lineNum [2]int, startLabel, endLabel string, commentChars string, modFunc func(string, string) string) error {
	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	currentLine := 1
	inSection := false
	var err error

	for scanner.Scan() {
		lineContent := scanner.Text()
		if startLabel != "" && endLabel != "" {
			if strings.Contains(lineContent, startLabel) {
				inSection = true
			}
			if inSection {
				lineContent = modFunc(lineContent, commentChars)
			}
			if strings.Contains(lineContent, endLabel) {
				inSection = false
			}
		} else {
			if lineNum[0] <= currentLine && currentLine <= lineNum[1] {
				lineContent = modFunc(lineContent, commentChars)
			}
		}

		if _, err = writer.WriteString(lineContent + "\n"); err != nil {
			return err
		}

		currentLine++
	}

	if lineNum[1] > currentLine && startLabel == "" && endLabel == "" {
		return errors.New("line number is out of range")
	}

	if err = scanner.Err(); err != nil {
		return err
	}

	return writer.Flush()
}

func printChanges(inputFile *os.File, lineNum [2]int, startLabel, endLabel, commentChars string, modFunc func(string, string) string) error {
	scanner := bufio.NewScanner(inputFile)
	currentLine := 1
	inSection := false

	for scanner.Scan() {
		lineContent := scanner.Text()

		// Determine if we are processing based on line numbers or labels
		if startLabel != "" && endLabel != "" {
			if strings.Contains(lineContent, startLabel) {
				inSection = true
			}
			if inSection {
				modified := modFunc(lineContent, commentChars)
				fmt.Printf("%d: %s -> %s\n", currentLine, lineContent, modified)
			}
			if strings.Contains(lineContent, endLabel) {
				inSection = false
			}
		} else {
			if lineNum[0] <= currentLine && currentLine <= lineNum[1] {
				modified := modFunc(lineContent, commentChars)
				fmt.Printf("%d: %s -> %s\n", currentLine, lineContent, modified)
			}
		}

		currentLine++
	}

	if lineNum[1] > currentLine && startLabel == "" && endLabel == "" {
		return errors.New("line number is out of range")
	}

	return scanner.Err()
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
func ProcessSingleFile(filename string, lineStr, startLabel, endLabel string, modFunc func(string, string) string, dryRun bool) error {
	commentChars, err := selectCommentChars(filename, "")
	if err != nil {
		return err
	}

	var lineNum [2]int
	if startLabel == "" && endLabel == "" {
		startLine, endLine, err := extractLines(lineStr)
		if err != nil {
			return err
		}
		lineNum = [2]int{startLine, endLine}
	}

	return ProcessFile(filename, lineNum, startLabel, endLabel, commentChars, modFunc, dryRun)
}

// ProcessMultipleFiles processes multiple files specified by comma-separated filenames.
func ProcessMultipleFiles(filename string, dryRun bool) error {
	fileLine := strings.Split(filename, ",")
	for _, fileInfo := range fileLine {
		if err := processFileWithLines(fileInfo, dryRun); err != nil {
			return err
		}
	}
	return nil
}

func processFileWithLines(fileInfo string, dryRun bool) error {
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

	commentChars, err := selectCommentChars(file, "")
	if err != nil {
		return err
	}

	return ProcessFile(file, lineNum, "", "", commentChars, comment.ToggleComments, dryRun)
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

func selectCommentChars(filename, lang string) (string, error) {
	if lang != "" {
		lang = strings.ToLower(lang)
		commentChars, ok := language.CommentChars[lang]
		if !ok {
			return "", fmt.Errorf("unsupported language: %s", lang)
		}
		return commentChars, nil
	}

	if filename != "" {
		extension := filepath.Ext(filename)
		switch extension {
		case ".go":
			return language.CommentChars["golang"], nil
		case ".js":
			return language.CommentChars["js"], nil
		case ".sh", ".bash":
			return language.CommentChars["bash"], nil
		case ".cpp", ".cc", ".h", ".c":
			return language.CommentChars["C"], nil
		case ".java":
			return language.CommentChars["java"], nil
		case ".py":
			return language.CommentChars["python"], nil
		case ".rb":
			return language.CommentChars["ruby"], nil
		case ".pl":
			return language.CommentChars["perl"], nil
		case ".php":
			return language.CommentChars["php"], nil
		case ".swift":
			return language.CommentChars["swift"], nil
		case ".kt", ".kts":
			return language.CommentChars["kotlin"], nil
		case ".R":
			return language.CommentChars["r"], nil
		case ".hs":
			return language.CommentChars["haskell"], nil
		case ".sql":
			return language.CommentChars["sql"], nil
		case ".rs":
			return language.CommentChars["rust"], nil
		case ".scala":
			return language.CommentChars["scala"], nil
		case ".dart":
			return language.CommentChars["dart"], nil
		case ".mm":
			return language.CommentChars["objective-c"], nil
		case ".m":
			return language.CommentChars["matlab"], nil
		case ".lua":
			return language.CommentChars["lua"], nil
		case ".erl":
			return language.CommentChars["erlang"], nil
		case ".ex", ".exs":
			return language.CommentChars["elixir"], nil
		case ".ts":
			return language.CommentChars["ts"], nil
		case ".vhdl", ".vhd":
			return language.CommentChars["vhdl"], nil
		case ".v", ".sv":
			return language.CommentChars["verilog"], nil
		default:
			return "", fmt.Errorf("unsupported file extension: %s", extension)
		}
	}

	return "", fmt.Errorf("language not specified and no filename provided")
}
