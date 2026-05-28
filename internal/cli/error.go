package cli

import "fmt"

type cliError struct {
	code int
	err  error
}

func (e cliError) Error() string { return e.err.Error() }
func (e cliError) Unwrap() error { return e.err }
func (e cliError) ExitCode() int { return e.code }

func exitError(code int, err error) error { return cliError{code: code, err: fmt.Errorf("%w", err)} }

// ExitCodeFrom extracts the exit code from a cliError, returning 1 for generic errors.
func ExitCodeFrom(err error) int {
	type coded interface{ ExitCode() int }
	if asErr, ok := err.(coded); ok {
		return asErr.ExitCode()
	}
	return 1
}
