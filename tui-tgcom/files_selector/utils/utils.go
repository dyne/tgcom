package utils

import(
	"os"
	"path/filepath"
	"io/fs"
	"fmt"
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
	if  normalizedPath == "/" {
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