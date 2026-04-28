package cli

import (
	"context"
	"io"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/iox"
)

// newCommentCmd builds `archetipo comment post` -> post_comment.
//
// Stdin is the raw markdown body. The file connector is a no-op.
func newCommentCmd(s streams) *cobra.Command {
	root := &cobra.Command{Use: "comment", Short: "Comment operations"}
	root.AddCommand(newCommentPostCmd(s))
	return root
}

func newCommentPostCmd(s streams) *cobra.Command {
	var ref string
	cmd := &cobra.Command{
		Use:   "post",
		Short: "Post a comment on a story (no-op for the file connector)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if ref == "" {
				return errInvalidUsage("missing --ref", "pass --ref US-XXX")
			}
			body, err := io.ReadAll(s.in)
			if err != nil {
				return iox.NewInvalidInput("reading stdin", "", err)
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return c.PostComment(ctx, ref, string(body))
			})
		},
	}
	cmd.Flags().StringVar(&ref, "ref", "", "story reference (US-XXX)")
	return cmd
}
