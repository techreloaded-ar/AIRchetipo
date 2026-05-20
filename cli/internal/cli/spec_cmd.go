package cli

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/iox"
)

// newSpecCmd builds `archetipo spec ...` with eight leaves:
//
//	spec add    -> idempotent backlog create-or-append (stdin: {"stories":[...]})
//	spec show   -> read story body + tasks by code
//	spec next   -> auto-pick first eligible story by --status (priority+code)
//	spec list   -> aggregated read: items (optionally filtered) + summary metadata
//	spec plan   -> save plan + transition TODO → PLANNED (stdin: {"plan_body","tasks"})
//	spec start  -> transition PLANNED → IN PROGRESS (idempotent)
//	spec review -> transition IN PROGRESS → REVIEW; --file (optional) is a closing comment
//	spec move   -> reposition a story within the board or across workflow columns
func newSpecCmd(s streams) *cobra.Command {
	root := &cobra.Command{Use: "spec", Short: "Spec (user story) operations"}
	root.AddCommand(
		newSpecAddCmd(s),
		newSpecShowCmd(s),
		newSpecNextCmd(s),
		newSpecListCmd(s),
		newSpecPlanCmd(s),
		newSpecStartCmd(s),
		newSpecReviewCmd(s),
		newSpecMoveCmd(s),
	)
	return root
}

// storiesPayload is the canonical stdin shape for `spec add`.
type storiesPayload struct {
	Schema  string         `json:"schema,omitempty"`
	Kind    string         `json:"kind,omitempty"`
	Stories []domain.Story `json:"stories"`
}

func newSpecAddCmd(s streams) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add stories to the backlog (idempotent: skips duplicate codes)",
		Long: "Reads a YAML or JSON payload from --file and writes the stories to the backlog. " +
			"Creates the backlog when missing, appends otherwise. Stories whose code is " +
			"already present are skipped and reported in data.skipped.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if filePath == "" {
				return errInvalidUsage("missing input file", "pass --file path/to/stories.yaml or --file -")
			}
			payload, err := readStoriesPayload(s.in, filePath)
			if err != nil {
				return err
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return idempotentBacklogWrite(ctx, c, payload.Stories)
			})
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "path to a YAML or JSON payload file, or - for stdin")
	return cmd
}

func readStoriesPayload(stdin io.Reader, path string) (storiesPayload, error) {
	var p storiesPayload
	if err := readStructuredInput(stdin, path, &p); err != nil {
		return storiesPayload{}, err
	}
	if len(p.Stories) == 0 {
		return storiesPayload{}, iox.NewInvalidInput("no stories in input payload", "expected {stories:[...]}", nil)
	}
	return p, nil
}

// idempotentBacklogWrite implements the create-or-append semantics of
// `spec add`: a fresh backlog is initialized, an existing backlog is
// extended skipping codes already present.
func idempotentBacklogWrite(ctx context.Context, c connector.Connector, stories []domain.Story) (domain.WriteResult, error) {
	summary, err := c.ReadExistingBacklog(ctx)
	backlogEmpty := false
	if err != nil {
		if ce, ok := err.(*iox.CodedError); ok && ce.Code == iox.CodePreconditionMissing {
			backlogEmpty = true
		} else {
			return domain.WriteResult{}, err
		}
	} else if len(summary.Codes) == 0 {
		backlogEmpty = true
	}

	if backlogEmpty {
		return c.SaveInitialBacklog(ctx, stories)
	}

	existing := make(map[string]struct{}, len(summary.Codes))
	for _, code := range summary.Codes {
		existing[code] = struct{}{}
	}
	fresh := make([]domain.Story, 0, len(stories))
	skipped := make([]string, 0)
	for _, st := range stories {
		if _, ok := existing[st.Code]; ok {
			skipped = append(skipped, st.Code)
			continue
		}
		fresh = append(fresh, st)
	}
	if len(fresh) == 0 {
		return domain.WriteResult{OK: true, Skipped: skipped}, nil
	}
	res, err := c.AppendStories(ctx, fresh)
	if err != nil {
		return domain.WriteResult{}, err
	}
	if len(skipped) > 0 {
		res.Skipped = skipped
	}
	return res, nil
}

