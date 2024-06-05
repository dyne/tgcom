package file

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func ProcessFile(filePath string, lineNum int, commentChars map[string]string, modFunc func(string, map[string]string) string) error {
	inputFile, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer inputFile.Close()

	outputPath := filePath + ".tmp"
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	_, err = inputFile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	err = writeChanges(inputFile, outputFile, lineNum, commentChars, modFunc)
	if err != nil {
		return err
	}

	err = os.Rename(outputPath, filePath)
	if err != nil {
		return err
	}

	return nil
}

func writeChanges(inputFile *os.File, outputFile *os.File, lineNum int, commentChars map[string]string, modFunc func(string, map[string]string) string) error {

	scanner := bufio.NewScanner(inputFile)
	writer := bufio.NewWriter(outputFile)
	currentLine := 1
	for scanner.Scan() {
		lineContent := scanner.Text()
		if currentLine == lineNum {
			lineContent = modFunc(lineContent, commentChars)
		}

		_, err := writer.WriteString(lineContent + "\n")
		if err != nil {
			fmt.Println("Error writing the modified file:", err)
			return err
		}

		currentLine++

	}
	if lineNum > currentLine {
		return errors.New("line number is out of range")
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return writer.Flush()
}
