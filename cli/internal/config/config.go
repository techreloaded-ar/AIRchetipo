// Package config loads and validates .archetipo/config.yaml.
//
// The config file lives in the *target project* (the project where the user
// runs the CLI), not in the CLI repo. It selects which connector implements
// the contract, where artifacts live, and how workflow statuses are labelled.
//
// Defaults: when config.yaml does not exist, the file connector is selected
// with the canonical paths and statuses documented in contracts.md.
package config

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
)

// Path of the config file relative to the project root.
const RelativePath = ".archetipo/config.yaml"

// Connector identifiers recognized by the registry.
const (
	ConnectorFile   = "file"
	ConnectorGitHub = "github"
)

// Config is the parsed shape of .archetipo/config.yaml.
type Config struct {
	Connector string             `yaml:"connector" json:"connector"`
	Paths     domain.ConfigPaths `yaml:"paths" json:"paths"`
	Workflow  domain.WorkflowConfig `yaml:"workflow" json:"workflow"`
	GitHub    GitHubConfig       `yaml:"github" json:"github,omitempty"`
	// ProjectRoot is the absolute path of the directory that contains
	// .archetipo/. Set by Load; not present in the YAML file.
	ProjectRoot string `yaml:"-" json:"project_root"`
}

// GitHubConfig holds connector-specific overrides. Owner and project number
// are auto-detected from `gh` when empty.
type GitHubConfig struct {
	Owner         string `yaml:"owner,omitempty" json:"owner,omitempty"`
	ProjectNumber int    `yaml:"project_number,omitempty" json:"project_number,omitempty"`
}

// Default returns the canonical default config (file connector, English status
// labels). Used when the project has no config.yaml.
func Default() Config {
	return Config{
		Connector: ConnectorFile,
		Paths: domain.ConfigPaths{
			PRD:         "docs/PRD.md",
			Backlog:     "docs/BACKLOG.md",
			Planning:    "docs/planning/",
			Mockups:     "docs/mockups/",
			TestResults: "docs/test-results/",
		},
		Workflow: domain.WorkflowConfig{
			Statuses: domain.StatusLabels{
				Todo:       string(domain.StatusTodo),
				Planned:    string(domain.StatusPlanned),
				InProgress: string(domain.StatusInProgress),
				Review:     string(domain.StatusReview),
				Done:       string(domain.StatusDone),
			},
		},
	}
}

// Load locates `.archetipo/config.yaml` starting from startDir, walking up
// the directory tree until found or the filesystem root is reached. When
// not found, the default config rooted at startDir is returned.
func Load(startDir string) (Config, error) {
	root, cfgPath, err := find(startDir)
	if err != nil {
		return Config{}, err
	}
	if cfgPath == "" {
		// No config: use default rooted at startDir.
		c := Default()
		abs, _ := filepath.Abs(startDir)
		c.ProjectRoot = abs
		return c, nil
	}
	raw, err := os.ReadFile(cfgPath)
	if err != nil {
		return Config{}, fmt.Errorf("reading %s: %w", cfgPath, err)
	}
	c := Default()
	if err := yaml.Unmarshal(raw, &c); err != nil {
		return Config{}, fmt.Errorf("parsing %s: %w", cfgPath, err)
	}
	c.applyDefaults()
	if err := c.validate(); err != nil {
		return Config{}, err
	}
	c.ProjectRoot = root
	return c, nil
}

// applyDefaults fills empty fields with canonical defaults. Lets the user
// omit unchanged keys from config.yaml.
func (c *Config) applyDefaults() {
	d := Default()
	if c.Connector == "" {
		c.Connector = d.Connector
	}
	if c.Paths.PRD == "" {
		c.Paths.PRD = d.Paths.PRD
	}
	if c.Paths.Backlog == "" {
		c.Paths.Backlog = d.Paths.Backlog
	}
	if c.Paths.Planning == "" {
		c.Paths.Planning = d.Paths.Planning
	}
	if c.Paths.Mockups == "" {
		c.Paths.Mockups = d.Paths.Mockups
	}
	if c.Paths.TestResults == "" {
		c.Paths.TestResults = d.Paths.TestResults
	}
	if c.Workflow.Statuses.Todo == "" {
		c.Workflow.Statuses.Todo = d.Workflow.Statuses.Todo
	}
	if c.Workflow.Statuses.Planned == "" {
		c.Workflow.Statuses.Planned = d.Workflow.Statuses.Planned
	}
	if c.Workflow.Statuses.InProgress == "" {
		c.Workflow.Statuses.InProgress = d.Workflow.Statuses.InProgress
	}
	if c.Workflow.Statuses.Review == "" {
		c.Workflow.Statuses.Review = d.Workflow.Statuses.Review
	}
	if c.Workflow.Statuses.Done == "" {
		c.Workflow.Statuses.Done = d.Workflow.Statuses.Done
	}
}

func (c *Config) validate() error {
	switch c.Connector {
	case ConnectorFile, ConnectorGitHub:
	default:
		return fmt.Errorf("unknown connector %q (allowed: file, github)", c.Connector)
	}
	return nil
}

// AbsPath joins p against the project root if p is relative.
func (c Config) AbsPath(p string) string {
	if filepath.IsAbs(p) {
		return p
	}
	return filepath.Join(c.ProjectRoot, p)
}

// find walks up from start looking for .archetipo/config.yaml. Returns the
// project root (the directory that contains .archetipo/) and the absolute
// path of the config file. If neither is found, returns ("", "", nil).
func find(start string) (root, cfg string, err error) {
	abs, err := filepath.Abs(start)
	if err != nil {
		return "", "", err
	}
	dir := abs
	for {
		candidate := filepath.Join(dir, RelativePath)
		info, statErr := os.Stat(candidate)
		if statErr == nil && !info.IsDir() {
			return dir, candidate, nil
		}
		if statErr != nil && !errors.Is(statErr, fs.ErrNotExist) {
			return "", "", statErr
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", "", nil
		}
		dir = parent
	}
}
