package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"

	"github.com/jeanmolossi/vibe-and-kalika-code/internal/app"
	"github.com/jeanmolossi/vibe-and-kalika-code/internal/state"
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

// interactiveCommands lists subcommands that require direct TTY access (e.g. huh forms).
// These are executed via tea.Exec so the terminal is temporarily handed over.
var interactiveCommands = map[string]struct{}{
	"install":   {},
	"init":      {},
	"uninstall": {}, // may show a huh confirmation form for outside-root files
}

// Layout constants.
const (
	maxSuggestions = 6
	defaultWidth   = 80
	defaultHeight  = 24
	// headerLines counts the lines rendered by renderHeader:
	// 6 logo + 1 blank + 1 title + 1 subtitle + 1 divider = 10.
	headerLines = 10
	// inputAreaLines counts the lines in the input section:
	// cwd + top-border + input + bottom-border + version = 5.
	inputAreaLines   = 5
	inputWidthOffset = 4
	minViewportLines = 2
	// inputAndVpCmdsCount is the fixed capacity for the sub-model command slice
	// (text input + viewport = 2 commands per update cycle).
	inputAndVpCmdsCount = 2
)

// asciiLogo is the VKC banner displayed at the top of the fullscreen REPL.
const asciiLogo = `
██╗   ██╗██╗  ██╗ ██████╗
██║   ██║██║ ██╔╝██╔════╝
██║   ██║█████╔╝ ██║     
╚██╗ ██╔╝██╔═██╗ ██║     
 ╚████╔╝ ██║  ██╗╚██████╗
  ╚═══╝  ╚═╝  ╚═╝ ╚═════╝`

// Color palette.
const (
	colourCWD          = lipgloss.Color("33")
	colourBorder       = lipgloss.Color("240")
	colourSuggestion   = lipgloss.Color("245")
	colourSelected     = lipgloss.Color("212")
	colourVersionLabel = lipgloss.Color("245")
	colourUpdate       = lipgloss.Color("208")
	colourLogo         = lipgloss.Color("99")
	colourTitle        = lipgloss.Color("212")
	colourSubtitle     = lipgloss.Color("245")
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

	logoStyle = lipgloss.NewStyle().
			Foreground(colourLogo).
			Bold(true)

	titleStyle = lipgloss.NewStyle().
			Foreground(colourTitle).
			Bold(true)

	subtitleStyle = lipgloss.NewStyle().
			Foreground(colourSubtitle)
)

// Message types ---------------------------------------------------------------

// updateCheckMsg carries the result of the background update check.
type updateCheckMsg struct{ notice string }

// cmdDoneMsg carries the result of an interactive (tea.Exec) command.
type cmdDoneMsg struct{ err error }

// cmdOutputMsg carries the captured output of a non-interactive command.
type cmdOutputMsg struct {
	output string
	err    error
}

// installedPackagesMsg carries the list of installed package names.
type installedPackagesMsg struct{ packages []string }

// replModel -------------------------------------------------------------------

// replModel is the Bubble Tea model for the fullscreen REPL UI.
type replModel struct {
	input             textinput.Model
	vp                viewport.Model
	suggestions       []replCommand
	selIdx            int
	cwd               string
	updateNotice      string
	width             int
	height            int
	baseFactory       func() *cobra.Command
	quitting          bool
	installedPackages []string
	vpLines           []string
	ready             bool
}

// newReplModel creates the initial REPL model with a welcome message pre-loaded.
func newReplModel(baseFactory func() *cobra.Command, cwd string) replModel {
	ti := textinput.New()
	ti.Placeholder = "ex: /help, /install, /uninstall"
	ti.Width = defaultWidth - inputWidthOffset
	ti.CharLimit = 512
	ti.Focus()

	return replModel{
		input:       ti,
		cwd:         cwd,
		baseFactory: baseFactory,
		width:       defaultWidth,
		height:      defaultHeight,
		vpLines:     []string{"Bem-vindo ao vkc. Digite /help para ver os comandos disponíveis."},
	}
}

// Init launches background tasks: cursor blink, update check, package list.
func (m replModel) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		fetchUpdateCheck(m.cwd),
		fetchInstalledPackages(),
	)
}

// fetchUpdateCheck performs the CLI/package update check in a background goroutine.
func fetchUpdateCheck(cwd string) tea.Cmd {
	return func() tea.Msg {
		result := app.CheckUpdates(cwd)
		if result == nil {
			return updateCheckMsg{}
		}
		var parts []string
		if result.CLIUpdateAvailable {
			parts = append(parts, fmt.Sprintf("Nova versão %s disponível! Execute /update --self", result.CLILatestVersion))
		}
		for _, pkg := range result.PackageUpdates {
			parts = append(parts, fmt.Sprintf("Pacote %s: %s → %s disponível", pkg.Name, pkg.CurrentVersion, pkg.LatestVersion))
		}
		return updateCheckMsg{notice: strings.Join(parts, " | ")}
	}
}

