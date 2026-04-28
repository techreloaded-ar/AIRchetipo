package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
)

// newStatusCmd builds `archetipo status set` -> transition_status.
func newStatusCmd(s streams) *cobra.Command {
	root := &cobra.Command{Use: "status", Short: "Status operations"}
	root.AddCommand(newStatusSetCmd(s))
	return root
}

func newStatusSetCmd(s streams) *cobra.Command {
	var (
		ref string
		to  string
	)
	cmd := &cobra.Command{
		Use:   "set",
		Short: "Transition a story to a new workflow status",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ref == "" || to == "" {
				return errInvalidUsage("missing --ref or --to", "pass --ref US-XXX --to PLANNED")
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.TransitionStatus(ctx, ref, domain.Status(to))
			})
		},
	}
	cmd.Flags().StringVar(&ref, "ref", "", "story reference (US-XXX)")
	cmd.Flags().StringVar(&to, "to", "", "target status (e.g. PLANNED)")
	return cmd
}
