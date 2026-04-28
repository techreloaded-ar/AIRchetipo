package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
)

// newTasksCmd builds `archetipo tasks read` -> read_story_tasks.
func newTasksCmd(s streams) *cobra.Command {
	root := &cobra.Command{Use: "tasks", Short: "Task list operations"}
	root.AddCommand(newTasksReadCmd(s))
	return root
}

func newTasksReadCmd(s streams) *cobra.Command {
	var ref string
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read the implementation task list for a story",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ref == "" {
				return errInvalidUsage("missing --ref", "pass --ref US-XXX")
			}
			return withConnector(cmd, s, "tasks", func(ctx context.Context, c connector.Connector) (any, error) {
				items, err := c.ReadStoryTasks(ctx, ref)
				if err != nil {
					return nil, err
				}
				return map[string]any{"items": items}, nil
			})
		},
	}
	cmd.Flags().StringVar(&ref, "ref", "", "parent story reference (US-XXX)")
	return cmd
}
