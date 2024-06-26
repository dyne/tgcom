package options_selector

import(
	"fmt"
	//"os/exec"
	//"runtime"
	//"os"
	"github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cursor   int
	Options  []string
	Selected string
}

func (m Model) Init() tea.Cmd {
    return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			}

		// The "down" key move the cursor down
		case "down":
			if m.cursor < len(m.Options)-1 {
				m.cursor++
			}
		// Press "enter" when you have Selected all the files to modify
		case "enter":
			m.Selected = m.Options[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	// the header
	s := paint("silver").Render("Choose among the following Options...") + "\n\n"
	for i := 0; i < len(m.Options); i++ {
		cursor := " " 
		if m.cursor == i {
			cursor = paint("red").Render(" ➪")
		}

		s += fmt.Sprintf("%s %s\n", cursor, paint("silver").Render(m.Options[i]))
	}
	return s
}

/*
func clearScreen() {
    var cmd *exec.Cmd
    switch runtime.GOOS {
    case "windows":
        cmd = exec.Command("cmd", "/c", "cls")
    default:
        cmd = exec.Command("clear")
    }
    cmd.Stdout = os.Stdout
    cmd.Run()
}

// JUST FOR TEST
func main() {
	clearScreen()
	example := []string{"Opzione 1", "Opzione 2", "Opzione 3", "Opzione 4"}
    p, err := tea.NewProgram(Model{Options: example}).Run()
    if err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
	
	fmt.Printf("Tipo di directory: %T\n", p)

	// JUST FOR TEST
	mm := p.(Model)
	fmt.Println(mm.Selected)
	fmt.Printf("Tipo di directory: %T\n", mm.Selected)
}
*/

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