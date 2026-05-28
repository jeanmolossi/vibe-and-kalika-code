package app

import (
	"fmt"
	"strings"
)

func errWithIssues(err error, issues []string) error {
	return fmt.Errorf("%w: %s", err, strings.Join(issues, "; "))
}
