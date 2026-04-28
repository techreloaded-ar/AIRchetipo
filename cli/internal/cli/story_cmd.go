package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
)

// newStoryCmd builds `archetipo story ...` with two leaves:
//
//	story select   -> select_story
//	story read     -> read_story_detail
func newStoryCmd(s streams) *cobra.Command {
	root := &cobra.Command{Use: "story", Short: "Story operations"}
	root.AddCommand(newStorySelectCmd(s), newStoryReadCmd(s))
	return root
}

func newStorySelectCmd(s streams) *cobra.Command {
	var (
		storyCode string
		auto      bool
		eligible  string
	)
	cmd := &cobra.Command{
		Use:   "select",
		Short: "Pick a story by code, or auto-select by priority",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return withConnector(cmd, s, "story", func(ctx context.Context, c connector.Connector) (any, error) {
				q := domain.SelectQuery{}
				switch {
				case storyCode != "" && auto:
					return nil, errInvalidUsage("--story and --auto are mutually exclusive", "")
				case storyCode != "":
					q.StoryCode = storyCode
				default:
					q.EligibleStatuses = parseStatuses(eligible)
				}
				return c.SelectStory(ctx, q)
			})
		},
	}
	cmd.Flags().StringVar(&storyCode, "story", "", "select a specific story (e.g. US-005)")
	cmd.Flags().BoolVar(&auto, "auto", false, "auto-select highest-priority story among --eligible")
	cmd.Flags().StringVar(&eligible, "eligible", "TODO", "comma-separated eligible statuses for auto-select")
	return cmd
}

func newStoryReadCmd(s streams) *cobra.Command {
	var ref string
	cmd := &cobra.Command{
		Use:   "read",
		Short: "Read full body/content of a story",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ref == "" {
				return errInvalidUsage("missing --ref", "pass --ref US-XXX")
			}
			return withConnector(cmd, s, "story", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.ReadStoryDetail(ctx, ref)
			})
		},
	}
	cmd.Flags().StringVar(&ref, "ref", "", "story reference (US-XXX or connector ref)")
	return cmd
}

func parseStatuses(csv string) []domain.Status {
	parts := strings.Split(csv, ",")
	out := make([]domain.Status, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			out = append(out, domain.Status(p))
		}
	}
	return out
}
