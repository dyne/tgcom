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

func main(){
	// interface for the user
	
	// definition of the flags for the command line: "file", "line"
	fileFlag := flag.String("file", "", "The file to process")
    lineFlag := flag.String("line", "", "The line number or range to modify (e.g., 4 or 10-20)")
	actionFlag := flag.String("action", "", "can be comment, uncomment or toogle")
    
    // Analyize the flags
    flag.Parse()

    // take arguments of the flags
    filename := *fileFlag
	lineStr := *lineFlag
    action := *actionFlag

    var startLine, endLine int
    var err error
    var commentChars map[string]string
    var modFunc func(string, map[string]string) string

    // check if in filename there are different files
    if strings.Contains(filename, ","){
        modFunc = comment.ToggleComments
        if lineStr != "" && action != ""{
            fmt.Println("Invalid syntax")
            return
        }
        fileLine := strings.Split(filename, ",")
        for i:=0; i<len(fileLine); i++ {
            if strings.Contains(fileLine[i], ":"){
                sub := strings.Split(fileLine[i], ":")
                if len(sub) != 2 {
                    fmt.Println("Invalid syntax format. Use '<filename>:<lines>'")
                    return
                }
                
                filego := sub[0]
                lineString := sub[1]

                // extract starLine and endLine from lineString
                if strings.Contains(lineString, "-") {
                    parts := strings.Split(lineString, "-")
                    if len(parts) != 2 {
                        fmt.Println("Invalid range format. Use 'start-end'")
                        return
                    }
                    startLine, err = strconv.Atoi(parts[0])
                    if err != nil || startLine <= 0 {
                        fmt.Println("Invalid start line number")
                        return
                    }
                    endLine, err = strconv.Atoi(parts[1])
                    if err != nil || endLine < startLine {
                        fmt.Println("Invalid end line number.")
                    return
                    }
                } else {
                    startLine, err = strconv.Atoi(lineString)
                    if err != nil || startLine <= 0 {
                        fmt.Println("Please provide a valid positive integer for the line number.")
                        return
                    }
                    endLine = startLine
                }

                // store startLine and endLine in an array
                lineNum := [2]int{startLine, endLine} 

                // need to know which file you are dealing with, so extract extension of the file
                extension := filepath.Ext(filego)

                // depending on the extension we select a different map from go.go
                switch extension {
                    case ".go":
                    commentChars = language.GoCommentChars
                    /* TODO: more possible extensions
                    case ".js":
                    case ".html":
                    */
                }

                file.ProcessFile(filego, lineNum, commentChars, modFunc)

            } else {
                fmt.Println("Invalid syntax")
            }
        }
    } else { // just one filename

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
            // TODO: more possible extensions
            // case ".js":
            // case ".html":
            //
        }
    
        // depending on the action we select a different function from comment.go
        
        switch action {
            case "comment":
                modFunc = comment.Comment
            case "uncomment":
                modFunc = comment.Uncomment
            case "toggle":
                modFunc = comment.ToggleComments
        }
    
        // operate on the file
        file.ProcessFile(filename, lineNum, commentChars, modFunc)

    }
}