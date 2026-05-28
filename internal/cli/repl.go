package cli

import (
	"bufio"
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

func runREPL(baseFactory func() *cobra.Command, out io.Writer) {
	fmt.Fprintln(out, "Welcome to vkc interactive mode. Type /help for available commands.")

	isTTY := isatty.IsTerminal(os.Stdin.Fd())
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if isTTY {
			fmt.Fprint(out, "> ")
		}

		if !scanner.Scan() {
			break
		}

		line := strings.TrimSpace(scanner.Text())
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
