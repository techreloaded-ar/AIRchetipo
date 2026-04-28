package inmemory_test

import (
	"testing"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/config"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector/conformance"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector/inmemory"
)

func TestInmemoryConformance(t *testing.T) {
	conformance.Run(t, func(t *testing.T) connector.Connector {
		t.Helper()
		return inmemory.New(config.Default())
	})
}
