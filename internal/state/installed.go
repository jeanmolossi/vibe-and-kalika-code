package state

type Store struct {
	Installations []Installation `yaml:"installations"`
}

// AgentBlock records a managed block that was merged into an AGENTS.md file.
type AgentBlock struct {
	Path      string `yaml:"path"`
	AgentName string `yaml:"agent_name"`
}

type Installation struct {
	Package        string       `yaml:"package"`
	Version        string       `yaml:"version"`
	Source         string       `yaml:"source"`
	InstalledAt    string       `yaml:"installed_at"`
	Platforms      []string     `yaml:"platforms"`
	Files          []string     `yaml:"files"`
	CreatedFiles   []string     `yaml:"created_files,omitempty"`
	ManagedMarkers []string     `yaml:"managed_markers,omitempty"`
	AgentBlocks    []AgentBlock `yaml:"agent_blocks,omitempty"`
	ReportPath     string       `yaml:"report_path,omitempty"`
	BackupPath     string       `yaml:"backup_path,omitempty"`
}
