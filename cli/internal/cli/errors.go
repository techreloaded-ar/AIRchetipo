package cli

import "github.com/techreloaded-ar/ARchetipo/cli/internal/iox"

// errInvalidUsage is a thin shortcut around iox.NewInvalidInput used by
// sub-commands when args/flags do not satisfy the documented contract. The
// error envelope is *not* written here: cobra returns it from RunE and
// cli.Execute lets withConnector decide whether to render it. To stay
// consistent we render it eagerly so the consumer always sees JSON on stderr.
func errInvalidUsage(message, hint string) error {
	return iox.NewInvalidInput(message, hint, nil)
}
