package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/dyne/tgcom/utils/modfile"
	"github.com/dyne/tgcom/utils/tui/modelutils"
)

// Model represents the main application model
type Model struct {
	// State variables
	State      string // Current state: "FileSelection", "ModeSelection", "ActionSelection", "LabelInput", "ApplyChanges", "Error"
	Files      []string
	Actions    []string
	Labels     []string
	LabelType  []bool
	CurrentDir string // Current directory for file selection
	Error      error

	// Models for different selection steps
	FilesSelector  modelutils.FilesSelector
	SpeedSelector  modelutils.ModeSelector
	ActionSelector modelutils.ModeSelector
	LabelInput     modelutils.LabelInput
}

// Init initializes the model
func (m Model) Init() tea.Cmd {
	return m.FilesSelector.Init()
}

// Update updates the model based on messages
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.Error != nil {
		// If an error occurred, display the error and quit
		return m, tea.Quit
	}

	switch m.State {
	case "FileSelection":
		newFilesSelector, cmd := m.FilesSelector.Update(msg)
		m.FilesSelector = newFilesSelector.(modelutils.FilesSelector)
		if m.FilesSelector.Done {
			if m.FilesSelector.Error != nil {
				m.Error = m.FilesSelector.Error
				return m, tea.Quit
			}
			m.State = "ModeSelection"
			m.SpeedSelector = modelutils.NewModeSelector([]string{"Fast mode", "Slow mode"}, "", "")
		}
		return m, cmd

	case "ModeSelection":
		newSpeedSelector, cmd := m.SpeedSelector.Update(msg)
		m.SpeedSelector = newSpeedSelector.(modelutils.ModeSelector)
		if m.SpeedSelector.Done {
			m.Files = m.FilesSelector.FilesPath
			if len(m.Files) == 0 {
				m.State = "ApplyChanges"
			} else {
				m.State = "ActionSelection"
				m.ActionSelector = modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, filepath.Base(m.Files[0]), m.SpeedSelector.Selected)
			}
		}
		return m, cmd

	case "ActionSelection":
		switch m.SpeedSelector.Selected {
		case "Slow mode":
			counter := 1
			newActionSelector, cmd := m.ActionSelector.Update(msg)
			m.ActionSelector = newActionSelector.(modelutils.ModeSelector)
			if m.ActionSelector.Done {
				m.Actions = append(m.Actions, m.ActionSelector.Selected)
				if len(m.Actions) == len(m.Files) {
					m.State = "LabelInput"
					m.LabelInput = modelutils.NewLabelInput(filepath.Base(m.Files[0]))
				} else {
					m.ActionSelector = modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, filepath.Base(m.Files[counter]), m.SpeedSelector.Selected)
					counter++
				}
			}
			return m, cmd
		case "Fast mode":
			newActionSelector, cmd := m.ActionSelector.Update(msg)
			m.ActionSelector = newActionSelector.(modelutils.ModeSelector)
			if m.ActionSelector.Done {
				for i := 0; i < len(m.Files); i++ {
					m.Actions = append(m.Actions, m.ActionSelector.Selected)
				}
				m.State = "LabelInput"
				m.LabelInput = modelutils.NewLabelInput(filepath.Base(m.Files[0]))
			}
			return m, cmd
		}

	case "LabelInput":
		switch m.SpeedSelector.Selected {
		case "Slow mode":
			counter := 1
			newLabelInput, cmd := m.LabelInput.Update(msg)
			m.LabelInput = newLabelInput.(modelutils.LabelInput)
			if m.LabelInput.Done {
				if m.LabelInput.Error != nil {
					m.Error = m.LabelInput.Error
					return m, tea.Quit
				}
				m.Labels = append(m.Labels, m.LabelInput.Input)
				m.LabelType = append(m.LabelType, m.LabelInput.IsLabel)
				if len(m.Labels) == len(m.Files) {
					m.State = "ApplyChanges"
				} else {
					m.LabelInput = modelutils.NewLabelInput(filepath.Base(m.Files[counter]))
					counter++
				}
			}
			return m, cmd
		case "Fast mode":
			newLabelInput, cmd := m.LabelInput.Update(msg)
			m.LabelInput = newLabelInput.(modelutils.LabelInput)
			m.LabelType = append(m.LabelType, m.LabelInput.IsLabel)
			if m.LabelInput.Done {
				for i := 0; i < len(m.Files); i++ {
					m.Labels = append(m.Labels, m.LabelInput.Input)
				}
				m.State = "ApplyChanges"
			}
			return m, cmd
		}

	case "ApplyChanges":
		err := m.applyChanges()
		m.Error = err
		m.State = "Final"
	case "Final":

		return m, tea.Quit // Quit the program after applying changes

	}
	return m, nil
}

// View renders the view based on the current state
func (m Model) View() string {
	if m.Error != nil {
		return fmt.Sprintf("An error occurred: %v", m.Error)
	}
	switch m.State {
	case "FileSelection":
		return m.FilesSelector.View()
	case "ModeSelection":
		return m.SpeedSelector.View()
	case "ActionSelection":
		return m.ActionSelector.View()
	case "LabelInput":
		return m.LabelInput.View()
	}
	return ""
}

// applyChanges applies changes to selected files based on user inputs
func (m *Model) applyChanges() error {
	for i := 0; i < len(m.Files); i++ {
		currentFilePath, err := AbsToRel(m.Files[i])
		if err != nil {
			return fmt.Errorf("failed to convert to relative path: %w", err)
		}
		if !m.LabelType[i] {
			conf := modfile.Config{
				Filename: currentFilePath,
				LineNum:  m.Labels[i],
				Action:   m.Actions[i],
			}
			err = modfile.ChangeFile(conf)
		} else {
			parts := strings.Split(m.Labels[i], ";")
			conf := modfile.Config{
				Filename:   currentFilePath,
				StartLabel: parts[0],
				EndLabel:   parts[1],
				Action:     m.Actions[i],
			}
			err = modfile.ChangeFile(conf)
		}

		if err != nil {
			return fmt.Errorf("failed to apply changes to file %s: %w", m.Files[i], err)
		}
	}

	return nil
}

// Helper function to convert absolute path to relative path
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
