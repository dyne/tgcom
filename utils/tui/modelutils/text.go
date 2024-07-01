package modelutils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type LabelInput struct {
	File    string
	Input   string
	Done    bool
	IsLabel bool // Added to distinguish between labels and line numbers
	flash   bool
	Error   error
}

func NewLabelInput(File string) LabelInput {
	return LabelInput{
		File:    File,
		Input:   "",
		Done:    false,
		IsLabel: false,
	}
}

func (m LabelInput) Init() tea.Cmd {
	return StartTicker()
}

func (m LabelInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
			if err := m.validateInput(); err != nil {
				m.Error = err
				return m, nil
			}
			m.Done = true

		case "backspace":
			if len(m.Input) > 0 {
				m.Input = m.Input[:len(m.Input)-1]
			}

		default:
			m.Input += msg.String()
		}
	case tickMsg:
		m.flash = !m.flash
	}
	return m, nil
}

func (m LabelInput) View() string {
	flash := ""
	if m.flash {
		flash = paint("green").Render("▎")
	}

	s := paint("silver").Render("Type below the section to modify. You can insert your start label\nand your end label using the syntax 'start';'end' or you can modify\n a single line by entering the line number or a range of lines using the syntax x-y") + "\n\n"
	if m.File != "" {
		s += paint("green").Render(m.File+": ✏ "+m.Input) + flash + "\n"
	} else {
		s += paint("green").Render("✏ "+m.Input) + flash + "\n"
	}

	if m.Error != nil {
		s += paint("red").Render("\nError: "+m.Error.Error()) + "\n"
	}

	s += paint("silver").Render("\n 'q' to quit     'enter' to select the lines/labels indicated\n '↑' to go up\n '↓' to go down")
	return s
}

func StartTicker() tea.Cmd {
	return tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

func (m *LabelInput) validateInput() error {
	// Trim spaces from input
	input := strings.TrimSpace(m.Input)

	// Check if input contains ';', indicating labels
	if strings.Contains(input, ";") {
		m.IsLabel = true
		return nil
	}

	// Check if input contains '-', indicating a range of lines
	if strings.Contains(input, "-") {
		parts := strings.Split(input, "-")
		if len(parts) != 2 {
			return fmt.Errorf("invalid input format for line range (expected 'x-y')")
		}

		// Check if both parts are integers
		startLine, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
		endLine, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))

		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid input format for line range (expected 'x-y' where x and y are integers)")
		}

		// Check if start line is less than or equal to end line
		if startLine > endLine {
			return fmt.Errorf("invalid input format for line range (start line should be less than or equal to end line)")
		}

		m.IsLabel = false
		return nil
	}

	// Check if input is a single line number
	_, err := strconv.Atoi(input)
	if err != nil {
		return fmt.Errorf("input does not match expected format (e.g., 'start';'end' or 'x-y' or single line number)")
	}

	m.IsLabel = false
	return nil
}