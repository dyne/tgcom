package modelutils

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type FilesSelector struct {
	Current_Dir            string
	Files_And_Dir          []string
	Selected_Files_And_Dir map[int]bool
	Files_Path             []string
	cursor                 int
	scrollOffset           int
	Done                   bool
	WindowHeight           int
}

func InitialModel(currentDir string, windowHeight int) FilesSelector {
	var files_and_dir []string
	selected_files_and_dir := make(map[int]bool)

	entries, err := os.ReadDir(currentDir)
	if err != nil {
		log.Printf("Error occurred: %v", err)
	}

	for _, entry := range entries {
		entry_Path, err := GetPathOfEntry(entry, currentDir)
		if err != nil {
			log.Printf("Error occurred: %v", err)
		}
		files_and_dir = append(files_and_dir, entry_Path)
	}

	for i := 0; i < len(files_and_dir); i++ {
		selected_files_and_dir[i] = false
	}

	return FilesSelector{
		Current_Dir:            currentDir,
		Files_And_Dir:          files_and_dir,
		Selected_Files_And_Dir: selected_files_and_dir,
		WindowHeight:           windowHeight,
	}
}

func (m FilesSelector) Init() tea.Cmd {
	return nil
}

func (m FilesSelector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit

		// The "up" key move the cursor up
		case "up":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scrollOffset {
					m.scrollOffset--
				}
			}

		// The "down" key move the cursor down
		case "down":
			if m.cursor < len(m.Files_And_Dir)-1 {
				m.cursor++
				if m.cursor >= m.scrollOffset+m.WindowHeight {
					m.scrollOffset++
				}
			}

		// The "enter" key moves you to the next directory or select/unselect a file
		case "enter":
			check_dir, err := IsDirectory(m.Files_And_Dir[m.cursor])
			if err != nil {
				log.Fatal(err)
			}

			if check_dir {
				moveToNextDir(&m, m.Files_And_Dir[m.cursor])
			} else {
				// update Files_Path
				if Contains(m.Files_Path, m.Files_And_Dir[m.cursor]) {
					m.Files_Path = Remove(m.Files_Path, m.Files_And_Dir[m.cursor])
				} else {
					m.Files_Path = append(m.Files_Path, m.Files_And_Dir[m.cursor])
				}
				m.Selected_Files_And_Dir[m.cursor] = !m.Selected_Files_And_Dir[m.cursor]
			}
		case "esc":
			moveToPrevDir(&m)
		// Press x to confirm
		case "x":
			m.Done = true
		}
	}
	return m, nil
}

func (m FilesSelector) View() string {
	// The header
	s := paint("silver").Render("\n Select the files you want to modify...") + "\n"

	s += paint("silver").Render("\n Selected files till now:") + "\n"
	for i := 0; i < len(m.Files_Path); i++ {
		s += fmt.Sprintf(" %s\n", paint("green").Render(m.Files_Path[i]))
	}

	s += "\n"

	for i := m.scrollOffset; i < m.scrollOffset+m.WindowHeight && i < len(m.Files_And_Dir); i++ {
		choice := m.Files_And_Dir[i]
		check_dir, err := IsDirectory(choice)
		if err != nil {
			log.Fatal(err)
		}
		if check_dir {
			choice = paint("blue").Render("❒ " + choice)
		} else if Contains(m.Files_Path, choice) {
			choice = paint("lime").Render("❒ " + choice)
		} else {
			choice = paint("silver").Render("❒ " + choice)
		}

		cursor := " "
		if m.cursor == i {
			cursor = paint("red").Render(" ➪")
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	// The footer
	s += paint("silver").Render("\n 'q' to quit      'esc' to move to parent directory\n '↑' to go up     'x' to modify selected files\n '↓' to go down   'enter' to select pointed file/move to pointed sub folder")
	return s
}

func paint(color string) lipgloss.Style {
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
