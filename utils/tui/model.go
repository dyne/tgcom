package tui

import (
	"fmt"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/dyne/tgcom/utils/modfile"
	"github.com/dyne/tgcom/utils/tui/modelutils"
)

// Model represents the main application model
type Model struct {
	// State variables
	State      string // Current state: "FileSelection", "ModeSelection", "ActionSelection", "LabelInput", "ApplyChanges", "Error", "Final"
	Files      []string
	Actions    []string
	Labels     []string
	LabelType  []bool
	CurrentDir string // Current directory for file selection
	Error      error
	Width      int
	Height     int

	// Models for different selection steps
	FilesSelector  modelutils.FilesSelector
	SpeedSelector  modelutils.ModeSelector
	ActionSelector modelutils.ModeSelector
	LabelInput     modelutils.LabelInput
}

// applyChangesMsg represents a message indicating that changes have been applied.
type applyChangesMsg struct {
	err error
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

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.State == "Final" {
			return m, tea.Quit
		}
	case applyChangesMsg:
		if msg.err != nil {
			m.Error = msg.err
		}
		m.State = "Final"
		return m, nil
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
			m.Files = m.FilesSelector.FilesPath
			if len(m.Files) == 1 {
				m.SpeedSelector = modelutils.ModeSelector{
					File:     m.Files[0],
					Choices:  []string{"Fast mode", "Slow mode"},
					Selected: "Fast mode",
					Speed:    "",
				}
				m.State = "ActionSelection"
				m.ActionSelector = modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, filepath.Base(m.Files[0]), m.SpeedSelector.Selected, m.Width, m.Height)
			} else {
				m.State = "ModeSelection"
				m.SpeedSelector = modelutils.NewModeSelector([]string{"Fast mode", "Slow mode"}, "", "", m.Width, m.Height)
			}
		}
		return m, cmd

	case "ModeSelection":
		newSpeedSelector, cmd := m.SpeedSelector.Update(msg)
		m.SpeedSelector = newSpeedSelector.(modelutils.ModeSelector)
		if m.SpeedSelector.Back {
			m.State = "FileSelection"
			m.FilesSelector.Done = false
		}
		if m.SpeedSelector.Done {
			m.State = "ActionSelection"
			m.ActionSelector = modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, filepath.Base(m.Files[0]), m.SpeedSelector.Selected, m.Width, m.Height)
		}
		return m, cmd

	case "ActionSelection":
		switch m.SpeedSelector.Selected {
		case "Slow mode":
			newActionSelector, cmd := m.ActionSelector.Update(msg)
			m.ActionSelector = newActionSelector.(modelutils.ModeSelector)
			if m.ActionSelector.Back {
				if len(m.Actions) == 0 {
					m.SpeedSelector.Done = false
					m.SpeedSelector.Selected = ""
					m.State = "ModeSelection"
				} else {
					m.ActionSelector.Done = false
					m.Actions = m.Actions[:len(m.Actions)-1]
					m.State = "ActionSelection"
					m.ActionSelector = modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, filepath.Base(m.Files[len(m.Actions)]), m.SpeedSelector.Selected, m.Width, m.Height)
				}
			}
			if m.ActionSelector.Done {
				m.Actions = append(m.Actions, m.ActionSelector.Selected)
				if len(m.Actions) == len(m.Files) {
					m.State = "LabelInput"
					m.LabelInput = modelutils.NewLabelInput(filepath.Base(m.Files[0]), m.Width, m.Height)
				} else {
					m.ActionSelector = modelutils.NewModeSelector([]string{"toggle", "comment", "uncomment"}, filepath.Base(m.Files[len(m.Actions)]), m.SpeedSelector.Selected, m.Width, m.Height)
				}
			}
			return m, cmd
		case "Fast mode":
			newActionSelector, cmd := m.ActionSelector.Update(msg)
			m.ActionSelector = newActionSelector.(modelutils.ModeSelector)
			if m.ActionSelector.Back {
				if len(m.Files) == 1 {
					if !m.FilesSelector.MultipleSelection {
						m.FilesSelector.FilesPath = []string{}
					}
					m.State = "FileSelection"
					m.FilesSelector.Done = false
				} else {
					m.SpeedSelector.Done = false
					m.SpeedSelector.Selected = ""
					m.State = "ModeSelection"
				}
			}
			if m.ActionSelector.Done {
				for i := 0; i < len(m.Files); i++ {
					m.Actions = append(m.Actions, m.ActionSelector.Selected)
				}
				m.State = "LabelInput"
				m.LabelInput = modelutils.NewLabelInput("", m.Width, m.Height)
			}
			return m, cmd
		}

	case "LabelInput":
		switch m.SpeedSelector.Selected {
		case "Slow mode":
			newLabelInput, cmd := m.LabelInput.Update(msg)
			m.LabelInput = newLabelInput.(modelutils.LabelInput)
			if m.LabelInput.Back {
				if len(m.Labels) == 0 {
					m.ActionSelector.Done = false
					m.ActionSelector.Selected = ""
					m.Actions = m.Actions[:len(m.Actions)-1]
					m.State = "ActionSelection"
				} else {
					m.LabelInput.Done = false
					m.Labels = m.Labels[:len(m.Labels)-1]
					m.LabelType = m.LabelType[:len(m.LabelType)-1]
					m.State = "LabelInput"
					m.LabelInput = modelutils.NewLabelInput(filepath.Base(m.Files[len(m.Labels)]), m.Width, m.Height)
				}
			}
			if m.LabelInput.Done {
				if m.LabelInput.Error != nil {
					m.Error = m.LabelInput.Error
					return m, tea.Quit
				}
				m.Labels = append(m.Labels, m.LabelInput.Input)
				m.LabelType = append(m.LabelType, m.LabelInput.IsLabel)
				if len(m.Labels) == len(m.Files) {
					m.State = "ApplyChanges"
					return m, m.applyChanges()
				} else {
					m.LabelInput = modelutils.NewLabelInput(filepath.Base(m.Files[len(m.Labels)]), m.Width, m.Height)
				}
			}
			return m, cmd
		case "Fast mode":
			newLabelInput, cmd := m.LabelInput.Update(msg)
			m.LabelInput = newLabelInput.(modelutils.LabelInput)
			if m.LabelInput.Back {
				m.ActionSelector.Done = false
				m.ActionSelector.Selected = ""
				m.Actions = nil
				m.State = "ActionSelection"
			}
			if m.LabelInput.Done {
				for i := 0; i < len(m.Files); i++ {
					m.Labels = append(m.Labels, m.LabelInput.Input)
					m.LabelType = append(m.LabelType, m.LabelInput.IsLabel)
				}
				m.State = "ApplyChanges"
				return m, m.applyChanges()
			}
			return m, cmd
		}

	case "ApplyChanges":
		return m, m.applyChanges()

	case "Final":
		return m, nil
	}

	return m, nil
}

