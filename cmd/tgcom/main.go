//go:build !vcs

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/dyne/tgcom/internal/comment"
	"github.com/dyne/tgcom/internal/file"
)

func main() {
	fileFlag := flag.String("file", "", "The file to process")
	lineFlag := flag.String("line", "", "The line number or range to modify (e.g., 4 or 10-20)")
	startLabelFlag := flag.String("start-label", "", "The start label for a section")
	endLabelFlag := flag.String("end-label", "", "The end label for a section")
	actionFlag := flag.String("action", "", "can be comment, uncomment or toggle")
	dryRunFlag := flag.Bool("dry-run", false, "Perform a dry run without modifying the files")
	lang := flag.String("language", "", "Specify the programming language")
	flag.Parse()

	filename := *fileFlag
	lineStr := *lineFlag
	startLabel := *startLabelFlag
	endLabel := *endLabelFlag
	action := *actionFlag
	dryRun := *dryRunFlag
	langStr := *lang
	info, _ := os.Stdin.Stat()
	isStdin := (info.Mode() & os.ModeCharDevice) == 0

	var modFunc func(string, string) string

	switch action {
	case "comment":
		modFunc = comment.Comment
	case "uncomment":
		modFunc = comment.Uncomment
	case "toggle":
		modFunc = comment.ToggleComments
	case "":
		// If no action provided, assume toggle
		modFunc = comment.ToggleComments
	default:
		fmt.Println("Invalid action. Please provide 'comment', 'uncomment', or 'toggle'.")
		flag.PrintDefaults()
		return
	}

	if startLabel == "" && endLabel != "" {
		fmt.Println("Error: 'startLabel' is required when 'endLabel' is provided.")
		return
	} else if startLabel != "" && endLabel == "" {
		fmt.Println("Error: 'endLabel' is required when 'startLabel' is provided.")
		return
	}
	if startLabel != "" && lineStr != "" {
		fmt.Println("Error: Specify either line number/range OR label, not both.")
		return
	}

	if isStdin {
		if err := file.ProcessStdin(lineStr, langStr, startLabel, endLabel, modFunc, dryRun); err != nil {
			fmt.Println("Error processing files:", err)
		}
	} else {
		if filename == "" {
			fmt.Println("Please provide a filename to process.")
			flag.PrintDefaults()
			return
		}

		if strings.Contains(filename, ",") {
			if err := file.ProcessMultipleFiles(filename, dryRun); err != nil {
				fmt.Println("Error processing files:", err)
			}
		} else {
			if strings.Contains(filename, ":") {
				parts := strings.Split(filename, ":")
				if len(parts) != 2 {
					fmt.Println("Invalid syntax format. Use '<filename>:<lines>'")
					return
				}
				filename = parts[0]
				lineStr = parts[1]
			}
			if err := file.ProcessSingleFile(filename, lineStr, startLabel, endLabel, modFunc, dryRun); err != nil {
				fmt.Println("Error processing file:", err)
			}
		}
	}
}