// fetchInstalledPackages loads the installed package names from the local state file.
func fetchInstalledPackages() tea.Cmd {
	return func() tea.Msg {
		st, err := state.Read()
		if err != nil {
			return installedPackagesMsg{}
		}
		names := make([]string, 0, len(st.Installations))
		for _, inst := range st.Installations {
			names = append(names, inst.Package)
		}
		return installedPackagesMsg{packages: names}
	}
}

// execCaptured runs a non-interactive cobra command, capturing its output.
func execCaptured(args []string, baseFactory func() *cobra.Command) tea.Cmd {
	return func() tea.Msg {
		var buf bytes.Buffer
		root := baseFactory()
		root.SetArgs(args)
		root.SetOut(&buf)
		root.SetErr(&buf)
		root.SilenceErrors = true
		root.SilenceUsage = true
		err := root.Execute()
		return cmdOutputMsg{output: buf.String(), err: err}
	}
}

// Update ----------------------------------------------------------------------

// Update handles all incoming Bubble Tea messages.
func (m replModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.input.Width = msg.Width - inputWidthOffset
		vpH := m.viewportHeight()
		if !m.ready {
			m.vp = viewport.New(msg.Width, vpH)
			m.vp.SetContent(strings.Join(m.vpLines, "\n"))
			m.vp.GotoBottom()
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = vpH
		}
		return m, nil

	case updateCheckMsg:
		m.updateNotice = msg.notice
		return m, nil

	case installedPackagesMsg:
		m.installedPackages = msg.packages
		return m, nil

	case cmdOutputMsg:
		output := strings.TrimRight(msg.output, "\n")
		if output != "" {
			m = m.appendToViewport(output)
		}
		if msg.err != nil {
			m = m.appendToViewport("erro: " + msg.err.Error())
		}
		return m, fetchInstalledPackages()

	case cmdDoneMsg:
		if msg.err != nil {
			m = m.appendToViewport("erro: " + msg.err.Error())
		}
		return m, fetchInstalledPackages()

	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	var cmds []tea.Cmd
	cmds = make([]tea.Cmd, 0, inputAndVpCmdsCount)
	var inputCmd tea.Cmd
	m.input, inputCmd = m.input.Update(msg)
	cmds = append(cmds, inputCmd)

	var vpCmd tea.Cmd
	m.vp, vpCmd = m.vp.Update(msg)
	cmds = append(cmds, vpCmd)

	return m, tea.Batch(cmds...)
}

// handleKey dispatches keyboard events.
func (m replModel) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.Type {
	case tea.KeyCtrlC:
		m.quitting = true
		return m, tea.Quit

	case tea.KeyEsc:
		if len(m.suggestions) > 0 {
			m.suggestions = nil
			m.selIdx = 0
			m = m.syncViewportHeight()
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
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd

	case tea.KeyDown:
		if len(m.suggestions) > 0 {
			m.selIdx = (m.selIdx + 1) % len(m.suggestions)
			return m, nil
		}
		var cmd tea.Cmd
		m.vp, cmd = m.vp.Update(msg)
		return m, cmd

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
		// All other keys are forwarded to the text input and suggestions are refreshed.
	}

	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	m = m.refreshSuggestions()
	return m, cmd
}

// handleEnter executes the current input or highlighted suggestion.
func (m replModel) handleEnter() (tea.Model, tea.Cmd) {
	line := strings.TrimSpace(m.input.Value())
	if len(m.suggestions) > 0 {
		line = m.suggestions[m.selIdx].cmd
	}

	m.input.SetValue("")
	m.suggestions = nil
	m.selIdx = 0
	m = m.syncViewportHeight()

	if line == "" {
		return m, nil
	}

	switch line {
	case cmdExit, cmdQuit:
		m.quitting = true
		return m, tea.Quit
	case cmdHelp:
		m = m.appendToViewport(replHelp)
		return m, nil
	}

	if !strings.HasPrefix(line, "/") {
		m = m.appendToViewport("dica: comandos começam com /. Digite /help para ver os disponíveis.")
		return m, nil
	}

	args := strings.Fields(strings.TrimPrefix(line, "/"))
	if len(args) == 0 {
		return m, nil
	}

	m = m.appendToViewport("\n❯ " + line)

	if _, isInteractive := interactiveCommands[args[0]]; isInteractive {
		exec := &replExecCmd{args: args, baseFactory: m.baseFactory}
		return m, tea.Exec(exec, func(err error) tea.Msg { return cmdDoneMsg{err: err} })
	}

	return m, execCaptured(args, m.baseFactory)
}

