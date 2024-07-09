package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dyne/tgcom/utils/tui/modelutils"
	"github.com/stretchr/testify/assert"
)

func TestModel(t *testing.T) {
	t.Run("Init", func(t *testing.T) {
		model := Model{FilesSelector: modelutils.InitialModel(".", 10)}
		cmd := model.Init()
		assert.Nil(t, cmd)
	})

	t.Run("Update", func(t *testing.T) {
		type updateTest struct {
			name   string
			model  Model
			setup  func(m *Model)
			msg    tea.Msg
			verify func(*testing.T, Model)
		}

		tests := []updateTest{
			{
				name: "FileSelection to ModeSelection",
				model: Model{
					State:         "FileSelection",
					FilesSelector: modelutils.InitialModel(".", 10),
				},
				setup: func(m *Model) {
					m.FilesSelector.FilesPath = []string{"path/test/file1", "path/test/file2"}

				},
				msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
				verify: func(t *testing.T, m Model) {
					assert.True(t, m.FilesSelector.Done)
					assert.Contains(t, m.Files, "path/test/file1")
					assert.Contains(t, m.Files, "path/test/file2")
					assert.Equal(t, "ModeSelection", m.State)

				},
			},
			{
				name: "FileSelection to ActionSelection",
				model: Model{
					State:         "FileSelection",
					FilesSelector: modelutils.InitialModel(".", 10),
				},
				setup: func(m *Model) {
					m.FilesSelector.FilesPath = []string{"path/test/file1"}

				},
				msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
				verify: func(t *testing.T, m Model) {
					assert.True(t, m.FilesSelector.Done)
					assert.Contains(t, m.Files, "path/test/file1")
					assert.Equal(t, "ActionSelection", m.State)

				},
			},
			{
				name: "No file selected",
				model: Model{
					State:         "FileSelection",
					FilesSelector: modelutils.InitialModel(".", 10),
				},
				msg: tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'x'}},
				verify: func(t *testing.T, m Model) {
					assert.False(t, m.FilesSelector.Done)
					assert.Equal(t, "FileSelection", m.State)

				},
			},
			{
				name: "ModeSelection to ActionSelection",
				model: Model{
					State:         "ModeSelection",
					SpeedSelector: modelutils.NewModeSelector([]string{"Fast mode", "Slow mode"}, "", ""),
					Files:         []string{"file1.txt", "file2.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEnter},
				verify: func(t *testing.T, m Model) {
					assert.True(t, m.SpeedSelector.Done)
					assert.Equal(t, "ActionSelection", m.State)

				},
			},
			{
				name: "ActionSelection to ActionSelection",
				model: Model{
					State:          "ActionSelection",
					SpeedSelector:  modelutils.ModeSelector{Selected: "Slow mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", "Slow mode"),
					Files:          []string{"file1.txt", "file2.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEnter},
				verify: func(t *testing.T, m Model) {
					assert.Equal(t, "ActionSelection", m.State)
					assert.Equal(t, 1, len(m.Actions))
				},
			},
			{
				name: "ActionSelection to LabelInput fast",
				model: Model{
					State:          "ActionSelection",
					SpeedSelector:  modelutils.ModeSelector{Selected: "Fast mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", "Fast mode"),
					Files:          []string{"file1.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEnter},
				verify: func(t *testing.T, m Model) {
					assert.True(t, m.ActionSelector.Done)
					assert.Equal(t, "LabelInput", m.State)
				},
			},
			{
				name: "ActionSelection to LabelInput slow",
				model: Model{
					State:          "ActionSelection",
					SpeedSelector:  modelutils.ModeSelector{Selected: "Slow mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", "Slow mode"),
					Files:          []string{"file1.txt", "file2.txt"},
					Actions:        []string{"test"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEnter},
				verify: func(t *testing.T, m Model) {
					assert.True(t, m.ActionSelector.Done)
					assert.Equal(t, "LabelInput", m.State)
					assert.Equal(t, 2, len(m.Actions))
				},
			},
			{
				name: "LabelInput to LabelInput",
				model: Model{
					State:         "LabelInput",
					SpeedSelector: modelutils.ModeSelector{Selected: "Slow mode"},
					LabelInput:    modelutils.LabelInput{Input: "1-3"},
					Files:         []string{"file1.txt", "file2.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEnter},
				verify: func(t *testing.T, m Model) {
					assert.Equal(t, "LabelInput", m.State)
					assert.Equal(t, 1, len(m.Labels))
					assert.Equal(t, 1, len(m.LabelType))
					assert.False(t, m.LabelType[0])
				},
			},
			{
				name: "ModeSelection to FileSelection",
				model: Model{
					State:         "ModeSelection",
					FilesSelector: modelutils.InitialModel(".", 10),
					SpeedSelector: modelutils.NewModeSelector([]string{"Fast mode", "Slow mode"}, "", ""),
					Files:         []string{"file1.txt", "file2.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEsc},
				verify: func(t *testing.T, m Model) {
					assert.False(t, m.FilesSelector.Done)
					assert.Contains(t, m.Files, "file1.txt")
					assert.Contains(t, m.Files, "file2.txt")
					assert.Equal(t, "FileSelection", m.State)
					assert.Equal(t, "", m.SpeedSelector.Selected)

				},
			},
			{
				name: "ActionSelection to FileSelection",
				model: Model{
					State:          "ActionSelection",
					FilesSelector:  modelutils.InitialModel(".", 10),
					SpeedSelector:  modelutils.ModeSelector{Selected: "Fast mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", ""),
					Files:          []string{"file1.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEsc},
				verify: func(t *testing.T, m Model) {
					assert.False(t, m.FilesSelector.Done)
					assert.Equal(t, "FileSelection", m.State)
					assert.Equal(t, "", m.ActionSelector.Selected)

				},
			},
			{
				name: "ActionSelection to ModeSelection fast",
				model: Model{
					State:          "ActionSelection",
					SpeedSelector:  modelutils.ModeSelector{Selected: "Fast mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", ""),
					Files:          []string{"file1.txt", "file2.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEsc},
				verify: func(t *testing.T, m Model) {
					assert.False(t, m.SpeedSelector.Done)
					assert.Equal(t, m.State, "ModeSelection")
					assert.Equal(t, "", m.ActionSelector.Selected)

				},
			},
			{
				name: "ActionSelection to ModeSelection slow",
				model: Model{
					State:          "ActionSelection",
					SpeedSelector:  modelutils.ModeSelector{Selected: "Slow mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", ""),
					Files:          []string{"file1.txt", "file2.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEsc},
				verify: func(t *testing.T, m Model) {
					assert.False(t, m.SpeedSelector.Done)
					assert.Equal(t, "ModeSelection", m.State)
					assert.Equal(t, "", m.ActionSelector.Selected)
				},
			},
			{
				name: "ActionSelection to ActionSelection back",
				model: Model{
					State:          "ActionSelection",
					SpeedSelector:  modelutils.ModeSelector{Selected: "Slow mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", "Slow mode"),
					Files:          []string{"file1.txt", "file2.txt"},
					Actions:        []string{"test", "comment"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEsc},
				verify: func(t *testing.T, m Model) {
					assert.False(t, m.ActionSelector.Done)
					assert.Equal(t, "ActionSelection", m.State)
					assert.Equal(t, 1, len(m.Actions))
					assert.Equal(t, "test", m.Actions[0])
				},
			},
			{
				name: "LabelInput to ActionSelection fast",
				model: Model{
					State:          "LabelInput",
					SpeedSelector:  modelutils.ModeSelector{Selected: "Fast mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", "Fast mode"),
					LabelInput:     modelutils.LabelInput{Input: "1-3"},
					Files:          []string{"file1.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEsc},
				verify: func(t *testing.T, m Model) {
					assert.False(t, m.ActionSelector.Done)
					assert.Equal(t, "ActionSelection", m.State)
				},
			},
			{
				name: "LabelInput to ActionSelection slow",
				model: Model{
					State:          "LabelInput",
					SpeedSelector:  modelutils.ModeSelector{Selected: "Slow mode"},
					ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", "Slow mode"),
					LabelInput:     modelutils.LabelInput{Input: "1-3"},
					Files:          []string{"file1.txt", "file2.txt"},
					Actions:        []string{"test", "comment"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEsc},
				verify: func(t *testing.T, m Model) {
					assert.False(t, m.ActionSelector.Done)
					assert.Equal(t, "ActionSelection", m.State)
					assert.Equal(t, []string{"test"}, m.Actions)

				},
			},
			{
				name: "LabelInput to LabelInput back",
				model: Model{
					State:         "LabelInput",
					SpeedSelector: modelutils.ModeSelector{Selected: "Slow mode"},
					LabelInput:    modelutils.LabelInput{Input: "start;end"},
					Files:         []string{"file1.txt", "file2.txt"},
					Labels:        []string{"1-3"},
					LabelType:     []bool{false},
				},
				msg: tea.KeyMsg{Type: tea.KeyEsc},
				verify: func(t *testing.T, m Model) {
					assert.Equal(t, m.State, "LabelInput")
					assert.Equal(t, []string{}, m.Labels)
					assert.Equal(t, []bool{}, m.LabelType)
				},
			},

			{
				name: "LabelInput to ApplyChanges fast",
				model: Model{
					State:         "LabelInput",
					SpeedSelector: modelutils.ModeSelector{Selected: "Fast mode"},
					LabelInput:    modelutils.LabelInput{Input: "1-3"},
					Files:         []string{"file1.txt"},
				},
				msg: tea.KeyMsg{Type: tea.KeyEnter},
				verify: func(t *testing.T, m Model) {
					assert.True(t, m.LabelInput.Done)
					assert.Equal(t, m.State, "ApplyChanges")
					assert.False(t, m.LabelType[0])

				},
			},
			{
				name: "LabelInput to ApplyChanges slow",
				model: Model{
					State:         "LabelInput",
					SpeedSelector: modelutils.ModeSelector{Selected: "Slow mode"},
					LabelInput:    modelutils.LabelInput{Input: "start;end"},
					Files:         []string{"file1.txt", "file2.txt"},
					Labels:        []string{"1-3"},
					LabelType:     []bool{false},
				},
				msg: tea.KeyMsg{Type: tea.KeyEnter},
				verify: func(t *testing.T, m Model) {
					assert.True(t, m.LabelInput.Done)
					assert.Equal(t, m.State, "ApplyChanges")
					assert.Equal(t, 2, len(m.Labels))
					assert.Equal(t, 2, len(m.LabelType))
					assert.True(t, m.LabelType[1])

				},
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if tt.setup != nil {
					tt.setup(&tt.model)
				}
				if tt.msg != nil {
					newModel, _ := tt.model.Update(tt.msg)
					tt.model = newModel.(Model)
				}

				if tt.verify != nil {
					tt.verify(t, tt.model)
				}
			})
		}
	})
	t.Run("applyChanges", func(t *testing.T) {

		type applyChangesTest struct {
			name          string
			model         Model
			expected      []string
			expectedError error
		}

		tests := []applyChangesTest{
			{
				name: "Successful Changes",
				model: Model{
					Actions:   []string{"comment"},
					Labels:    []string{"start;end"},
					LabelType: []bool{true},
				},
				expected: []string{"start\n// Line 1\n// Line 2\n// Line 3\nend\nLine 4\n"},
			},
			{
				name: "Multiple Files",
				model: Model{
					Actions:   []string{"comment", "toggle"},
					Labels:    []string{"start;end", "1-3"},
					LabelType: []bool{true, false},
				},
				expected: []string{"start\n// Line 1\n// Line 2\n// Line 3\nend\nLine 4\n", "// start\n// Line 1\n// Line 2\nLine 3\nend\nLine 4\n"},
			},
			{
				name: "Error Applying Changes",
				model: Model{
					Files:     []string{"error.txt"},
					Actions:   []string{"comment"},
					Labels:    []string{"label"},
					LabelType: []bool{false},
				},
				expectedError: fmt.Errorf("ailed to convert to relative path"),
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				for i := 0; i < len(tt.model.Actions); i++ {
					filename := fmt.Sprintf("file%d.go", i)
					tmpFile, cleanup := createTempFile(t, "start\nLine 1\nLine 2\nLine 3\nend\nLine 4", filename)
					defer cleanup()
					tt.model.Files = append(tt.model.Files, tmpFile.Name())
				}
				cmd := tt.model.applyChanges()
				msg := cmd()
				err := msg.(applyChangesMsg).err
				actual := []string{}

				if tt.expectedError == nil {
					for i := 0; i < len(tt.model.Files); i++ {
						content, err := os.ReadFile(tt.model.Files[i])
						if err != nil {
							t.Fatalf("Failed to read temporary file: %v", err)
						}
						actual = append(actual, string(content))
					}
					assert.Nil(t, err)

					assert.Equal(t, tt.expected, actual)
				}
				if tt.expectedError != nil {
					assert.Contains(t, err.Error(), tt.expectedError.Error())
				} else {
					assert.Nil(t, err)
				}
			})
		}
	})

	t.Run("View", func(t *testing.T) {
		type viewTest struct {
			name     string
			model    Model
			expected string
		}

		tests := []viewTest{
			{
				name:     "FileSelection View",
				model:    Model{State: "FileSelection", FilesSelector: modelutils.InitialModel(".", 10)},
				expected: "Select the files you want to modify",
			},
			{
				name:     "ModeSelection View",
				model:    Model{State: "ModeSelection", SpeedSelector: modelutils.NewModeSelector([]string{"Fast mode", "Slow mode"}, "", "")},
				expected: "Select 'Fast mode'",
			},
			{
				name:     "ActionSelection View",
				model:    Model{State: "ActionSelection", ActionSelector: modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, "", "Fast mode")},
				expected: "Select action",
			},
			{
				name:     "LabelInput View",
				model:    Model{State: "LabelInput", LabelInput: modelutils.NewLabelInput("")},
				expected: "Type below the section to modify",
			},

			{
				name:     "Final View with Error",
				model:    Model{State: "Final", Error: fmt.Errorf("test error")},
				expected: "An error occurred: test error",
			},
			{
				name:     "Final View without Error",
				model:    Model{State: "Final"},
				expected: "Changes applied successfully!",
			},
		}

		for _, test := range tests {
			t.Run(test.name, func(t *testing.T) {
				view := test.model.View()
				assert.Contains(t, view, test.expected)
			})
		}
	})

}

func createTempFile(t testing.TB, content string, name string) (*os.File, func()) {
	t.Helper()
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, name)
	temp, err := os.Create(tempFile)
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	if _, err := temp.WriteString(content); err != nil {
		t.Fatalf("Failed to write to temporary file: %v", err)
	}
	if err := temp.Close(); err != nil {
		t.Fatalf("Failed to close temporary file: %v", err)
	}
	return temp, func() { os.Remove(temp.Name()) }
}
