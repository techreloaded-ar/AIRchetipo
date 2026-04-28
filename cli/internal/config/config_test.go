package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultWhenConfigMissing(t *testing.T) {
	dir := t.TempDir()
	c, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if c.Connector != ConnectorFile {
		t.Errorf("expected default connector %q, got %q", ConnectorFile, c.Connector)
	}
	if c.Paths.Backlog != "docs/BACKLOG.md" {
		t.Errorf("default backlog path: %q", c.Paths.Backlog)
	}
}

func TestLoadFromConfigFile(t *testing.T) {
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, ".archetipo"), 0o755))
	must(t, os.WriteFile(filepath.Join(root, RelativePath), []byte(`connector: github
paths:
  backlog: my/BL.md
workflow:
  statuses:
    todo: A_FARE
    planned: PIANIFICATO
    in_progress: IN CORSO
    review: REVISIONE
    done: FATTO
`), 0o644))

	c, err := Load(root)
	if err != nil {
		t.Fatal(err)
	}
	if c.Connector != ConnectorGitHub {
		t.Errorf("connector: %q", c.Connector)
	}
	if c.Paths.Backlog != "my/BL.md" {
		t.Errorf("backlog: %q", c.Paths.Backlog)
	}
	// Defaults preserved for unspecified path keys.
	if c.Paths.PRD != "docs/PRD.md" {
		t.Errorf("PRD default lost: %q", c.Paths.PRD)
	}
	if c.Workflow.Statuses.Todo != "A_FARE" {
		t.Errorf("status override lost: %q", c.Workflow.Statuses.Todo)
	}
	if c.ProjectRoot != root {
		t.Errorf("project root: %q want %q", c.ProjectRoot, root)
	}
}

func TestLoadFromSubdirectoryWalksUp(t *testing.T) {
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, ".archetipo"), 0o755))
	must(t, os.WriteFile(filepath.Join(root, RelativePath), []byte(`connector: file
`), 0o644))
	sub := filepath.Join(root, "src", "deep")
	must(t, os.MkdirAll(sub, 0o755))

	c, err := Load(sub)
	if err != nil {
		t.Fatal(err)
	}
	if c.ProjectRoot != root {
		t.Errorf("project root: %q want %q", c.ProjectRoot, root)
	}
}

func TestUnknownConnectorRejected(t *testing.T) {
	root := t.TempDir()
	must(t, os.MkdirAll(filepath.Join(root, ".archetipo"), 0o755))
	must(t, os.WriteFile(filepath.Join(root, RelativePath), []byte(`connector: gitlab
`), 0o644))
	_, err := Load(root)
	if err == nil {
		t.Fatal("expected error for unknown connector")
	}
}

func must(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatal(err)
	}
}
