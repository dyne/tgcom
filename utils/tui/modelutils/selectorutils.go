package modelutils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

func Contains(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func Remove(slice []string, target string) []string {
	result := make([]string, 0, len(slice))
	for _, item := range slice {
		if item != target {
			result = append(result, item)
		}
	}
	return result
}

func IsDirectory(path string) (bool, error) {
	if path == "/" {
		return true, nil
	}
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

func GetParentDirectory(directoryPath string) (string, error) {

	parentDir := filepath.Dir(directoryPath)
	if parentDir == directoryPath || parentDir == "/" {
		return "", fmt.Errorf("cannot move above the root directory")
	}

	return parentDir, nil
}

func GetPathOfEntry(entry fs.DirEntry, baseDir string) (string, error) {
	_, err := entry.Info()
	if err != nil {
		return "", err
	}

	absPath, err := filepath.Abs(filepath.Join(baseDir, entry.Name()))
	if err != nil {
		return "", err
	}

	return absPath, nil
}

func moveToNextDir(filesSelector *FilesSelector, nextDirPath string) error {
	var filesAndDirs []string
	selectedFilesAndDirs := make(map[int]bool)

	entries, err := os.ReadDir(nextDirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath, err := GetPathOfEntry(entry, nextDirPath)
		if err != nil {
			return err
		}
		filesAndDirs = append(filesAndDirs, entryPath)
	}

	for i := 0; i < len(filesAndDirs); i++ {
		selectedFilesAndDirs[i] = false
	}

	filesSelector.CurrentDir = nextDirPath
	filesSelector.FilesAndDir = filesAndDirs
	filesSelector.SelectedFilesAndDir = selectedFilesAndDirs
	filesSelector.cursor = 0
	filesSelector.scrollOffset = 0
	return nil
}

func moveToPreviousDir(filesSelector *FilesSelector) error {
	prevDirPath, err := GetParentDirectory(filesSelector.CurrentDir)
	if err != nil {
		return err
	}

	var filesAndDirs []string
	selectedFilesAndDirs := make(map[int]bool)

	entries, err := os.ReadDir(prevDirPath)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		entryPath, err := GetPathOfEntry(entry, prevDirPath)
		if err != nil {
			return err
		}
		filesAndDirs = append(filesAndDirs, entryPath)
	}

	for i := 0; i < len(filesAndDirs); i++ {
		selectedFilesAndDirs[i] = false
	}

	filesSelector.CurrentDir = prevDirPath
	filesSelector.FilesAndDir = filesAndDirs
	filesSelector.SelectedFilesAndDir = selectedFilesAndDirs
	filesSelector.cursor = 0
	filesSelector.scrollOffset = 0
	return nil
}
