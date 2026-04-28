package filefs_test

import (
	"path/filepath"
	"testing"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/config"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector/conformance"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector/filefs"
)

func TestFilefsConformance(t *testing.T) {
	conformance.Run(t, func(t *testing.T) connector.Connector {
		t.Helper()
		dir := t.TempDir()
		cfg := config.Default()
		cfg.ProjectRoot = dir
		// Use absolute paths so writes land inside the temp dir.
		cfg.Paths.Backlog = filepath.Join(dir, "BACKLOG.md")
		cfg.Paths.Planning = filepath.Join(dir, "planning")
		cfg.Paths.PRD = filepath.Join(dir, "PRD.md")
		return filefs.New(cfg)
	})
}