// View renders the view based on the current state
func (m Model) View() string {
	var rightPane string
	fileSelectionPane := m.FilesSelector.View()
	var halfWidth int
	switch m.State {
	case "FileSelection":
		return fileSelectionPane
	case "ModeSelection":
		halfWidth = m.SpeedSelector.Width
		rightPane = m.SpeedSelector.View()
	case "ActionSelection":
		halfWidth = m.ActionSelector.Width
		rightPane = m.ActionSelector.View()
	case "LabelInput":
		halfWidth = m.LabelInput.Width
		rightPane = m.LabelInput.View()
	case "ApplyChanges":
		rightPane = "Applying changes..."
	case "Final":
		if m.Error != nil {
			rightPane = modelutils.Paint("red").Render(fmt.Sprintf("An error occurred: %v\nPress any key to exit.", m.Error))
		} else {
			rightPane = "Changes applied successfully!\nPress any key to exit."
		}
	}

	// Use a style for the layout
	layout := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Width(halfWidth).Render(fileSelectionPane),
		lipgloss.NewStyle().Width(halfWidth).Render(rightPane),
	)

	return layout
}

// applyChanges applies changes to selected files based on user inputs
func (m *Model) applyChanges() tea.Cmd {
	return func() tea.Msg {
		for i := 0; i < len(m.Files); i++ {
			currentFilePath, err := AbsToRel(m.Files[i])
			if err != nil {
				return applyChangesMsg{err: fmt.Errorf("failed to convert to relative path: %w", err)}
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
				return applyChangesMsg{err: fmt.Errorf("failed to apply changes to file %s: %w", m.Files[i], err)}
			}
		}

		return applyChangesMsg{err: nil}
	}
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
