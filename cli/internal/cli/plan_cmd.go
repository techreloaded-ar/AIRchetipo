package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/iox"
)

// newPlanCmd builds `archetipo plan save` -> save_plan.
//
// Stdin payload: {"plan_body":"...","tasks":[...]}
func newPlanCmd(s streams) *cobra.Command {
	root := &cobra.Command{Use: "plan", Short: "Plan operations"}
	root.AddCommand(newPlanSaveCmd(s))
	return root
}

func newPlanSaveCmd(s streams) *cobra.Command {
	var ref string
	cmd := &cobra.Command{
		Use:   "save",
		Short: "Persist an implementation plan for a story",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ref == "" {
				return errInvalidUsage("missing --ref", "pass --ref US-XXX")
			}
			var input domain.PlanInput
			if err := iox.ReadJSON(s.in, &input); err != nil {
				iox.WriteError(s.err, err)
				return err
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.SavePlan(ctx, ref, input)
			})
		},
	}
	cmd.Flags().StringVar(&ref, "ref", "", "parent story reference (US-XXX)")
	return cmd
}
