package modfile

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/dyne/tgcom/utils/commenter"
)

// Config holds configuration settings for modifying files based on comments.
type Config struct {
	Filename   string
	LineNum    string
	StartLabel string
	EndLabel   string
	Lang       string
	Action     string
	DryRun     bool
}

func setModFunc(action string) (func(string, string) string, error) {
	switch action {
	case "comment":
		return commenter.Comment, nil
	case "uncomment":
		return commenter.Uncomment, nil
	case "toggle":
		return commenter.ToggleComments, nil
	case "":
		// If no action provided, assume toggle
		return commenter.ToggleComments, nil
	default:
		return nil, fmt.Errorf("invalid action. Please provide 'comment', 'uncomment', or 'toggle'")

	}
}

// This function process the input
func ChangeFile(conf Config) error {
	var file *os.File
	var err error
	var isStdin bool

	if conf.Filename == "" {
		// Read from stdin
		file = os.Stdin
		isStdin = true
	} else {
		// Open the file
		file, err = os.Open(conf.Filename)
		if err != nil {
			return err
		}
		defer file.Close()
	}

	char, err := selectCommentChars(conf.Filename, conf.Lang)
	if err != nil {
		return err
	}
	modFunc, err := setModFunc(conf.Action)
	if err != nil {
		return err
	}
	lines := [2]int{0, 0}
	if conf.LineNum != "" {
		lines, err = findLines(conf.LineNum)
		if err != nil {
			return err
		}
	}
	if conf.DryRun {
		err := printChanges(file, lines, conf.StartLabel, conf.EndLabel, char, modFunc)
		if err != nil {
			return fmt.Errorf("failed to process the file: %s", err)
		}
	} else {
		if isStdin {
			err := printOutput(file, lines, conf.StartLabel, conf.EndLabel, char, modFunc)
			if err != nil {
				return fmt.Errorf("failed to process the file: %s", err)
			}
		} else {
			// Create a backup of the original file
			backupFilename := conf.Filename + ".bak"
			if err := createBackup(conf.Filename, backupFilename); err != nil {
				return err
			}

			// Create a temporary file
			tmpFilename := conf.Filename + ".tmp"
			tmpFile, err := os.Create(tmpFilename)
			if err != nil {
				restoreBackup(conf.Filename, backupFilename)
				return err
			}
			defer tmpFile.Close()

			_, err = file.Seek(0, io.SeekStart)
			if err != nil {
				restoreBackup(conf.Filename, backupFilename)
				tmpFile.Close()
				os.Remove(tmpFilename)
				return err
			}

			err = writeChanges(file, tmpFile, lines, conf.StartLabel, conf.EndLabel, char, modFunc)
			if err != nil {
				restoreBackup(conf.Filename, backupFilename)
				tmpFile.Close()
				os.Remove(tmpFilename)
				return err
			}

			if err := file.Close(); err != nil {
				restoreBackup(conf.Filename, backupFilename)
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
			if err := os.Rename(tmpFilename, conf.Filename); err != nil {
				restoreBackup(conf.Filename, backupFilename)
				return err
			}

			// Remove backup file after successful processing
			os.Remove(backupFilename)
		}
	}
	return nil
}

func shouldProcessLine(currentLine int, lineNum [2]int, startLabel, endLabel string, inSection bool) bool {
	if startLabel != "" && endLabel != "" {
		return inSection
	}
	return lineNum[0] <= currentLine && currentLine <= lineNum[1]
}

func writeChanges(inputFile *os.File, outputFile *os.File, lineNum [2]int, startLabel, endLabel string, commentChars string, modFunc func(string, string) string) error {
	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	currentLine := 1
	inSection := false
	foundStart := false
	foundEnd := false
	var err error

	for scanner.Scan() {
		lineContent := scanner.Text()
		if strings.Contains(lineContent, endLabel) {
			foundEnd = true
			inSection = false
		}

		if shouldProcessLine(currentLine, lineNum, startLabel, endLabel, inSection) {
			lineContent = modFunc(lineContent, commentChars)
		}

		if strings.Contains(lineContent, startLabel) {
			foundStart = true
			inSection = true
		}

		if _, err = writer.WriteString(lineContent + "\n"); err != nil {
			return err
		}

		currentLine++
	}

	if lineNum[1] > currentLine && startLabel == "" && endLabel == "" {
		return errors.New("line number is out of range")
	}
	if !foundStart {
		return errors.New("start label not found in file")
	}

	if !foundEnd {
		return errors.New("end label not found in file")
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
	foundStart := false
	foundEnd := false

	for scanner.Scan() {

		lineContent := scanner.Text()

		if strings.Contains(lineContent, endLabel) {
			inSection = false
			foundEnd = true
		}

		if shouldProcessLine(currentLine, lineNum, startLabel, endLabel, inSection) {
			modified := modFunc(lineContent, commentChars)
			fmt.Printf("%d: %s -> %s\n", currentLine, lineContent, modified)
		}

		if strings.Contains(lineContent, startLabel) {
			inSection = true
			foundStart = true
		}

		currentLine++
	}

	if lineNum[1] > currentLine && startLabel == "" && endLabel == "" {
		return errors.New("line number is out of range")
	}
	if !foundStart {
		return errors.New("start label not found in file")
	}

	if !foundEnd {
		return errors.New("end label not found in file")
	}

	return scanner.Err()
}

func printOutput(input *os.File, lineNum [2]int, startLabel, endLabel, commentChars string, modFunc func(string, string) string) error {
	scanner := bufio.NewScanner(input)
	currentLine := 1
	inSection := false

	for scanner.Scan() {

		lineContent := scanner.Text()

		if strings.Contains(lineContent, endLabel) {
			inSection = false
		}

		if shouldProcessLine(currentLine, lineNum, startLabel, endLabel, inSection) {
			lineContent = modFunc(lineContent, commentChars)

		}

		if strings.Contains(lineContent, startLabel) {
			inSection = true
		}
		fmt.Println(lineContent)
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

func findLines(lineStr string) ([2]int, error) {
	if strings.Contains(lineStr, "-") {
		parts := strings.Split(lineStr, "-")
		if len(parts) != 2 {
			return [2]int{0, 0}, fmt.Errorf("invalid range format. Use 'start-end'")
		}
		startLine, err := strconv.Atoi(parts[0])
		if err != nil || startLine <= 0 {
			return [2]int{0, 0}, fmt.Errorf("invalid start line number")
		}
		endLine, err := strconv.Atoi(parts[1])
		if err != nil || endLine < startLine {
			return [2]int{0, 0}, fmt.Errorf("invalid end line number")
		}
		return [2]int{startLine, endLine}, nil
	} else {
		startLine, err := strconv.Atoi(lineStr)
		if err != nil || startLine <= 0 {
			return [2]int{0, 0}, fmt.Errorf("please provide a valid positive integer for the line number or a range")
		}
		endLine := startLine
		return [2]int{startLine, endLine}, nil
	}
}

func selectCommentChars(filename, lang string) (string, error) {
	if lang != "" {
		lang = strings.ToLower(lang)
		commentChars, ok := CommentChars[lang]
		if !ok {
			return "", fmt.Errorf("unsupported language: %s", lang)
		}
		return commentChars, nil
	}

	if filename != "" {
		extension := filepath.Ext(filename)
		switch extension {
		case ".go":
			return CommentChars["golang"], nil
		case ".js":
			return CommentChars["js"], nil
		case ".sh", ".bash":
			return CommentChars["bash"], nil
		case ".cpp", ".cc", ".h", ".c":
			return CommentChars["C"], nil
		case ".java":
			return CommentChars["java"], nil
		case ".py":
			return CommentChars["python"], nil
		case ".rb":
			return CommentChars["ruby"], nil
		case ".pl":
			return CommentChars["perl"], nil
		case ".php":
			return CommentChars["php"], nil
		case ".swift":
			return CommentChars["swift"], nil
		case ".kt", ".kts":
			return CommentChars["kotlin"], nil
		case ".R":
			return CommentChars["r"], nil
		case ".hs":
			return CommentChars["haskell"], nil
		case ".sql":
			return CommentChars["sql"], nil
		case ".rs":
			return CommentChars["rust"], nil
		case ".scala":
			return CommentChars["scala"], nil
		case ".dart":
			return CommentChars["dart"], nil
		case ".mm":
			return CommentChars["objective-c"], nil
		case ".m":
			return CommentChars["matlab"], nil
		case ".lua":
			return CommentChars["lua"], nil
		case ".erl":
			return CommentChars["erlang"], nil
		case ".ex", ".exs":
			return CommentChars["elixir"], nil
		case ".ts":
			return CommentChars["ts"], nil
		case ".vhdl", ".vhd":
			return CommentChars["vhdl"], nil
		case ".v", ".sv":
			return CommentChars["verilog"], nil
		case ".html":
			return CommentChars["html"], nil
		default:
			return "", fmt.Errorf("unsupported file extension: %s", extension)
		}
	}

	return "", fmt.Errorf("language not specified and no filename provided")
}

// CommentChars maps programming languages to their respective comment syntax.
var CommentChars = map[string]string{
	"golang":      "//",
	"go":          "//",
	"js":          "//",
	"bash":        "#",
	"c":           "//",
	"c++":         "//",
	"java":        "//",
	"python":      "#",
	"ruby":        "#",
	"perl":        "#",
	"php":         "//",
	"swift":       "//",
	"kotlin":      "//",
	"r":           "#",
	"haskell":     "--",
	"sql":         "--",
	"rust":        "//",
	"scala":       "//",
	"dart":        "//",
	"objective-c": "//",
	"matlab":      "%",
	"lua":         "--",
	"erlang":      "%",
	"elixir":      "#",
	"ts":          "//",
	"vhdl":        "--",
	"verilog":     "//",
	"html":        "<!-- -->",
}
