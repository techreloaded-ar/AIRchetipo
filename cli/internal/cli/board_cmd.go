package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
)

func newBoardCmd(s streams) *cobra.Command {
	root := &cobra.Command{
		Use:   "board",
		Short: "Board operations",
	}
	root.AddCommand(newBoardMoveCmd(s))
	return root
}

func newBoardMoveCmd(s streams) *cobra.Command {
	var before string
	var after string
	var target string
	cmd := &cobra.Command{
		Use:   "move US-XXX",
		Short: "Move a story within the board or across workflow columns",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if target == "" {
				return errInvalidUsage("missing target column", "pass --to todo|planned|in_progress|review|done")
			}
			if before != "" && after != "" {
				return errInvalidUsage("before and after are mutually exclusive", "pass only one anchor")
			}
			ref := args[0]
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.MoveBoardCard(ctx, ref, target, domain.ReorderAnchor{Before: before, After: after})
			})
		},
	}
	cmd.Flags().StringVar(&target, "to", "", "target board column: todo|planned|in_progress|review|done")
	cmd.Flags().StringVar(&before, "before", "", "insert before the given story code in the target column")
	cmd.Flags().StringVar(&after, "after", "", "insert after the given story code in the target column")
	return cmd
}
