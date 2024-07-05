package modelutils

import (
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestFilesSelector(t *testing.T) {
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "file.txt")
	_, err := os.Create(tempFile)
	assert.NoError(t, err)
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	assert.NoError(t, err)
	tempFile = filepath.Join(subDir, "file.txt")
	_, err = os.Create(tempFile)
	assert.NoError(t, err)

	tests := []struct {
		name   string
		msg    tea.Msg
		setup  func(*FilesSelector)
		verify func(*testing.T, FilesSelector)
	}{
		{
			name: "InitialModel",
			setup: func(m *FilesSelector) {
				*m = InitialModel(tempDir, 10)
			},
			verify: func(t *testing.T, m FilesSelector) {
				assert.Equal(t, tempDir, m.CurrentDir)
				assert.Contains(t, m.FilesAndDir, subDir)
				assert.NotNil(t, m.SelectedFilesAndDir)
				assert.Equal(t, 0, m.cursor)
				assert.Equal(t, 10, m.WindowHeight)
				assert.NoError(t, m.Error)
			},
		},
		{
			name: "KeyDown",
			msg:  tea.KeyMsg{Type: tea.KeyDown},
			setup: func(m *FilesSelector) {
				*m = InitialModel(tempDir, 10)

			},
			verify: func(t *testing.T, m FilesSelector) {
				assert.Equal(t, 1, m.cursor)
			},
		},
		{
			name: "KeyUp",
			msg:  tea.KeyMsg{Type: tea.KeyUp},
			setup: func(m *FilesSelector) {
				*m = InitialModel(tempDir, 10)
				m.cursor = 1
			},
			verify: func(t *testing.T, m FilesSelector) {
				assert.Equal(t, 0, m.cursor)
			},
		},
		{
			name: "EnterDirectory",
			msg:  tea.KeyMsg{Type: tea.KeyEnter},
			setup: func(m *FilesSelector) {
				*m = InitialModel(tempDir, 10)
				m.cursor = 1
			},
			verify: func(t *testing.T, m FilesSelector) {
				assert.Equal(t, subDir, m.CurrentDir)
			},
		},
		{
			name: "SelectFile",
			msg:  tea.KeyMsg{Type: tea.KeyEnter},
			setup: func(m *FilesSelector) {
				*m = InitialModel(subDir, 10)
				m.cursor = 0
			},
			verify: func(t *testing.T, m FilesSelector) {
				assert.Contains(t, m.FilesPath, tempFile)
			},
		},
		{
			name: "Exit",
			msg:  tea.KeyMsg{Type: tea.KeyCtrlC},
			setup: func(m *FilesSelector) {
				*m = InitialModel(tempDir, 10)
			},
			verify: func(t *testing.T, m FilesSelector) {
				// Call Update with the exit message
				newModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
				model := newModel.(FilesSelector)
				msg := cmd()
				// Verify the model state after handling the exit message
				assert.Equal(t, tempDir, model.CurrentDir)
				assert.Equal(t, tea.Quit(), msg) // Verify that the command indicates quitting
			},
		},
		{
			name: "MoveToPreviousDir",
			msg:  tea.KeyMsg{Type: tea.KeyEsc},
			setup: func(m *FilesSelector) {
				*m = InitialModel(subDir, 10)
			},
			verify: func(t *testing.T, m FilesSelector) {
				assert.Equal(t, tempDir, m.CurrentDir)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := InitialModel(tempDir, 10)
			if tt.setup != nil {
				tt.setup(&model)
			}
			if tt.msg != nil {
				newModel, _ := model.Update(tt.msg)
				model = newModel.(FilesSelector)
			}
			if tt.verify != nil {
				tt.verify(t, model)
			}
		})
	}
}

func TestFilesSelectorView(t *testing.T) {
	tempDir := t.TempDir()
	tempFile1 := filepath.Join(tempDir, "file.txt")
	_, err := os.Create(tempFile1)
	assert.NoError(t, err)
	subDir := filepath.Join(tempDir, "subdir")
	err = os.Mkdir(subDir, 0755)
	assert.NoError(t, err)
	tempFile2 := filepath.Join(subDir, "file.txt")
	_, err = os.Create(tempFile2)
	assert.NoError(t, err)

	tests := []struct {
		name   string
		setup  func(*FilesSelector)
		verify func(*testing.T, string)
	}{
		{
			name: "No selected file",
			setup: func(m *FilesSelector) {
				*m = InitialModel(tempDir, 10)
			},
			verify: func(t *testing.T, view string) {
				assert.Contains(t, view, "Select the files you want to modify...")
				assert.Contains(t, view, "➪ ❒ "+tempFile1)
				assert.Contains(t, view, "❒ "+subDir)
			},
		},
		{
			name: " with  a selected file",
			setup: func(m *FilesSelector) {
				*m = InitialModel(tempDir, 10)
				m.cursor = 1
				m.FilesPath = append(m.FilesPath, tempFile1)
			},
			verify: func(t *testing.T, view string) {
				assert.Contains(t, view, "Select the files you want to modify...")
				assert.Contains(t, view, "Selected files till now: "+tempFile1)
				assert.Contains(t, view, "❒ "+tempFile1)
				assert.Contains(t, view, "➪ ❒ "+subDir)
			},
		},
		{
			name: " inside subdir",
			setup: func(m *FilesSelector) {
				*m = InitialModel(tempDir, 10)
				msg := tea.KeyMsg{Type: tea.KeyEnter}
				m.cursor = 1
				newModel, _ := m.Update(msg)
				*m = newModel.(FilesSelector)
			},
			verify: func(t *testing.T, view string) {
				assert.Contains(t, view, "Select the files you want to modify...")
				assert.Contains(t, view, "➪ ❒ "+tempFile2)
			},
		},
		{
			name: "Navigate above root directory",
			setup: func(m *FilesSelector) {
				*m = InitialModel("/", 10)
				msg := tea.KeyMsg{Type: tea.KeyEsc}
				m.cursor = 1
				newModel, _ := m.Update(msg)
				*m = newModel.(FilesSelector)
			},
			verify: func(t *testing.T, view string) {
				assert.Contains(t, view, "cannot move above the root directory")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := InitialModel(tempDir, 10)
			if tt.setup != nil {
				tt.setup(&model)
			}
			view := stripANSI(model.View())
			if tt.verify != nil {
				tt.verify(t, view)
			}
		})
	}
}
