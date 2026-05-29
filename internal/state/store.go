package state

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// StateDir returns the global .ai-setup directory, always under the user's
// home directory. Set VKC_STATE_DIR to override (useful in tests).
func StateDir() (string, error) {
	if override := os.Getenv("VKC_STATE_DIR"); override != "" {
		return override, nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home directory: %w", err)
	}
	return filepath.Join(home, ".ai-setup"), nil
}

// Path returns the path to the global installed.yaml state file.
func Path() (string, error) {
	dir, err := StateDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "installed.yaml"), nil
}

// Read loads the global installation state. Returns an empty Store if the
// state file does not yet exist.
func Read() (*Store, error) {
	path, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Store{}, nil
	}
	if err != nil {
		return nil, err
	}
	var store Store
	if err := yaml.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("parse installed state: %w", err)
	}
	return &store, nil
}

// Write persists the installation state to the global state file.
func Write(store *Store) error {
	path, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(store)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
