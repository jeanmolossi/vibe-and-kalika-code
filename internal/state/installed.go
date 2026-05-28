package state

type Store struct {
	Installations []Installation `yaml:"installations"`
}

type Installation struct {
	Package        string   `yaml:"package"`
	Version        string   `yaml:"version"`
	Source         string   `yaml:"source"`
	InstalledAt    string   `yaml:"installed_at"`
	Platforms      []string `yaml:"platforms"`
	Files          []string `yaml:"files"`
	ManagedMarkers []string `yaml:"managed_markers,omitempty"`
	ReportPath     string   `yaml:"report_path,omitempty"`
	BackupPath     string   `yaml:"backup_path,omitempty"`
}
