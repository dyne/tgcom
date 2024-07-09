package modelutils

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FilesSelector struct {
	CurrentDir          string
	FilesAndDir         []string
	SelectedFilesAndDir map[int]bool
	FilesPath           []string
	cursor              int
	scrollOffset        int
	Done                bool
	WindowHeight        int
	Error               error
	NoFileSelected      bool
}

func InitialModel(currentDir string, windowHeight int) FilesSelector {
	var filesAndDir []string
	selectedFilesAndDir := make(map[int]bool)

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		return FilesSelector{Error: fmt.Errorf("error reading directory: %w", err)}
	}

	for _, entry := range entries {
		entryPath, err := GetPathOfEntry(entry, currentDir)
		if err != nil {
			return FilesSelector{Error: fmt.Errorf("error getting path of entry: %w", err)}
		}
		filesAndDir = append(filesAndDir, entryPath)
	}

	for i := 0; i < len(filesAndDir); i++ {
		selectedFilesAndDir[i] = false
	}

	return FilesSelector{
		CurrentDir:          currentDir,
		FilesAndDir:         filesAndDir,
		SelectedFilesAndDir: selectedFilesAndDir,
		WindowHeight:        windowHeight,
	}
}

func (m FilesSelector) Init() tea.Cmd {
	return nil
}

func (m FilesSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.Error != nil {
			return m, tea.Quit
		}
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scrollOffset {
					m.scrollOffset--
				}
			}
		case "down":
			if m.cursor < len(m.FilesAndDir)-1 {
				m.cursor++
				if m.cursor >= m.scrollOffset+m.WindowHeight {
					m.scrollOffset++
				}
			}
		case "enter":
			m.NoFileSelected = false
			checkDir, err := IsDirectory(m.FilesAndDir[m.cursor])
			if err != nil {
				m.Error = fmt.Errorf("error checking directory: %w", err)
				return m, tea.Quit
			}
			if checkDir {
				err := moveToNextDir(&m, m.FilesAndDir[m.cursor])
				if err != nil {
					m.Error = fmt.Errorf("error checking directory: %w", err)
					return m, tea.Quit
				}
			} else {
				if Contains(m.FilesPath, m.FilesAndDir[m.cursor]) {
					m.FilesPath = Remove(m.FilesPath, m.FilesAndDir[m.cursor])
				} else {
					m.FilesPath = append(m.FilesPath, m.FilesAndDir[m.cursor])
				}
				m.SelectedFilesAndDir[m.cursor] = !m.SelectedFilesAndDir[m.cursor]
			}
		case "esc":
			err := moveToPreviousDir(&m)
			if err != nil {
				m.Error = fmt.Errorf("error moving back: %w", err)
				return m, tea.Quit
			}
		case "x":
			if len(m.FilesPath) == 0 {
				m.NoFileSelected = true
			} else {
				m.Done = true
			}
		}
	}
	return m, nil
}

func (m FilesSelector) View() string {
	if m.Error != nil {
		return Paint("red").Render(fmt.Sprintf("An error occurred: %v", m.Error))
	}

	s := Paint("silver").Render("\n Select the files you want to modify...") + "\n"
	s += Paint("silver").Render("\n Selected files till now:") + "\n"
	if m.NoFileSelected {
		s += Paint("red").Render("\n No file selected. Please select at least one file or quit.") + "\n"
	}
	for i := 0; i < len(m.FilesPath); i++ {
		s += fmt.Sprintf(" %s\n", Paint("green").Render(m.FilesPath[i]))
	}
	s += "\n"

	for i := m.scrollOffset; i < m.scrollOffset+m.WindowHeight && i < len(m.FilesAndDir); i++ {
		choice := m.FilesAndDir[i]
		checkDir, err := IsDirectory(choice)
		if err != nil {
			m.Error = fmt.Errorf("error checking directory: %w", err)
			return Paint("red").Render(fmt.Sprintf("An error occurred: %v", m.Error))
		}
		if checkDir {
			choice = Paint("blue").Render("❒ " + choice)
		} else if Contains(m.FilesPath, choice) {
			choice = Paint("lime").Render("❒ " + choice)
		} else {
			choice = Paint("silver").Render("❒ " + choice)
		}

		cursor := " "
		if m.cursor == i {
			cursor = Paint("red").Render(" ➪")
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	s += Paint("silver").Render("\n 'q' to quit      'esc' to move to parent directory\n '↑' to go up     'x' to modify selected files\n '↓' to go down   'enter' to select pointed file/move to pointed sub folder")
	return s
}

func Paint(color string) lipgloss.Style {
	switch color {
	case "lime":
		lime := lipgloss.Color("#00FF00")
		return lipgloss.NewStyle().Foreground(lime)
	case "blue":
		blue := lipgloss.Color("#0000FF")
		return lipgloss.NewStyle().Foreground(blue)
	case "green":
		green := lipgloss.Color("#008000")
		return lipgloss.NewStyle().Foreground(green)
	case "red":
		red := lipgloss.Color("#FF0000")
		return lipgloss.NewStyle().Foreground(red)
	case "silver":
		silver := lipgloss.Color("#C0C0C0")
		return lipgloss.NewStyle().Foreground(silver)
	default:
		white := lipgloss.Color("#FFFFFF")
		return lipgloss.NewStyle().Foreground(white)
	}
}
