package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/dyne/tgcom/internal/comment"
	"github.com/dyne/tgcom/internal/file"
)

func main() {

	fileFlag := flag.String("file", "", "The file to process")
	lineFlag := flag.String("line", "", "The line number or range to modify (e.g., 4 or 10-20)")
	actionFlag := flag.String("action", "", "can be comment, uncomment or toggle")

	flag.Parse()

	filename := *fileFlag
	lineStr := *lineFlag
	action := *actionFlag

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

	if filename == "" {
		fmt.Println("Please provide a filename to process.")
		flag.PrintDefaults()
		return
	}

	if strings.Contains(filename, ",") {
		if err := file.ProcessMultipleFiles(filename); err != nil {
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
		if err := file.ProcessSingleFile(filename, lineStr, modFunc); err != nil {
			fmt.Println("Error processing file:", err)
		}
	}
}
