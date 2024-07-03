package modelutils

import tea "github.com/charmbracelet/bubbletea"

type ModeSelector struct {
	File     string
	Choices  []string
	cursor   int
	Selected string
	Done     bool
	Speed    string
}

func NewModeSelector(choices []string, file string, speed string) ModeSelector {
	return ModeSelector{
		File:     file,
		Choices:  choices,
		Selected: "",
		Speed:    speed,
	}
}

func (m ModeSelector) Init() tea.Cmd {
	return nil
}

func (m ModeSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.Choices)-1 {
				m.cursor++
			}

		case "enter":
			m.Selected = m.Choices[m.cursor]
			m.Done = true
		}
	}
	return m, nil
}

func (m ModeSelector) View() string {
	if len(m.Choices) == 2 {
		s := Paint("silver").Render("Select 'Fast mode' if you want to toggle all your files by giving just indications about start label and end label.\nSelect 'Slow mode' if you want to specify what action to perform file by file.") + "\n"
		for i, choice := range m.Choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += cursor + " " + choice + "\n"
		}
		return s
	} else {
		s := ""
		switch m.Speed {

		case "Slow mode":
			s += Paint("silver").Render("Select action for file: "+m.File) + "\n\n"
			for i, choice := range m.Choices {
				cursor := " "
				if m.cursor == i {
					cursor = ">"
				}
				s += cursor + " " + choice + "\n"
			}

		case "Fast mode":
			s += Paint("silver").Render("Select action:") + "\n\n"
			for i, choice := range m.Choices {
				cursor := " "
				if m.cursor == i {
					cursor = ">"
				}
				s += cursor + " " + choice + "\n"
			}
		}
		return s + Paint("silver").Render("\n 'q' to quit     'enter' to modify selected files\n '↑' to go up\n '↓' to go down")
	}

}