func newSpecShowCmd(s streams) *cobra.Command {
	return &cobra.Command{
		Use:   "show US-XXX",
		Short: "Read a story's body and tasks by code",
		Long:  "Looks up the story by code and returns its body and current task list.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return errInvalidUsage("missing story code", "pass US-XXX as positional argument")
			}
			ref := strings.TrimSpace(args[0])
			if ref == "" {
				return errInvalidUsage("missing story code", "pass US-XXX as positional argument")
			}
			return withConnector(cmd, s, "story", func(ctx context.Context, c connector.Connector) (any, error) {
				st, err := c.ReadStoryDetail(ctx, ref)
				if err != nil {
					return nil, err
				}
				return loadStoryWithTasks(ctx, c, st)
			})
		},
	}
}

func newSpecNextCmd(s streams) *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:   "next",
		Short: "Auto-pick the first eligible story by --status (priority+code)",
		Long: "Selects the first eligible story whose workflow status matches --status, " +
			"ordered by priority and code, and returns its body and tasks.",
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if strings.TrimSpace(status) == "" {
				return errInvalidUsage("missing --status", "pass --status TODO|PLANNED|IN PROGRESS|REVIEW|DONE")
			}
			return withConnector(cmd, s, "story", func(ctx context.Context, c connector.Connector) (any, error) {
				st, err := c.SelectStory(ctx, domain.SelectQuery{EligibleStatuses: []domain.Status{domain.Status(status)}})
				if err != nil {
					return nil, err
				}
				return loadStoryWithTasks(ctx, c, st)
			})
		},
	}
	cmd.Flags().StringVar(&status, "status", "", "workflow status to auto-pick from")
	return cmd
}

// loadStoryWithTasks builds the `story` envelope payload shared by spec show
// and spec next. A story without a plan reports an empty task list rather than
// an error.
func loadStoryWithTasks(ctx context.Context, c connector.Connector, st domain.Story) (map[string]any, error) {
	tasks, err := c.ReadStoryTasks(ctx, st.Code)
	if err != nil {
		if ce, ok := err.(*iox.CodedError); ok && ce.Code == iox.CodePreconditionMissing {
			tasks = []domain.Task{}
		} else {
			return nil, err
		}
	}
	return map[string]any{"story": st, "tasks": tasks}, nil
}

func newSpecListCmd(s streams) *cobra.Command {
	var status string
	cmd := &cobra.Command{
		Use:   "list",
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

func newSpecPlanCmd(s streams) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "plan US-XXX",
		Short: "Save the implementation plan for a story and transition to PLANNED",
		Long: "Reads a YAML or JSON payload from --file. " +
			"Idempotent: re-running on a PLANNED story upserts the plan body without erroring. " +
			"Errors with E_CONFLICT when the story is past PLANNED.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := strings.TrimSpace(args[0])
			if ref == "" {
				return errInvalidUsage("missing story code", "pass US-XXX as positional argument")
			}
			if filePath == "" {
				return errInvalidUsage("missing input file", "pass --file path/to/plan.yaml or --file -")
			}
			var input domain.PlanInput
			if err := readStructuredInput(s.in, filePath, &input); err != nil {
				return err
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				story, err := c.ReadStoryDetail(ctx, ref)
				if err != nil {
					return nil, err
				}
				switch story.Status {
				case domain.StatusTodo:
					res, err := c.SavePlan(ctx, ref, input)
					if err != nil {
						return nil, err
					}
					if _, err := c.TransitionStatus(ctx, ref, domain.StatusPlanned); err != nil {
						return nil, err
					}
					return res, nil
				case domain.StatusPlanned:
					return c.SavePlan(ctx, ref, input)
				default:
					return nil, iox.NewConflict(
						fmt.Sprintf("cannot plan story %s: status is %s, expected TODO or PLANNED", ref, story.Status),
						"inspect the current status with `archetipo spec show "+ref+"`", nil)
				}
			})
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "path to a YAML or JSON payload file, or - for stdin")
	return cmd
}

func newSpecStartCmd(s streams) *cobra.Command {
	return &cobra.Command{
		Use:   "start US-XXX",
		Short: "Transition a planned story to IN PROGRESS (idempotent)",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := strings.TrimSpace(args[0])
			if ref == "" {
				return errInvalidUsage("missing story code", "pass US-XXX as positional argument")
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				return transitionWithValidation(ctx, c, ref, "start", domain.StatusPlanned, domain.StatusInProgress)
			})
		},
	}
}

