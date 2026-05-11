package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/iox"
)

// newBacklogCmd builds `archetipo backlog ...` with one leaf:
//
//	backlog show -> aggregated read: items (optionally filtered by status) +
//	                idempotency summary (codes, last_code, epics, titles).
//
// Story-level writes live under `archetipo story add`.
func newBacklogCmd(s streams) *cobra.Command {
	root := &cobra.Command{
		Use:   "backlog",
		Short: "Backlog read operations",
	}
	root.AddCommand(
		newBacklogShowCmd(s),
		newBacklogReorderCmd(s),
	)
	return root
}

func newBacklogShowCmd(s streams) *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:   "show",
		Short: "List backlog stories (optionally filtered by status) with summary metadata",
		Long: "Returns {items, summary} in a single envelope. items is filtered by --status when provided; " +
			"summary is always the full backlog metadata (codes, last_code, epics, titles).",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return withConnector(cmd, s, "backlog", func(ctx context.Context, c connector.Connector) (any, error) {
				items, err := c.FetchBacklogItems(ctx, domain.Status(status))
				if err != nil {
					return nil, err
				}
				summary, err := c.ReadExistingBacklog(ctx)
				if err != nil {
					if ce, ok := err.(*iox.CodedError); ok && ce.Code == iox.CodePreconditionMissing {
						summary = domain.BacklogSummary{}
					} else {
						return nil, err
					}
				}
				return map[string]any{
					"items":   items,
					"summary": summary,
				}, nil
			})
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "filter items by workflow status (e.g. TODO)")
	return cmd
}

func newBacklogReorderCmd(s streams) *cobra.Command {
	var before string
	var after string
	cmd := &cobra.Command{
		Use:   "reorder US-XXX",
		Short: "Move a story within the backlog order",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if before != "" && after != "" {
				return errInvalidUsage("before and after are mutually exclusive", "pass only one anchor")
			}
			ref := args[0]
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.ReorderBacklog(ctx, ref, domain.ReorderAnchor{Before: before, After: after})
			})
		},
	}
	cmd.Flags().StringVar(&before, "before", "", "insert before the given story code")
	cmd.Flags().StringVar(&after, "after", "", "insert after the given story code")
	return cmd
}
