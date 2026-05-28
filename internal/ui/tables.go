package ui

import "strings"

func Bullets(items []string) string { return "- " + strings.Join(items, "\n- ") }
