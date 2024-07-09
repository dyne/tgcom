package modelutils

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsDirectory(t *testing.T) {
	t.Run("Existing Directory", func(t *testing.T) {
		isDir, err := IsDirectory(".")
		assert.NoError(t, err)
		assert.True(t, isDir)
	})

	t.Run("Non-Existing Path", func(t *testing.T) {
		_, err := IsDirectory("/nonexistent/path")
		assert.Error(t, err)
	})
}
func TestGetParentDirectory(t *testing.T) {
	t.Run("Valid Directory", func(t *testing.T) {
		tempDir := t.TempDir()
		subDir := filepath.Join(tempDir, "subdir")
		err := os.Mkdir(subDir, 0755)
		assert.NoError(t, err)

		parentDir, err := GetParentDirectory(subDir)
		assert.NoError(t, err)
		assert.Equal(t, tempDir, parentDir)
	})

	t.Run("Root Directory", func(t *testing.T) {
		rootDir := "/"
		if runtime.GOOS == "windows" {
			rootDir = filepath.VolumeName(rootDir) + "\\"
		}

		parentDir, err := GetParentDirectory(rootDir)
		assert.Error(t, err)
		assert.Equal(t, "", parentDir)
		assert.Equal(t, "cannot move above the root directory", err.Error())
	})

	t.Run("Parent of Root Directory", func(t *testing.T) {
		parentDir, err := GetParentDirectory("/")
		assert.Error(t, err)
		assert.Equal(t, "", parentDir)
		assert.Equal(t, "cannot move above the root directory", err.Error())
	})

	t.Run("Windows Root Directory", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			rootDir := "C:\\"
			parentDir, err := GetParentDirectory(rootDir)
			assert.Error(t, err)
			assert.Equal(t, "", parentDir)
			assert.Equal(t, "cannot move above the root directory", err.Error())
		}
	})
}

func TestGetPathOfEntry(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()

	// Create a temporary file in the temporary directory
	tempFile := filepath.Join(tempDir, "file.txt")
	_, err := os.Create(tempFile)
	assert.NoError(t, err)

	entry, err := os.ReadDir(tempDir)
	assert.NoError(t, err)
	assert.NotEmpty(t, entry)

	t.Run("Existing Entry", func(t *testing.T) {
		path, err := GetPathOfEntry(entry[0], tempDir)
		assert.NoError(t, err)
		assert.Equal(t, tempFile, path)
	})

	t.Run("Non-Existing Entry", func(t *testing.T) {
		nonexistentEntry := filepath.Join(tempDir, "nonexistent.txt")
		entry, err := os.ReadDir(nonexistentEntry)
		assert.Error(t, err)
		if len(entry) > 0 {
			_, err := GetPathOfEntry(entry[0], nonexistentEntry)
			assert.Error(t, err)
		}
	})
}
func TestMoveToNextDir(t *testing.T) {
	// Create a temporary directory structure
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	assert.NoError(t, err)
	tempFile := filepath.Join(subDir, "file.txt")
	_, err = os.Create(tempFile)
	assert.NoError(t, err)

	filesSelector := &FilesSelector{}

	t.Run("Valid Directory", func(t *testing.T) {
		err := moveToNextDir(filesSelector, subDir)
		assert.NoError(t, err)
		assert.Equal(t, subDir, filesSelector.CurrentDir)
		assert.NotEmpty(t, filesSelector.FilesAndDir)
	})

	t.Run("Non-Existing Directory", func(t *testing.T) {
		err := moveToNextDir(filesSelector, "/nonexistent/path")
		assert.Error(t, err)
	})
}
func TestMoveToPreviousDir(t *testing.T) {
	tempDir := t.TempDir()
	subDir := filepath.Join(tempDir, "subdir")
	err := os.Mkdir(subDir, 0755)
	assert.NoError(t, err)
	tempFile := filepath.Join(subDir, "file.txt")
	_, err = os.Create(tempFile)
	assert.NoError(t, err)

	filesSelector := &FilesSelector{CurrentDir: subDir}

	t.Run("Valid Directory", func(t *testing.T) {
		err := moveToPreviousDir(filesSelector)
		assert.NoError(t, err)
		assert.Equal(t, tempDir, filesSelector.CurrentDir)
		assert.NotEmpty(t, filesSelector.FilesAndDir)
	})
	t.Run("Root_Directory", func(t *testing.T) {
		// Use the root directory for the current OS as the initial directory
		rootDir := "/"
		if runtime.GOOS == "windows" {
			rootDir = filepath.VolumeName(rootDir) + "\\"
		}

		filesSelector := &FilesSelector{
			CurrentDir: rootDir,
		}

		err := moveToPreviousDir(filesSelector)
		assert.Error(t, err, "An error is expected when moving to the previous directory from the root directory")
	})
}
