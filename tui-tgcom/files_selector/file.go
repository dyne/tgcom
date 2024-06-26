package files_selector

import(
	"fmt"
    "os"
	"log"
	//"os/exec"
    //"runtime"
	"github.com/dyne/tgcom/tui-tgcom/files_selector/utils"
	"github.com/charmbracelet/lipgloss"
    tea "github.com/charmbracelet/bubbletea"
)

const windowHeight int = 20

type Model struct {
	Current_Dir				string
    Files_And_Dir			[]string  
    Selected_Files_And_Dir	map[int]bool
	
	Files_Path	[]string
	
    cursor 	     int                
	scrollOffset int
}

func InitialModel() Model {
	var files_and_dir []string
	selected_files_and_dir := make(map[int]bool)

	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		entry_Path, err := utils.GetPathOfEntry(entry, currentDir)
		if err != nil {
			log.Fatal(err)
		}
		files_and_dir = append(files_and_dir, entry_Path)
	}

	for i := 0; i < len(files_and_dir); i++ {
		selected_files_and_dir[i] = false
	}

	return Model{
		Current_Dir: currentDir,
		Files_And_Dir:  files_and_dir,
		Selected_Files_And_Dir: selected_files_and_dir,
	}
}

/*
func newModel(mOld Model) Model {
	var files_and_dir []string
	selected_files_and_dir := make(map[int]bool)

	currentDir := mOld.Current_Dir
	
	entries, err := os.ReadDir(currentDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		entry_Path, err := utils.GetPathOfEntry(entry, currentDir)
		if err != nil {
			log.Fatal(err)
		}
		files_and_dir = append(files_and_dir, entry_Path)
	}

	for i := 0; i < len(files_and_dir); i++ {
		selected_files_and_dir[i] = false
	}

	return Model{
		Current_Dir: currentDir,
		Files_And_Dir:  files_and_dir,
		Selected_Files_And_Dir: selected_files_and_dir,
	}
}
*/

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
				if m.cursor < m.scrollOffset {
					m.scrollOffset--
				}
			}

		// The "down" key move the cursor down
		case "down":
			if m.cursor < len(m.Files_And_Dir)-1 {
				m.cursor++
				if m.cursor >= m.scrollOffset+windowHeight {
					m.scrollOffset++
				}
			}

		// The "x" key moves you to the next directory or select/unselect a file
		case "x":
			check_dir, err := utils.IsDirectory(m.Files_And_Dir[m.cursor])
			if err != nil {
				log.Fatal(err)
			}
			
			if check_dir {
				moveToNextDir(&m, m.Files_And_Dir[m.cursor])
			} else {
				// update Files_Path
				if utils.Contains(m.Files_Path, m.Files_And_Dir[m.cursor]){
					m.Files_Path = utils.Remove(m.Files_Path, m.Files_And_Dir[m.cursor])
				} else {
					m.Files_Path = append(m.Files_Path, m.Files_And_Dir[m.cursor])
				}
				// update Selected_Files_And_Dir
				file_already_selected := m.Selected_Files_And_Dir[m.cursor]
				if file_already_selected {
					m.Selected_Files_And_Dir[m.cursor] = false
				} else {
					m.Selected_Files_And_Dir[m.cursor] = true
				}
			}

		// The "esc" key moves you to the parent directory
		case "esc":
			moveToPrevDir(&m)

		// Press "enter" when you have selected all the files to modify
		case "enter":
			m.cursor = 0
			//m.File_Flag = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	// The header
	s := paint("silver").Render("\n Select the files you want to modify...") + "\n"

	s += paint("silver").Render("\n Selected files till now:") + "\n"	
	for i := 0; i < len(m.Files_Path); i++ {	
		s += fmt.Sprintf(" %s\n", paint("green").Render(m.Files_Path[i]))
	}

	s+= "\n"

	for i := m.scrollOffset; i < m.scrollOffset+windowHeight && i < len(m.Files_And_Dir); i++ {
		choice := m.Files_And_Dir[i]
		check_dir, err := utils.IsDirectory(choice)
		if err != nil {
			log.Fatal(err)
		}
		if check_dir {
			choice = paint("blue").Render("❒ " + choice)
		} else if utils.Contains(m.Files_Path, choice) {
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
	s += paint("silver").Render("\n 'q' to quit      'esc' to move to parent directory\n '↑' to go up     'enter' to modify selected files\n '↓' to go down   'x' to select pointed file/move to pointed sub folder")
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

func main() {
	clearScreen()
    p, err := tea.NewProgram(initialModel()).Run()
    if err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
	
	// RESTART FROM NEW WINDOW
	// JUST FOR TEST
	clearScreen()
	mm := p.(Model)
	fmt.Println(mm.Current_Dir)
	fmt.Printf("Tipo di directory: %T\n", mm.Current_Dir)
	fmt.Println(mm.Files_Path)
	
    p, err = tea.NewProgram(newModel(mm)).Run()
    if err != nil {
        fmt.Printf("Alas, there's been an error: %v", err)
        os.Exit(1)
    }
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

func moveToNextDir(m *Model, nextDirPath string) {
	var files_and_dir []string
	selected_files_and_dir := make(map[int]bool)

	entries, err := os.ReadDir(nextDirPath)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, entry := range entries {
		entry_Path, err := utils.GetPathOfEntry(entry, nextDirPath)
		if err != nil {
			log.Fatal(err)
		}
		files_and_dir = append(files_and_dir, entry_Path)
	}


	for i := 0; i < len(files_and_dir); i++ {
		selected_files_and_dir[i] = false
	}

	// update values of m
	m.Current_Dir = nextDirPath
	m.Files_And_Dir = files_and_dir
	m.Selected_Files_And_Dir = selected_files_and_dir
	m.cursor = 0
	m.scrollOffset = 0
}

func moveToPrevDir(m *Model) {
	prevDirPath, err := utils.GetParentDirectory(m.Current_Dir)
	if err != nil {
		os.Exit(0)
		log.Fatal(err)
	}

	var files_and_dir []string
	selected_files_and_dir := make(map[int]bool)

	entries, err := os.ReadDir(prevDirPath)
	if err != nil {
		log.Fatal(err)
	}
	
	for _, entry := range entries {
		baseDir := prevDirPath
		entry_Path, err := utils.GetPathOfEntry(entry, baseDir)
		if err != nil {
			log.Fatal(err)
		}
		files_and_dir = append(files_and_dir, entry_Path)
	}

	for i := 0; i < len(files_and_dir); i++ {
		selected_files_and_dir[i] = false
	}
	
	// update values of m
	m.Current_Dir = prevDirPath
	m.Files_And_Dir = files_and_dir
	m.Selected_Files_And_Dir = selected_files_and_dir
	m.cursor = 0
	m.scrollOffset = 0
}

// come incapsulare file_selector in un'altra tui:
// scrivi un'altra tui presupponendo di partire da un elenco di files [file_1, file_2, file_3, ...]
// determina il nuovo Modello che ti interessa e aggiungi questo elenco di files al Modello nuovo.