func newSpecReviewCmd(s streams) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "review US-XXX",
		Short: "Transition a story to REVIEW; --file (or stdin) is posted as a closing comment",
		Long: "Transitions the story from IN PROGRESS to REVIEW and, when a non-empty body is provided " +
			"via --file or stdin, posts it as a closing comment on the parent issue. Connectors " +
			"without comment support silently ignore the body.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := strings.TrimSpace(args[0])
			if ref == "" {
				return errInvalidUsage("missing story code", "pass US-XXX as positional argument")
			}
			comment, err := readRawInput(s.in, filePath)
			if err != nil {
				return err
			}
			return withConnector(cmd, s, "write_result", func(ctx context.Context, c connector.Connector) (any, error) {
				res, err := transitionWithValidation(ctx, c, ref, "review", domain.StatusInProgress, domain.StatusReview)
				if err != nil {
					return nil, err
				}
				if len(strings.TrimSpace(string(comment))) > 0 {
					if _, err := c.PostComment(ctx, ref, string(comment)); err != nil {
						return nil, err
					}
				}
				return res, nil
			})
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "path to the closing comment, or - for stdin (default: stdin)")
	return cmd
}

// validMoveTargets lists the board columns accepted by `spec move --to`.
// The list mirrors the mapping in the connector implementations.
var validMoveTargets = []string{"todo", "planned", "in_progress", "review", "done"}

func newSpecMoveCmd(s streams) *cobra.Command {
	var before string
	var after string
	var target string
	cmd := &cobra.Command{
		Use:   "move US-XXX",
		Short: "Move a story within the board or across workflow columns",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if target == "" {
				return errInvalidUsage("missing target column", "pass --to "+strings.Join(validMoveTargets, "|"))
			}
			if !isValidMoveTarget(target) {
				return errInvalidUsage(
					fmt.Sprintf("invalid --to value %q", target),
					"valid columns: "+strings.Join(validMoveTargets, "|"),
				)
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
	cmd.Flags().StringVar(&target, "to", "", "target board column: "+strings.Join(validMoveTargets, "|"))
	cmd.Flags().StringVar(&before, "before", "", "insert before the given story code in the target column")
	cmd.Flags().StringVar(&after, "after", "", "insert after the given story code in the target column")
	return cmd
}

func isValidMoveTarget(t string) bool {
	for _, v := range validMoveTargets {
		if v == t {
			return true
		}
	}
	return false
}

// transitionWithValidation enforces the idempotent + validated transition rules
// shared by `spec start` and `spec review`. Calling the verb when the story
// is already at the target state is a no-op success; calling it from any
// status other than the expected source returns E_CONFLICT.
func transitionWithValidation(ctx context.Context, c connector.Connector, ref, verb string, source, target domain.Status) (domain.WriteResult, error) {
	story, err := c.ReadStoryDetail(ctx, ref)
	if err != nil {
		return domain.WriteResult{}, err
	}
	if story.Status == target {
		return domain.WriteResult{OK: true, Refs: []domain.Ref{{Code: story.Code}}}, nil
	}
	if story.Status != source {
		return domain.WriteResult{}, iox.NewConflict(
			fmt.Sprintf("cannot %s story %s: status is %s, expected %s", verb, ref, story.Status, source),
			fmt.Sprintf("transition the story to %s before running `archetipo spec %s`", source, verb),
			nil)
	}
	return c.TransitionStatus(ctx, ref, target)
}

func readStructuredInput(stdin io.Reader, path string, v any) error {
	raw, err := readRawInput(stdin, path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(raw, v); err != nil {
		return iox.NewInvalidInput("invalid structured input", "expected YAML or JSON payload", err)
	}
	return nil
}

// readRawInput reads raw bytes from a file path or stdin. When path is empty
// or "-", input is taken from stdin. Returns iox-typed errors on read failure.
func readRawInput(stdin io.Reader, path string) ([]byte, error) {
	if path == "" || path == "-" {
		raw, err := io.ReadAll(stdin)
		if err != nil {
			return nil, iox.NewInvalidInput("reading stdin", "", err)
		}
		return raw, nil
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, iox.NewInvalidInput(fmt.Sprintf("reading input file %s", path), "", err)
	}
	return raw, nil
}
