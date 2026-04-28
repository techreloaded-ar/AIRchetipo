package cli

import (
	"context"
	"fmt"
	"io"

	"github.com/spf13/cobra"

	// Concrete connectors register themselves via init().
	_ "github.com/techreloaded-ar/ARchetipo/cli/internal/connector/builtin"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/version"
)

// Execute runs the archetipo CLI with the given args and returns the process
// exit code. Stdin/stdout/stderr are taken as parameters so tests can drive the
// CLI without touching the real OS streams.
//
// Exit codes follow contracts.md:
//
//	0  ok
//	1  generic error
//	2  input/validation error
//	3  connector error (auth, network, gh)
//	4  precondition missing (e.g. backlog absent)
func Execute(args []string, stdin io.Reader, stdout, stderr io.Writer) int {
	root := newRootCmd(stdin, stdout, stderr)
	root.SetArgs(args)
	root.SetIn(stdin)
	root.SetOut(stdout)
	root.SetErr(stderr)
	if err := root.Execute(); err != nil {
		// cobra has already written the error to stderr via the default
		// behavior; we just translate to an exit code.
		return exitCodeFor(err)
	}
	return 0
}

func newRootCmd(stdin io.Reader, stdout, stderr io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:           "archetipo",
		Short:         "ARchetipo connector CLI",
		Long:          "Deterministic CLI implementing the ARchetipo connector contracts (file and github).",
		SilenceUsage:  true,
		SilenceErrors: false,
		Version:       version.Version,
	}
	cmd.SetVersionTemplate(fmt.Sprintf("archetipo %s\n", version.Version))
	// Cancel sub-command contexts when cobra is told to abort.
	cmd.SetContext(context.Background())

	s := streams{in: stdin, out: stdout, err: stderr}
	cmd.AddCommand(
		newInitCmd(s),
		newBacklogCmd(s),
		newStoryCmd(s),
		newTasksCmd(s),
		newPRDCmd(s),
		newPlanCmd(s),
		newStatusCmd(s),
		newTaskCmd(s),
		newCommentCmd(s),
	)
	return cmd
}

// exitCodeFor maps a returned error to the documented exit code. Specific
// error types defined in internal/iox carry their own code; everything else is
// a generic error (1).
func exitCodeFor(err error) int {
	if err == nil {
		return 0
	}
	if coded, ok := err.(interface{ ExitCode() int }); ok {
		return coded.ExitCode()
	}
	return 1
}
