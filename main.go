package main

import (
	"github.com/dyne/tgcom/cmd"
	"os/exec"
	"os"
	"runtime"
	"fmt"
	"path/filepath"
	"github.com/charmbracelet/bubbletea"
	"github.com/dyne/tgcom/tui-tgcom/options_selector"
	"github.com/dyne/tgcom/tui-tgcom/files_selector"
	"github.com/dyne/tgcom/tui-tgcom/texts_selector"
	"github.com/dyne/tgcom/utils/modfile" // o solo modfile
)

func main() {
	//function main must be modified, but it shows that it works both with
	//flags and as a tui
	var mod int = 0
	var num int
	if mod == 0 {
		cmd.Execute()
		_, err := fmt.Scanf("%d", &num)
    	if err != nil {
        	fmt.Println("Errore:", err)
    	} 
	}

	clearScreen() // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

    // initialize model for file selection
	model1 := file_selector.InitialModel()

    // select files
	p1, _ := tea.NewProgram(model1).Run()
	model1 = p1.(file_selector.Model)

    // Files []string contain the path for all the files user selects
	Files := model1.Files_Path // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

	// Ask if user wants fast or slow mode
    clearScreen() // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

    optionsz := []string{"Fast mode", "Slow mode"}
    
    model2 := option_selector.Model{
        Options: optionsz,
    }

    // Esegui Init
    /*
    cmd2 := model2.Init()
    if cmd2 == nil {
        fmt.Println("Init command not executed")
    }
    */

    p2, _ := tea.NewProgram(model2).Run()
    model2 = p2.(option_selector.Model)

    // Speed is the string "Fast mode o Slow mode"
    speed := model2.Selected

    // Array that contains informations about how the user wants to modify each file
    var Actions []string
    var Labels []string

    switch speed {
    case "Fast mode":
        // Ask the user for the labels he wants to assign
        clearScreen() // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
        model3 := text_selector.Model{}

        // Esegui Init
        cmd3 := model3.Init()
        if cmd3 == nil {
            fmt.Println("Init command not executed")
        }

	    p3, _ := tea.NewProgram(model3).Run()

	    model3 = p3.(text_selector.Model)

        // update of Actions and Label
        for i := 0; i < len(Files); i++{
            Actions = append(Actions, "toggle")
            Labels = append(Labels, model3.Input)
        }
		
		/*
        fmt.Println("Files:")
        fmt.Println(Files)
        fmt.Println("Azioni:")
        fmt.Println(Actions)
        fmt.Println("Labels:")
        fmt.Println(Labels)
		*/

    case "Slow mode":
        var p4 tea.Model
        for i := 0; i < len(Files); i++ {
            // Ask the user the action to perform
            clearScreen() // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

            optionsz = []string{"toggle", "comment", "uncomment"}
    
            model4 := option_selector.Model{
                Options: optionsz,
            }

            // Esegui Init
            /*
            cmd4 := model4.Init()
            if cmd4 == nil {
                fmt.Println("Init command not executed")
            }
            */

            p4, _ = tea.NewProgram(model4).Run()
            model4 = p4.(option_selector.Model)

            // Speed is the string "Fast mode o Slow mode"
            Actions = append(Actions, model4.Selected) // <<<<<<<<<<<<<<<<<<<<<

            // Ask the user for the lines/labels
            clearScreen() // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<
            model5 := text_selector.Model{}

            // Esegui Init
            cmd5 := model5.Init()
            if cmd5 == nil {
                fmt.Println("Init command not executed")
            }

	        p5, _ := tea.NewProgram(model5).Run()

	        model5 = p5.(text_selector.Model)
            // Label contiene le stringhe di start e end
	        Labels = append(Labels, model5.Input) // <<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<

			/*
            fmt.Println("Files:")
            fmt.Println(Files)
            fmt.Println("Azioni:")
            fmt.Println(Actions)
            fmt.Println("Labels:")
            fmt.Println(Labels)
			*/
        }

    }

    // ora modifica i files
    var conf modfile.Config
    for i := 0; i < len(Files); i++ {
		currentFilePath, err := AbsToRel(Files[i])
		if err != nil{
			os.Exit(1)
		}
        conf = modfile.Config{Filename: currentFilePath, LineNum: Labels[i], Action: Actions[i]}
        fmt.Println(conf)

        err = modfile.ChangeFile(conf)
        if err != nil {
            os.Exit(0)
        }
    }
}

// clean the screen
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

// extract relative path from absolute path
// try to adapt better
func AbsToRel(absPath string) (string, error) {
    // Get the current working directory
    currentDir, err := filepath.Abs(".")
    if err != nil {
        return "", fmt.Errorf("failed to get the current directory: %w", err)
    }

    // Convert the absolute path to a relative path
    relPath, err := filepath.Rel(currentDir, absPath)
    if err != nil {
        return "", fmt.Errorf("failed to convert to relative path: %w", err)
    }

    return relPath, nil
}
