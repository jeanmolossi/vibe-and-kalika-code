package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/version"
)

// replCommand pairs a slash-command with a short hint shown in the autocomplete dropdown.
type replCommand struct {
	cmd  string
	hint string
}

// replCommands is the ordered list of all REPL slash-commands used for autocomplete.
var replCommands = []replCommand{
	{cmd: "/detect", hint: "Detecta plataformas disponíveis no projeto"},
	{cmd: "/doctor", hint: "Executa verificações de saúde do ambiente"},
	{cmd: cmdHelp, hint: "Lista os comandos disponíveis"},
	{cmd: "/init", hint: "Inicializa um novo pacote interativamente"},
	{cmd: "/install", hint: "Instala um pacote a partir de um diretório local ou URL git"},
	{cmd: "/uninstall", hint: "Remove um pacote instalado anteriormente"},
	{cmd: "/update", hint: "Atualiza pacotes instalados"},
	{cmd: "/update --self", hint: "Atualiza o próprio CLI"},
	{cmd: "/validate", hint: "Valida um pacote"},
	{cmd: cmdExit, hint: "Sai do REPL"},
	{cmd: cmdQuit, hint: "Sai do REPL"},
}

const (
	maxSuggestions   = 6
	defaultWidth     = 80
	inputWidthOffset = 4
)

// Color tokens.
const (
	colourCWD          = lipgloss.Color("33")
	colourBorder       = lipgloss.Color("240")
	colourSuggestion   = lipgloss.Color("245")
	colourSelected     = lipgloss.Color("212")
	colourVersionLabel = lipgloss.Color("245")
	colourUpdate       = lipgloss.Color("208")
)

var (
	cwdStyle = lipgloss.NewStyle().
			Foreground(colourCWD).
			Bold(true)

	borderLineStyle = lipgloss.NewStyle().
			Foreground(colourBorder)

	suggestionNormalStyle = lipgloss.NewStyle().
				Foreground(colourSuggestion)

	suggestionSelectedStyle = lipgloss.NewStyle().
				Foreground(colourSelected).
				Bold(true)

	versionInfoStyle = lipgloss.NewStyle().
				Foreground(colourVersionLabel)

	updateNoticeStyle = lipgloss.NewStyle().
				Foreground(colourUpdate).
				Bold(true)
)

// updateCheckMsg carries the result of the background update check.
type updateCheckMsg struct {
	notice string
}

// cmdDoneMsg carries the error (if any) returned by a REPL command execution.
type cmdDoneMsg struct {
	err error
}

// replModel is the Bubble Tea model for the rich REPL UI.
type replModel struct {
	input        textinput.Model
	suggestions  []replCommand
	selIdx       int
	cwd          string
	updateNotice string
	width        int
	baseFactory  func() *cobra.Command
	quitting     bool
}

// newReplModel creates the initial REPL model.
func newReplModel(baseFactory func() *cobra.Command, cwd string) replModel {
	ti := textinput.New()
	ti.Placeholder = "type a command, e.g. /help"
	ti.Width = defaultWidth - inputWidthOffset
	ti.CharLimit = 512
	ti.Focus()

	return replModel{
		input:       ti,
		cwd:         cwd,
		baseFactory: baseFactory,
	}
}

// Init starts the blink ticker and launches the background update check.
func (m replModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		fetchUpdateCheck(m.cwd),
		tea.Println("Welcome to vkc interactive mode. Type /help for available commands."),
	)
}

// fetchUpdateCheck performs the update check in a background goroutine.
func fetchUpdateCheck(cwd string) tea.Cmd {
	return func() tea.Msg {
		result := app.CheckUpdates(cwd)
		if result == nil || !result.CLIUpdateAvailable {
			return updateCheckMsg{}
		}
		return updateCheckMsg{
			notice: fmt.Sprintf("New version %s available! Run /update --self", result.CLILatestVersion),
		}
	}
}

