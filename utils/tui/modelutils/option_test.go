package modelutils

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func TestNewModeSelector(t *testing.T) {
	tests := []struct {
		name     string
		choices  []string
		file     string
		speed    string
		expected ModeSelector
	}{
		{
			name:    "Basic Initialization",
			choices: []string{"Option1", "Option2"},
			file:    "",
			speed:   "",
			expected: ModeSelector{
				File:     "",
				Choices:  []string{"Option1", "Option2"},
				Selected: "",
				Speed:    "",
				Done:     false,
				Back:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			selector := NewModeSelector(tt.choices, tt.file, tt.speed)
			assert.Equal(t, tt.expected.File, selector.File)
			assert.Equal(t, tt.expected.Choices, selector.Choices)
			assert.Equal(t, tt.expected.Selected, selector.Selected)
			assert.Equal(t, tt.expected.Speed, selector.Speed)
			assert.False(t, selector.Done)
			assert.False(t, selector.Back)
		})
	}
}

func TestInit(t *testing.T) {
	selector := NewModeSelector([]string{"Option1", "Option2"}, "", "")
	cmd := selector.Init()
	assert.Nil(t, cmd)
}

func TestUpdate(t *testing.T) {
	tests := []struct {
		name       string
		initial    ModeSelector
		msg        tea.Msg
		expected   ModeSelector
		cmdChecker func(tea.Cmd)
	}{
		{
			name:    "Test up key",
			initial: NewModeSelector([]string{"Option1", "Option2"}, "", ""),
			msg:     tea.KeyMsg{Type: tea.KeyUp},
			expected: ModeSelector{
				Choices: []string{"Option1", "Option2"},
				cursor:  0,
			},
		},
		{
			name:    "Test down key",
			initial: NewModeSelector([]string{"Option1", "Option2"}, "", ""),
			msg:     tea.KeyMsg{Type: tea.KeyDown},
			expected: ModeSelector{
				Choices: []string{"Option1", "Option2"},
				cursor:  1,
			},
		},
		{
			name:    "Test enter key",
			initial: NewModeSelector([]string{"Option1", "Option2"}, "", ""),
			msg:     tea.KeyMsg{Type: tea.KeyEnter},
			expected: ModeSelector{
				Choices:  []string{"Option1", "Option2"},
				cursor:   0,
				Selected: "Option1",
				Done:     true,
			},
		},
		{
			name:    "Test esc key",
			initial: NewModeSelector([]string{"Option1", "Option2"}, "", ""),
			msg:     tea.KeyMsg{Type: tea.KeyEsc},
			expected: ModeSelector{
				Choices: []string{"Option1", "Option2"},
				cursor:  0,
				Back:    true,
			},
		},
		{
			name:    "Test quit keys",
			initial: NewModeSelector([]string{"Option1", "Option2"}, "", ""),
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'q'}, Alt: false},
			expected: ModeSelector{
				Choices: []string{"Option1", "Option2"},
				cursor:  0,
			},
			cmdChecker: func(cmd tea.Cmd) {
				if cmd != nil {
					msg := cmd()
					assert.Equal(t, tea.Quit(), msg)
				} else {
					t.Error("Expected tea.Quit command, got nil")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, cmd := tt.initial.Update(tt.msg)
			selector := model.(ModeSelector) // type assertion
			assert.Equal(t, tt.expected, selector)
			if tt.cmdChecker != nil {
				tt.cmdChecker(cmd)
			}
		})
	}
}

func TestView(t *testing.T) {
	tests := []struct {
		name     string
		selector ModeSelector
		expected string
	}{
		{
			name:     "View with cursor at default position",
			selector: NewModeSelector([]string{"Option1", "Option2"}, "testfile", "Fast mode"),
			expected: "Select 'Fast mode' if you want to toggle all your files by giving just indications about start label and end label. Select 'Slow mode' if you want to specify what action to perform file by file. > Option1 Option2",
		},
		{
			name: "View with cursor at position 1",
			selector: ModeSelector{
				File:    "",
				Choices: []string{"Option1", "Option2"},
				cursor:  1,
				Speed:   "Fast mode",
			},
			expected: "Select 'Fast mode' if you want to toggle all your files by giving just indications about start label and end label. Select 'Slow mode' if you want to specify what action to perform file by file. Option1 > Option2",
		},
		{
			name: "View with more options in Slow mode",
			selector: ModeSelector{
				File:    "testfile",
				Choices: []string{"Option1", "Option2", "Option3"},
				Speed:   "Slow mode",
			},
			expected: "Select action for file: testfile > Option1 Option2 Option3 'q' to quit 'enter' to modify selected files 'esc' to go back '↑' to go up '↓' to go down",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.selector.View()
			assert.Equal(t, tt.expected, stripANSI(view))
		})
	}
}
