package manifest

type Manifest struct {
	Name        string   `yaml:"name"`
	Version     string   `yaml:"version"`
	Description string   `yaml:"description"`
	Author      string   `yaml:"author,omitempty"`
	License     string   `yaml:"license,omitempty"`
	Targets     []string `yaml:"targets"`
	Agents      []Agent  `yaml:"agents,omitempty"`
	Skills      []Skill  `yaml:"skills,omitempty"`
}

type Agent struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description,omitempty"`
	Source      string                 `yaml:"source"`
	Targets     map[string]AgentTarget `yaml:"targets"`
}

type Skill struct {
	Name        string                 `yaml:"name"`
	Description string                 `yaml:"description,omitempty"`
	Source      string                 `yaml:"source"`
	Targets     map[string]SkillTarget `yaml:"targets"`
}

type AgentTarget struct {
	Scope string `yaml:"scope,omitempty"`
	Mode  string `yaml:"mode,omitempty"`
}

type SkillTarget struct {
	Scope string `yaml:"scope,omitempty"`
}
