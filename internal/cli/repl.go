package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
)

const replHelp = `Available commands:
  /detect        — detect platforms
  /doctor        — run health checks
  /help          — list available commands
  /init          — interactive init
  /install       — install a package
  /uninstall     — uninstall a package
  /update        — update packages
  /update --self — update the CLI itself
  /validate      — validate a package
  /exit, /quit   — exit the REPL`

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

func runREPL(baseFactory func() *cobra.Command, out io.Writer) {
	fmt.Fprintln(out, "Welcome to vkc interactive mode. Type /help for available commands.")

	isTTY := isatty.IsTerminal(os.Stdin.Fd())

	for {
		if isTTY {
			fmt.Fprint(out, "> ")
		}

		line, err := readLine(os.Stdin)
		if err != nil {
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		switch line {
		case "/exit", "/quit":
			return
		case "/help":
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
