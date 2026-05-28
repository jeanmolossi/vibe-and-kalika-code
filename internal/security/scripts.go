package security

import (
	"io/fs"
	"path/filepath"
	"strings"
)

var scriptExts = map[string]struct{}{
	".sh":  {},
	".ps1": {},
	".bat": {},
	".cmd": {},
}

func IsScriptPath(path string) bool {
	_, ok := scriptExts[strings.ToLower(filepath.Ext(path))]
	return ok
}

func FindScripts(root string) ([]string, error) {
	var scripts []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if IsScriptPath(path) {
			scripts = append(scripts, path)
		}
		return nil
	})
	return scripts, err
}
