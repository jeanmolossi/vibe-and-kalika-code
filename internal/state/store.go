package state

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

func Path(projectRoot string) string {
	return filepath.Join(projectRoot, ".ai-setup", "installed.yaml")
}

func Read(projectRoot string) (*Store, error) {
	path := Path(projectRoot)
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

func Write(projectRoot string, store *Store) error {
	path := Path(projectRoot)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := yaml.Marshal(store)
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}
