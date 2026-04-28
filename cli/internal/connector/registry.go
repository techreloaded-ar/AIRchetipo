package connector

import (
	"fmt"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/config"
)

// Builder constructs a Connector from a parsed config. Implementations
// register themselves via Register at init time so the cli package does not
// need to import every connector concrete type.
type Builder func(cfg config.Config) (Connector, error)

var builders = map[string]Builder{}

// Register associates a connector name with its constructor. Concrete
// connector packages call this in their init().
func Register(name string, b Builder) {
	if _, dup := builders[name]; dup {
		panic("connector already registered: " + name)
	}
	builders[name] = b
}

// New builds the connector selected by cfg.Connector.
func New(cfg config.Config) (Connector, error) {
	b, ok := builders[cfg.Connector]
	if !ok {
		return nil, fmt.Errorf("unknown connector %q (registered: %v)", cfg.Connector, registeredNames())
	}
	return b(cfg)
}

func registeredNames() []string {
	out := make([]string, 0, len(builders))
	for k := range builders {
		out = append(out, k)
	}
	return out
}
