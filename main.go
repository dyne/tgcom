package main

import (
    "fmt"
    "strconv"
    "strings"
    "flag"
    "path/filepath"
    "github.com/dyne/tgcom/internal/file"
    "github.com/dyne/tgcom/internal/comment"
    "github.com/dyne/tgcom/internal/language"
)
func main() {
	// interface for the user
	
	// definition of the flags for the command line: "file", "line"
	fileFlag := flag.String("file", "", "The file to process")
    lineFlag := flag.String("line", "", "The line number or range to modify (e.g., 4 or 10-20)")
	actionFlag := flag.String("action", "", "can be comment, uncomment or toogle")
    
    // Analyize the flags
    flag.Parse()

	// Print error message if flags arguments are empty
	if *fileFlag == "" || *lineFlag == "" || *actionFlag == "" {
        fmt.Println("Usage: go run main.go --file <filename> --line <line> --action <action>")
        return
    }

	// take arguments of the flags
    filename := *fileFlag
	lineStr := *lineFlag
    action := *actionFlag

	// find the startLine and the endLine specified by user, if only one
	// argument is given then startLine = endLine
	var startLine, endLine int
    var err error

    if strings.Contains(lineStr, "-") {
        parts := strings.Split(lineStr, "-")
        if len(parts) != 2 {
            fmt.Println("Invalid range format. Use 'start-end'.")
            return
        }
        startLine, err = strconv.Atoi(parts[0])
        if err != nil || startLine <= 0 {
            fmt.Println("Invalid start line number.")
            return
        }
        endLine, err = strconv.Atoi(parts[1])
        if err != nil || endLine < startLine {
            fmt.Println("Invalid end line number.")
            return
        }
    } else {
        startLine, err = strconv.Atoi(lineStr)
        if err != nil || startLine <= 0 {
            fmt.Println("Please provide a valid positive integer for the line number.")
            return
        }
        endLine = startLine
    }

	// Prepare the array that indicates lines to be commented
	lineNum := [2]int{startLine, endLine} 

    // need to know which file are you dealing with, so extract extension of the file
    extension := filepath.Ext(filename)

    // depending on the extension we select a different map from go.go
    var commentChars map[string]string
    switch extension {
        case ".go":
            commentChars = language.GoCommentChars
        /* TODO: more possible extensions
        case ".js":
        case ".html":
        */
    }

    // depending on the action we select a different function from comment.go
    var modFunc func(string, map[string]string) string
    switch action {
        case "comment":
            modFunc = comment.Comment
        case "uncomment":
            modFunc = comment.Uncomment
        case "toggle":
            modFunc = comment.ToggleComments
    }

    // operate on the file
    file.ProcessFile(filename, lineNum, commentChars, modFunc);
}
