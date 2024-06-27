package text_selector

import(
	//"fmt"
	//"os/exec"
	//"runtime"
	//"os"
	"time"
	"github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Input string
	flash bool

	Header string
	Help   string
}

func (m Model) Init() tea.Cmd {
    return tea.Batch(StartTicker())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		// Press "enter" ...
		case "enter":
			return m, tea.Quit
		case "backspace":
			if len(m.Input) > 0 {
				m.Input = m.Input[:len(m.Input) - 1]
			}
		default:
			m.Input += msg.String()
		}
	case tickMsg:
		m.flash= !m.flash
	}
	return m, nil
}

func (m Model) View() string {
	// flash command
	flash := ""
	if m.flash {
		flash = paint("green").Render("▎")
	}

	// the header
	s := m.Header
	
	s += paint("green").Render(" ✏ " + m.Input) +  flash + "\n"
	
	// The footer
	s += m.Help
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
    p, err := tea.NewProgram(Model{}).Run()
    if err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
	
	// JUST FOR TEST
	mm := p.(Model)
	fmt.Println(mm.Input)
	fmt.Printf("Tipo di directory: %T\n", mm.Input)
}
*/

func StartTicker() tea.Cmd {
	return tea.Tick(time.Millisecond*2500, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

type tickMsg time.Time

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