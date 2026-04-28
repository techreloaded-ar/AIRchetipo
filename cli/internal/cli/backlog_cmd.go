package cli

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/iox"
)

// newBacklogCmd builds `archetipo backlog ...` and its 4 leaves:
//
//	backlog list      -> fetch_backlog_items
//	backlog existing  -> read_existing_backlog
//	backlog save      -> save_initial_backlog (stdin JSON: {"stories":[...]})
//	backlog append    -> append_stories       (stdin JSON: {"stories":[...]})
func newBacklogCmd(s streams) *cobra.Command {
	root := &cobra.Command{
		Use:   "backlog",
		Short: "Backlog operations",
	}
	root.AddCommand(newBacklogListCmd(s), newBacklogExistingCmd(s), newBacklogSaveCmd(s), newBacklogAppendCmd(s))
	return root
}

func newBacklogListCmd(s streams) *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List backlog stories, optionally filtered by status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return withConnector(cmd, s, "stories", func(ctx context.Context, c connector.Connector) (any, error) {
				items, err := c.FetchBacklogItems(ctx, domain.Status(status))
				if err != nil {
					return nil, err
				}
				return map[string]any{"items": items}, nil
			})
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "filter by workflow status (e.g. TODO)")
	return cmd
}

func newBacklogExistingCmd(s streams) *cobra.Command {
	return &cobra.Command{
		Use:   "existing",
		Short: "Idempotency metadata about the existing backlog",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return withConnector(cmd, s, "backlog_summary", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.ReadExistingBacklog(ctx)
			})
		},
	}
}

// storiesPayload is the canonical stdin shape for save/append.
type storiesPayload struct {
	Schema  string          `json:"schema,omitempty"`
	Kind    string          `json:"kind,omitempty"`
	Stories []domain.Story `json:"stories"`
}

func newBacklogSaveCmd(s streams) *cobra.Command {
	return &cobra.Command{
		Use:   "save",
		Short: "Create the initial backlog from stdin JSON",
		Long:  "Reads {\"stories\":[...]} from stdin and creates the backlog. Fails if a non-empty backlog already exists.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			payload, err := readStories(s.in)
			if err != nil {
				iox.WriteError(s.err, err)
				return err
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.SaveInitialBacklog(ctx, payload.Stories)
			})
		},
	}
}

func newBacklogAppendCmd(s streams) *cobra.Command {
	return &cobra.Command{
		Use:   "append",
		Short: "Append stories from stdin JSON to an existing backlog",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			payload, err := readStories(s.in)
			if err != nil {
				iox.WriteError(s.err, err)
				return err
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.AppendStories(ctx, payload.Stories)
			})
		},
	}
}

func readStories(r io.Reader) (storiesPayload, error) {
	var p storiesPayload
	if err := iox.ReadJSON(r, &p); err != nil {
		return storiesPayload{}, err
	}
	if len(p.Stories) == 0 {
		return storiesPayload{}, iox.NewInvalidInput("no stories on stdin", "expected {\"stories\":[...]}", nil)
	}
	return p, nil
}
