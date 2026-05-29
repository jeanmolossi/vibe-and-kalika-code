package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

// REPL command constants used across repl.go and repl_tui.go.
const (
	cmdExit = "/exit"
	cmdQuit = "/quit"
	cmdHelp = "/help"
)

const replHelp = `Comandos disponíveis:
  /detect        — Detecta plataformas disponíveis no projeto
  /doctor        — Executa verificações de saúde do ambiente
  /help          — Lista os comandos disponíveis
  /init          — Inicializa um novo pacote interativamente
  /install       — Instala um pacote a partir de um diretório local ou URL git
  /uninstall     — Remove um pacote instalado anteriormente
  /update        — Atualiza pacotes instalados
  /update --self — Atualiza o próprio CLI
  /validate      — Valida um pacote
  /exit, /quit   — Sai do REPL`

// readLine reads exactly one newline-terminated line from r, one byte at a time.
// This avoids buffering stdin ahead of subcommands that also need to read from it.
func readLine(r io.Reader) (string, error) {
	var buf []byte
	b := make([]byte, 1)
	for {
		n, err := r.Read(b)
		if n > 0 {
			if b[0] == '\n' {
				return strings.TrimRight(string(buf), "\r"), nil
			}
			buf = append(buf, b[0])
		}
		if err != nil {
			if err == io.EOF && len(buf) > 0 {
				return string(buf), nil
			}
			return "", err
		}
	}
}

// runREPL starts the interactive REPL. When stdin is a TTY it launches the rich
// Bubble Tea UI; otherwise it falls back to a plain line-based loop.
func runREPL(baseFactory func() *cobra.Command, out io.Writer) {
	if !isatty.IsTerminal(os.Stdin.Fd()) {
		runReplBasic(baseFactory, out)
		return
	}

	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}

	m := newReplModel(baseFactory, cwd)
	p := tea.NewProgram(m)
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
	}
}

// runReplBasic is the plain line-based fallback used when stdin is not a TTY.
func runReplBasic(baseFactory func() *cobra.Command, out io.Writer) {
	fmt.Fprintln(out, "Welcome to vkc interactive mode. Type /help for available commands.")

	for {
		fmt.Fprint(out, "> ")

		line, err := readLine(os.Stdin)
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch line {
		case cmdExit, cmdQuit:
			return
		case cmdHelp:
			fmt.Fprintln(out, replHelp)
			continue
		}

		if !strings.HasPrefix(line, "/") {
			fmt.Fprintln(out, "hint: commands start with /. Type /help for available commands.")
			continue
		}

		args := strings.Fields(strings.TrimPrefix(line, "/"))
		root := baseFactory()
		root.SetArgs(args)
		root.SetOut(out)
		root.SetErr(os.Stderr)
		root.SilenceErrors = true
		root.SilenceUsage = true

		if err := root.Execute(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
		}
	}
}
