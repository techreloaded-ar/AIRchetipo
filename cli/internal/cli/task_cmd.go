package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
)

// newTaskCmd builds `archetipo task complete` -> complete_task.
func newTaskCmd(s streams) *cobra.Command {
	root := &cobra.Command{Use: "task", Short: "Task operations"}
	root.AddCommand(newTaskCompleteCmd(s))
	return root
}

func newTaskCompleteCmd(s streams) *cobra.Command {
	var (
		parent string
		ref    string
	)
	cmd := &cobra.Command{
		Use:   "complete",
		Short: "Mark a task as completed",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if parent == "" || ref == "" {
				return errInvalidUsage("missing --parent or --ref", "pass --parent US-XXX --ref TASK-NN")
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.CompleteTask(ctx, parent, ref)
			})
		},
	}
	cmd.Flags().StringVar(&parent, "parent", "", "parent story reference (US-XXX)")
	cmd.Flags().StringVar(&ref, "ref", "", "task reference (TASK-NN)")
	return cmd
}