// appendToViewport adds text to the scrollable history and scrolls to the bottom.
func (m replModel) appendToViewport(s string) replModel {
	m.vpLines = append(m.vpLines, strings.Split(s, "\n")...)
	if m.ready {
		m.vp.SetContent(strings.Join(m.vpLines, "\n"))
		m.vp.GotoBottom()
	}
	return m
}

// refreshSuggestions recomputes the autocomplete list from the current input.
func (m replModel) refreshSuggestions() replModel {
	m.suggestions = computeSuggestions(m.input.Value(), m.installedPackages)
	if m.selIdx >= len(m.suggestions) {
		m.selIdx = 0
	}
	return m.syncViewportHeight()
}

// syncViewportHeight adjusts the viewport height based on current suggestion count.
func (m replModel) syncViewportHeight() replModel {
	if m.ready {
		m.vp.Height = m.viewportHeight()
	}
	return m
}

// viewportHeight returns the number of lines available for the scrollable viewport.
func (m replModel) viewportHeight() int {
	suggCount := min(len(m.suggestions), maxSuggestions)
	h := m.height - headerLines - 1 - suggCount - inputAreaLines
	if h < minViewportLines {
		return minViewportLines
	}
	return h
}

// View ------------------------------------------------------------------------

// View renders the fullscreen REPL layout.
func (m replModel) View() string {
	if m.quitting {
		return ""
	}

	width := m.width
	if width <= 0 {
		width = defaultWidth
	}

	var sb strings.Builder

	sb.WriteString(renderHeader(width))

	if m.ready {
		sb.WriteString(m.vp.View())
		sb.WriteByte('\n')
	} else {
		// Before WindowSizeMsg, fill the viewport region with blank lines so
		// the input area stays anchored to the bottom of the alt screen.
		vpH := m.height - headerLines - inputAreaLines
		if vpH < 0 {
			vpH = 0
		}
		sb.WriteString(strings.Repeat("\n", vpH))
	}

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

	border := borderLineStyle.Render(strings.Repeat("─", width))
	sb.WriteString(cwdStyle.Render(" "+m.cwd) + "\n")
	sb.WriteString(border + "\n")
	sb.WriteString(" " + m.input.View() + "\n")
	sb.WriteString(border + "\n")

	verLine := versionInfoStyle.Render("vkc " + version.Version)
	if m.updateNotice != "" {
		verLine += "  " + updateNoticeStyle.Render("⚠ "+m.updateNotice)
	}
	sb.WriteString(verLine)

	return sb.String()
}

// renderHeader renders the ASCII art banner and application description.
// It outputs exactly headerLines (10) lines, left-aligned.
func renderHeader(width int) string {
	var sb strings.Builder

	logoLines := strings.Split(strings.TrimPrefix(asciiLogo, "\n"), "\n")
	for _, line := range logoLines {
		sb.WriteString(logoStyle.Render(strings.TrimRight(line, " ")))
		sb.WriteByte('\n')
	}

	sb.WriteByte('\n')

	sb.WriteString(titleStyle.Render("Vibe & Kalika Code"))
	sb.WriteByte('\n')

	sb.WriteString(subtitleStyle.Render("Gerenciador de ambientes de desenvolvimento com IA"))
	sb.WriteByte('\n')

	sb.WriteString(borderLineStyle.Render(strings.Repeat("─", width)))
	sb.WriteByte('\n')

	return sb.String()
}

// computeSuggestions returns matching commands or argument completions for the given input.
// When the input starts with "/uninstall " it matches installed package names.
func computeSuggestions(input string, installed []string) []replCommand {
	if input == "" {
		return nil
	}

	// Argument completion for /uninstall <package>.
	if rest, ok := strings.CutPrefix(input, "/uninstall "); ok {
		lower := strings.ToLower(rest)
		matches := make([]replCommand, 0, len(installed))
		for _, pkg := range installed {
			if strings.HasPrefix(strings.ToLower(pkg), lower) {
				matches = append(matches, replCommand{
					cmd:  "/uninstall " + pkg,
					hint: "pacote instalado",
				})
			}
		}
		return matches
	}

	// Command-level prefix completion.
	lower := strings.ToLower(input)
	matches := make([]replCommand, 0, len(replCommands))
	for _, c := range replCommands {
		if strings.HasPrefix(strings.ToLower(c.cmd), lower) && c.cmd != input {
			matches = append(matches, c)
		}
	}
	return matches
}

// replExecCmd implements tea.ExecCommand for interactive cobra subcommands
// that require direct terminal access (e.g. huh forms).
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

// Run builds and executes the cobra subcommand with the terminal handed over.
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
