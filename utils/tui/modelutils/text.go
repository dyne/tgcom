package modelutils

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type LabelInput struct {
	File  string
	Input string
	Done  bool
	flash bool
}

func NewLabelInput(File string) LabelInput {
	return LabelInput{
		File:  File,
		Input: "",
		Done:  false,
	}
}

func (m LabelInput) Init() tea.Cmd {
	return nil
}

func (m LabelInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "enter":
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
	// flash command
	flash := ""
	if m.flash {
		flash = paint("green").Render("▎")
	}

	// the header
	s := paint("silver").Render("Type below the section to modify. You can insert your start label\nand your end label using the syntax 'start';'end' or you can modify\n a single line by digiting the line number or a range of lines using the syntax x-y") + "\n\n"
	if m.File != "" {
		s += paint("green").Render(m.File+": ✏ "+m.Input) + flash + "\n"
	} else {
		s += paint("green").Render("✏ "+m.Input) + flash + "\n"
	}

	// The footer
	s += paint("silver").Render("\n 'q' to quit     'enter' to select the lines/labels indicated\n '↑' to go up\n '↓' to go down")
	return s
}
func StartTicker() tea.Cmd {
	return tea.Tick(time.Millisecond*2500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time