// Update handles all incoming messages.
func (m replModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.input.Width = msg.Width - inputWidthOffset
		return m, nil

	case updateCheckMsg:
		m.updateNotice = msg.notice
		return m, nil

	case cmdDoneMsg:
		if msg.err != nil {
			return m, tea.Println(fmt.Sprintf("error: %v", msg.err))
		}
		return m, nil

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	// Forward all other messages (e.g. blink timer) to the text input.
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

// handleKey dispatches keyboard events for the REPL.
func (m replModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyEsc:
		if len(m.suggestions) > 0 {
			m.suggestions = nil
			m.selIdx = 0
			return m, nil
		}
		m.quitting = true
		return m, tea.Quit

	case tea.KeyEnter:
		return m.handleEnter()

	case tea.KeyUp:
		if len(m.suggestions) > 0 {
			m.selIdx = (m.selIdx - 1 + len(m.suggestions)) % len(m.suggestions)
			return m, nil
		}

	case tea.KeyDown:
		if len(m.suggestions) > 0 {
			m.selIdx = (m.selIdx + 1) % len(m.suggestions)
			return m, nil
		}

	case tea.KeyTab:
		if len(m.suggestions) > 0 {
			m.selIdx = (m.selIdx + 1) % len(m.suggestions)
			m.input.SetValue(m.suggestions[m.selIdx].cmd)
			m.input.CursorEnd()
			return m, nil
		}

	case tea.KeyShiftTab:
		if len(m.suggestions) > 0 {
			m.selIdx = (m.selIdx - 1 + len(m.suggestions)) % len(m.suggestions)
			m.input.SetValue(m.suggestions[m.selIdx].cmd)
			m.input.CursorEnd()
			return m, nil
		}

	default:
		// All other key types are forwarded to the text input below.
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m.refreshSuggestions()
	return m, cmd
}

// handleEnter executes the current input or the highlighted suggestion.
func (m replModel) handleEnter() (tea.Model, tea.Cmd) {
	line := strings.TrimSpace(m.input.Value())
	if len(m.suggestions) > 0 {
		line = m.suggestions[m.selIdx].cmd
	}

	m.input.SetValue("")
	m.suggestions = nil
	m.selIdx = 0

	if line == "" {
		return m, nil
	}

	switch line {
	case cmdExit, cmdQuit:
		m.quitting = true
		return m, tea.Quit
	case cmdHelp:
		return m, tea.Println(replHelp)
	}

	if !strings.HasPrefix(line, "/") {
		return m, tea.Println("hint: commands start with /. Type /help for available commands.")
	}

	args := strings.Fields(strings.TrimPrefix(line, "/"))
	exec := &replExecCmd{args: args, baseFactory: m.baseFactory}
	return m, tea.Exec(exec, func(err error) tea.Msg { return cmdDoneMsg{err: err} })
}

// refreshSuggestions updates the autocomplete list from the current input value.
func (m *replModel) refreshSuggestions() {
	m.suggestions = computeSuggestions(m.input.Value())
	if m.selIdx >= len(m.suggestions) {
		m.selIdx = 0
	}
}

// View renders the full REPL UI below any previous terminal output.
func (m replModel) View() string {
	if m.quitting {
		return ""
	}

	width := m.width
	if width <= 0 {
		width = defaultWidth
	}

	var sb strings.Builder

	// Autocomplete dropdown rendered above the input.
	if len(m.suggestions) > 0 {
		limit := min(len(m.suggestions), maxSuggestions)
		for i := range limit {
			s := m.suggestions[i]
			label := s.cmd + "  — " + s.hint
			if i == m.selIdx {
				sb.WriteString(suggestionSelectedStyle.Render("▶ "+label) + "\n")
			} else {
				sb.WriteString(suggestionNormalStyle.Render("  "+label) + "\n")
			}
		}
	}

	// Current working directory immediately above the input.
	sb.WriteString(cwdStyle.Render(" "+m.cwd) + "\n")

	// Horizontal rule acting as the top border of the input area.
	border := borderLineStyle.Render(strings.Repeat("─", width))
	sb.WriteString(border + "\n")

	// The text input itself.
	sb.WriteString(" " + m.input.View() + "\n")

	// Horizontal rule acting as the bottom border of the input area.
	sb.WriteString(border + "\n")

	// Version label and optional update notice below the input.
	verLine := versionInfoStyle.Render("vkc " + version.Version)
	if m.updateNotice != "" {
		verLine += "  " + updateNoticeStyle.Render("⚠ "+m.updateNotice)
	}
	sb.WriteString(verLine)

	return sb.String()
}

// computeSuggestions returns all replCommands whose prefix matches the given input.
func computeSuggestions(input string) []replCommand {
	if input == "" {
		return nil
	}
	lower := strings.ToLower(input)
	matches := make([]replCommand, 0, len(replCommands))
	for _, c := range replCommands {
		if strings.HasPrefix(strings.ToLower(c.cmd), lower) && c.cmd != input {
			matches = append(matches, c)
		}
	}
	return matches
}

// replExecCmd implements tea.ExecCommand so that cobra sub-commands run with full
// terminal control, including nested Bubble Tea programs such as huh forms.
type replExecCmd struct {
	args        []string
	baseFactory func() *cobra.Command
	stdin       io.Reader
	stdout      io.Writer
	stderr      io.Writer
}

func (e *replExecCmd) SetStdin(r io.Reader)  { e.stdin = r }
func (e *replExecCmd) SetStdout(w io.Writer) { e.stdout = w }
func (e *replExecCmd) SetStderr(w io.Writer) { e.stderr = w }

// Run builds and executes the cobra sub-command with the terminal handed over.
func (e *replExecCmd) Run() error {
	stdin := e.stdin
	if stdin == nil {
		stdin = os.Stdin
	}
	stdout := e.stdout
	if stdout == nil {
		stdout = os.Stdout
	}
	stderr := e.stderr
	if stderr == nil {
		stderr = os.Stderr
	}

	root := e.baseFactory()
	root.SetArgs(e.args)
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SilenceErrors = true
	root.SilenceUsage = true
	return root.Execute()
}
