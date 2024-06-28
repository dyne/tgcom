package modelutils

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

func Contains(arr []string, str string) bool {
	for _, item := range arr {
		if item == str {
			return true
		}
	}
	return false
}

func Remove(slice []string, target string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if s != target {
			result = append(result, s)
		}
	}
	return result
}

func IsDirectory(path string) (bool, error) {
	// Use os.Stat to get file or directory info.
	fileInfo, err := os.Stat(path)
	if err != nil {
		// If there's an error (e.g., path doesn't exist), return the error.
		return false, err
	}

	// Check if the file mode indicates a directory.
	return fileInfo.IsDir(), nil
}

func GetParentDirectory(directoryPath string) (string, error) {
	normalizedPath := filepath.Clean(directoryPath)

	// Check for Unix-like root directory
	if normalizedPath == "/" {
		return directoryPath, nil
	}

	// Clean the directory path and get the parent directory.
	parentDir := filepath.Dir(directoryPath)

	// Check if the given path is a root directory.
	if parentDir == directoryPath {
		return "", fmt.Errorf("the given path '%s' is a root directory or invalid", directoryPath)
	}

	return parentDir, nil
}

func GetPathOfEntry(entry fs.DirEntry, baseDir string) (string, error) {
	_, err := entry.Info()
	if err != nil {
		return "", err
	}

	// Get the absolute path of the entry.
	absPath, err := filepath.Abs(filepath.Join(baseDir, entry.Name()))
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func moveToNextDir(m *FilesSelector, nextDirPath string) {
	var files_and_dir []string
	selected_files_and_dir := make(map[int]bool)

	entries, err := os.ReadDir(nextDirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		entry_Path, err := GetPathOfEntry(entry, nextDirPath)
		if err != nil {
			log.Fatal(err)
		}
		files_and_dir = append(files_and_dir, entry_Path)
	}

	for i := 0; i < len(files_and_dir); i++ {
		selected_files_and_dir[i] = false
	}

	// update values of m
	m.Current_Dir = nextDirPath
	m.Files_And_Dir = files_and_dir
	m.Selected_Files_And_Dir = selected_files_and_dir
	m.cursor = 0
	m.scrollOffset = 0
}

func moveToPrevDir(m *FilesSelector) {
	prevDirPath, err := GetParentDirectory(m.Current_Dir)
	if err != nil {
		os.Exit(0)
		log.Fatal(err)
	}

	var files_and_dir []string
	selected_files_and_dir := make(map[int]bool)

	entries, err := os.ReadDir(prevDirPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		baseDir := prevDirPath
		entry_Path, err := GetPathOfEntry(entry, baseDir)
		if err != nil {
			log.Fatal(err)
		}
		files_and_dir = append(files_and_dir, entry_Path)
	}

	for i := 0; i < len(files_and_dir); i++ {
		selected_files_and_dir[i] = false
	}

	// update values of m
	m.Current_Dir = prevDirPath
	m.Files_And_Dir = files_and_dir
	m.Selected_Files_And_Dir = selected_files_and_dir
	m.cursor = 0
	m.scrollOffset = 0
}
