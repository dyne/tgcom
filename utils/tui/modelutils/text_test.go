package modelutils

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
)

func setInput(m LabelInput, input string) LabelInput {
	m.Input = input
	return m
}
func stripANSI(str string) string {
	// Remove ANSI color codes
	re := regexp.MustCompile(`\x1b\[[0-9;]*m`)
	stripped := re.ReplaceAllString(str, "")

	// Replace consecutive whitespace characters with a single space
	stripped = regexp.MustCompile(`\s+`).ReplaceAllString(stripped, " ")

	// Trim leading and trailing whitespace
	stripped = strings.TrimSpace(stripped)

	return stripped
}

func TestNewLabelInput(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		expected LabelInput
	}{
		{
			name: "Basic Initialization",
			file: "testfile",
			expected: LabelInput{
				File:    "testfile",
				Input:   "",
				Done:    false,
				IsLabel: false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := NewLabelInput(tt.file, 10, 10)
			assert.Equal(t, tt.expected, input)
		})
	}
}

func TestInitLabelInput(t *testing.T) {
	input := NewLabelInput("testfile", 10, 10)
	cmd := input.Init()
	assert.NotNil(t, cmd)
}
func TestUpdateLabelInput(t *testing.T) {
	tests := []struct {
		name       string
		initial    LabelInput
		msg        tea.Msg
		expected   LabelInput
		cmdChecker func(tea.Cmd)
	}{
		{
			name:    "Test KeyEnter with valid input label",
			initial: setInput(NewLabelInput("", 10, 10), "test;test"),
			msg:     tea.KeyMsg{Type: tea.KeyEnter},
			expected: LabelInput{
				File:    "",
				Input:   "test;test",
				Done:    true,
				IsLabel: true,
			},
		},
		{
			name:    "Test KeyEnter with valid input lines",
			initial: setInput(NewLabelInput("", 10, 10), "1"),
			msg:     tea.KeyMsg{Type: tea.KeyEnter},
			expected: LabelInput{
				File:    "",
				Input:   "1",
				Done:    true,
				IsLabel: false,
			},
		},
		{
			name:    "Test KeyEnter with invalid input",
			initial: setInput(NewLabelInput("", 10, 10), ""),
			msg:     tea.KeyMsg{Type: tea.KeyEnter},
			expected: LabelInput{
				File:    "",
				Input:   "",
				Done:    false,
				IsLabel: false,
				Error:   fmt.Errorf("input does not match expected format (e.g., 'start';'end' or 'x-y' or single line number)"),
			},
		},
		{
			name:    "Test KeyBackspace",
			initial: NewLabelInput("", 10, 10),
			msg:     tea.KeyMsg{Type: tea.KeyBackspace},
			expected: LabelInput{
				File:    "",
				Input:   "",
				Done:    false,
				IsLabel: false,
			},
		},
		{
			name:    "Test KeyRunes",
			initial: NewLabelInput("", 10, 10),
			msg:     tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'t', 'e', 's', 't'}},
			expected: LabelInput{
				File:    "",
				Input:   "test",
				Done:    false,
				IsLabel: false,
			},
		},
		{
			name:    "Test KeyEsc",
			initial: NewLabelInput("", 10, 10),
			msg:     tea.KeyMsg{Type: tea.KeyEsc},
			expected: LabelInput{
				File:    "",
				Input:   "",
				Done:    false,
				IsLabel: false,
				Back:    true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model, cmd := tt.initial.Update(tt.msg)
			assert.Equal(t, tt.expected, model)
			if tt.cmdChecker != nil {
				tt.cmdChecker(cmd)
			}
		})
	}
}
func TestViewLabelInput(t *testing.T) {
	tests := []struct {
		name     string
		input    LabelInput
		expected string
	}{
		{
			name: "View with File",
			input: LabelInput{
				File:  "testfile",
				Input: "input",
				flash: true,
			},
			expected: "Type below the section to modify. You can insert your start label and your end label using the syntax 'start';'end' or you can modify a single line by entering the line number or a range of lines using the syntax x-y testfile: ✏ input▎ 'ctrl +c' to quit 'enter' to select the lines/labels indicated 'esc' to go back '↑' to go up '↓' to go down",
		},
		{
			name: "View without File",
			input: LabelInput{
				File:  "",
				Input: "input",
				flash: true,
			},
			expected: "Type below the section to modify. You can insert your start label and your end label using the syntax 'start';'end' or you can modify a single line by entering the line number or a range of lines using the syntax x-y ✏ input▎ 'ctrl +c' to quit 'enter' to select the lines/labels indicated 'esc' to go back '↑' to go up '↓' to go down",
		},
		{
			name: "View with Error",
			input: LabelInput{
				File:  "testfile",
				Input: "input",
				Error: fmt.Errorf("test error"),
			},
			expected: "An error occurred: test error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			view := tt.input.View()
			assert.Equal(t, tt.expected, stripANSI(view))
		})
	}

}

func TestValidateInput(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected error
		isLabel  bool
	}{
		{
			name:     "Valid Label Input",
			input:    "start;end",
			expected: nil,
			isLabel:  true,
		},
		{
			name:     "Valid Line Range Input",
			input:    "1-5",
			expected: nil,
			isLabel:  false,
		},
		{
			name:     "Valid Single Line Input",
			input:    "1",
			expected: nil,
			isLabel:  false,
		},
		{
			name:     "Invalid Input Format",
			input:    "invalid",
			expected: fmt.Errorf("input does not match expected format (e.g., 'start';'end' or 'x-y' or single line number)"),
			isLabel:  false,
		},
		{
			name:     "Invalid Line Range Format",
			input:    "1-",
			expected: fmt.Errorf("invalid input format for line range (expected 'x-y' where x and y are integers)"),
			isLabel:  false,
		},
		{
			name:     "Invalid Line Range Format 2",
			input:    "5-3",
			expected: fmt.Errorf("invalid input format for line range (start line should be less than or equal to end line)"),
			isLabel:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := LabelInput{
				Input: tt.input,
			}
			err := model.validateInput()
			assert.Equal(t, tt.expected, err)
			assert.Equal(t, tt.isLabel, model.IsLabel)
		})
	}
}
