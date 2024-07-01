package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dyne/tgcom/cmd"
	"github.com/dyne/tgcom/utils/tui"
	"github.com/dyne/tgcom/utils/tui/modelutils"
)

func main() {
	//function main must be modified, but it shows that it works both with
	//flags and as a tui
	var mod int = 0
	var num int
	if mod == 1 {
		cmd.Execute()
		_, err := fmt.Scanf("%d", &num)
		if err != nil {
			fmt.Println("Errore:", err)
		}
	}
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting current working directory: %v\n", err)
		os.Exit(1)
	}

	// Initialize your model with the current directory
	model := tui.Model{
		State:         "FileSelection",
		FilesSelector: modelutils.InitialModel(currentDir, 20),
	}

	// Bubble Tea program
	p := tea.NewProgram(model)

	// Start the program
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error starting Bubble Tea program: %v\n", err)
		os.Exit(1)
	}
}
