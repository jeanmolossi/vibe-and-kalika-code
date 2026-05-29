package ui

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const maxDropdownSuggestions = 8

// pathInputModel is a bubbletea model for an interactive path input with a
// live dropdown suggestion list.
type pathInputModel struct {
	title       string
	input       textinput.Model
	suggestions []string
	cursor      int
	confirmed   bool
	canceled    bool
}

func newPathInputModel(title, placeholder string) pathInputModel {
	ti := textinput.New()
	ti.Placeholder = placeholder
	ti.Focus()
	ti.CharLimit = 512

	return pathInputModel{
		title:       title,
		input:       ti,
		suggestions: pathSuggestions(""),
		cursor:      -1,
	}
}

// Init implements tea.Model.
func (m pathInputModel) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model.
func (m pathInputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		var cmd tea.Cmd
		m.input, cmd = m.input.Update(msg)
		return m, cmd
	}
	return m.handleKey(keyMsg)
}

func (m pathInputModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC, tea.KeyEsc:
		m.canceled = true
		return m, tea.Quit
	case tea.KeyEnter:
		return m.handleEnter()
	case tea.KeyDown:
		if len(m.suggestions) > 0 {
			m.cursor = min(m.cursor+1, len(m.suggestions)-1)
		}
		return m, nil
	case tea.KeyUp:
		if m.cursor > -1 {
			m.cursor--
		}
		return m, nil
	case tea.KeyTab:
		return m.handleTab()
	default:
		return m.handleTyping(msg)
	}
}

func (m pathInputModel) handleEnter() (tea.Model, tea.Cmd) {
	if m.cursor >= 0 && m.cursor < len(m.suggestions) {
		m.input.SetValue(m.suggestions[m.cursor])
	}
	m.confirmed = true
	return m, tea.Quit
}

func (m pathInputModel) handleTab() (tea.Model, tea.Cmd) {
	if m.cursor >= 0 && m.cursor < len(m.suggestions) {
		m.input.SetValue(m.suggestions[m.cursor] + "/")
		m.cursor = -1
		m.suggestions = pathSuggestions(m.input.Value())
		return m, nil
	}
	if len(m.suggestions) == 1 {
		m.input.SetValue(m.suggestions[0] + "/")
		m.cursor = -1
		m.suggestions = pathSuggestions(m.input.Value())
	}
	return m, nil
}

func (m pathInputModel) handleTyping(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.suggestions = pathSuggestions(m.input.Value())
	m.cursor = -1
	return m, cmd
}

// View implements tea.Model.
func (m pathInputModel) View() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("212"))
	normalStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("245"))
	selectedStyle := lipgloss.NewStyle().Background(lipgloss.Color("62")).Bold(true)

	var b strings.Builder
	b.WriteString(titleStyle.Render("? " + m.title))
	b.WriteString("\n> ")
	b.WriteString(m.input.View())

	shown := m.suggestions
	if len(shown) > maxDropdownSuggestions {
		shown = shown[:maxDropdownSuggestions]
	}

	if len(shown) > 0 {
		b.WriteString("\n")
		for i, s := range shown {
			b.WriteString("\n")
			if i == m.cursor {
				b.WriteString(selectedStyle.Render("  > " + s))
			} else {
				b.WriteString(normalStyle.Render("    " + s))
			}
		}
	}

	b.WriteString("\n")
	return b.String()
}

// RunPathInput runs an interactive path input prompt with a live directory
// suggestion dropdown. Returns the confirmed path or an error if the user
// cancels (Esc or Ctrl+C).
func RunPathInput(title, placeholder string) (string, error) {
	m := newPathInputModel(title, placeholder)
	p := tea.NewProgram(m, tea.WithOutput(os.Stderr))

	result, err := p.Run()
	if err != nil {
		return "", fmt.Errorf("path input: %w", err)
	}

	final, ok := result.(pathInputModel)
	if !ok {
		return "", fmt.Errorf("path input: unexpected model type")
	}

	if final.canceled {
		return "", fmt.Errorf("canceled")
	}

	return strings.TrimSpace(final.input.Value()), nil
}
