package cli

import (
	"fmt"
	"io"
	"regexp"
	"sort"
	"strings"

	"github.com/spf13/cobra"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/iox"
)

var (
	specCodeRE = regexp.MustCompile(`^US-\d{3,}$`)
	epicCodeRE = regexp.MustCompile(`^EP-\d{3,}$`)
	taskIDRE   = regexp.MustCompile(`^TASK-\d{2,}$`)
)

func newValidateCmd(s streams) *cobra.Command {
	root := &cobra.Command{Use: "validate", Short: "Validate generated ARchetipo artifacts without persisting them"}
	root.AddCommand(newValidateSpecsCmd(s), newValidatePlanCmd(s))
	return root
}

func newValidateSpecsCmd(s streams) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "specs",
		Short: "Validate a spec add payload",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if filePath == "" {
				return errInvalidUsage("missing input file", "pass --file path/to/specs.yaml or --file -")
			}
			var payload specsPayload
			if err := readStructuredInput(s.in, filePath, &payload); err != nil {
				return err
			}
			return writeValidationResult(s.out, validateSpecsPayload(payload))
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "path to a YAML or JSON specs payload, or - for stdin")
	return cmd
}

func newValidatePlanCmd(s streams) *cobra.Command {
	var filePath string
	cmd := &cobra.Command{
		Use:   "plan US-XXX",
		Short: "Validate a spec plan payload",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ref := strings.TrimSpace(args[0])
			if ref == "" {
				return errInvalidUsage("missing spec code", "pass US-XXX as positional argument")
			}
			if filePath == "" {
				return errInvalidUsage("missing input file", "pass --file path/to/plan.yaml or --file -")
			}
			var input domain.PlanInput
			if err := readStructuredInput(s.in, filePath, &input); err != nil {
				return err
			}
			return writeValidationResult(s.out, validatePlanPayload(ref, input))
		},
	}
	cmd.Flags().StringVar(&filePath, "file", "", "path to a YAML or JSON plan payload, or - for stdin")
	return cmd
}

func writeValidationResult(w io.Writer, issues []domain.ValidationIssue) error {
	res := domain.ValidationResult{OK: !hasValidationErrors(issues), Issues: issues}
	if res.Issues == nil {
		res.Issues = []domain.ValidationIssue{}
	}
	if err := iox.WriteOK(w, "validation_result", res); err != nil {
		return iox.NewInternal("encoding validation output", err)
	}
	return nil
}

func validateSpecsPayload(p specsPayload) []domain.ValidationIssue {
	issues := []domain.ValidationIssue{}
	if len(p.Specs) == 0 {
		return appendIssue(issues, "error", "E_SPECS_EMPTY", "specs", "payload must include at least one spec", "expected {specs:[...]}")
	}
	seen := map[string]struct{}{}
	for i, spec := range p.Specs {
		base := fmt.Sprintf("specs[%d]", i)
		if !specCodeRE.MatchString(spec.Code) {
			issues = appendIssue(issues, "error", "E_SPEC_CODE_INVALID", base+".code", "spec code must match US-NNN", "use a zero-padded code such as US-001")
		}
		if _, ok := seen[spec.Code]; spec.Code != "" && ok {
			issues = appendIssue(issues, "error", "E_SPEC_CODE_DUPLICATE", base+".code", "spec code is duplicated", "each spec in the payload must have a unique code")
		}
		seen[spec.Code] = struct{}{}
		if strings.TrimSpace(spec.Title) == "" {
			issues = appendIssue(issues, "error", "E_SPEC_TITLE_EMPTY", base+".title", "spec title is required", "")
		}
		if !epicCodeRE.MatchString(spec.Epic.Code) {
			issues = appendIssue(issues, "error", "E_SPEC_EPIC_INVALID", base+".epic.code", "epic code must match EP-NNN", "assign the spec to an explicit epic")
		}
		if !validPriority(spec.Priority) {
			issues = appendIssue(issues, "error", "E_SPEC_PRIORITY_INVALID", base+".priority", "priority must be HIGH, MEDIUM, or LOW", "")
		}
		if spec.Points <= 0 {
			issues = appendIssue(issues, "error", "E_SPEC_POINTS_INVALID", base+".points", "points must be greater than zero", "")
		}
		if spec.Status == "" {
			issues = appendIssue(issues, "error", "E_SPEC_STATUS_EMPTY", base+".status", "status is required", "use the configured TODO status for new specs")
		}
		body := strings.TrimSpace(spec.Body)
		if body == "" {
			issues = appendIssue(issues, "error", "E_SPEC_BODY_EMPTY", base+".body", "spec body is required", "include user story, Demonstrates, and acceptance criteria")
			continue
		}
		lower := strings.ToLower(body)
		if !strings.Contains(lower, "demonstr") && !strings.Contains(lower, "dimostra") {
			issues = appendIssue(issues, "error", "E_SPEC_DEMONSTRATES_MISSING", base+".body", "spec body must include a concrete Demonstrates section", "state what a reviewer can observe after implementation")
		}
		if !strings.Contains(body, "- [ ]") {
			issues = appendIssue(issues, "error", "E_SPEC_ACCEPTANCE_MISSING", base+".body", "spec body must include checklist acceptance criteria", "add one or more '- [ ]' acceptance criteria")
		}
	}
	for i, spec := range p.Specs {
		for _, dep := range spec.BlockedBy {
			if _, ok := seen[dep]; dep != "" && !ok {
				issues = appendIssue(issues, "warning", "W_SPEC_BLOCKER_UNKNOWN", fmt.Sprintf("specs[%d].blocked_by", i), fmt.Sprintf("blocked_by references %s, which is not in this payload", dep), "ensure the dependency already exists in the backlog")
			}
		}
	}
	return issues
}

func validatePlanPayload(specCode string, input domain.PlanInput) []domain.ValidationIssue {
	issues := []domain.ValidationIssue{}
	if !specCodeRE.MatchString(specCode) {
		issues = appendIssue(issues, "error", "E_PLAN_SPEC_CODE_INVALID", "spec_code", "spec code must match US-NNN", "")
	}
	if strings.TrimSpace(input.PlanBody) == "" {
		issues = appendIssue(issues, "error", "E_PLAN_BODY_EMPTY", "plan_body", "plan body is required", "")
	}
	if len(input.Tasks) == 0 {
		issues = appendIssue(issues, "error", "E_PLAN_TASKS_EMPTY", "tasks", "plan must include at least one task", "")
		return issues
	}
	if len(input.Tasks) > 15 {
		issues = appendIssue(issues, "warning", "W_PLAN_TOO_MANY_TASKS", "tasks", "plan has more than 15 tasks", "consider splitting the spec")
	}
	ids := map[string]int{}
	hasTest := false
	for i, task := range input.Tasks {
		base := fmt.Sprintf("tasks[%d]", i)
		if !taskIDRE.MatchString(task.ID) {
			issues = appendIssue(issues, "error", "E_PLAN_TASK_ID_INVALID", base+".id", "task id must match TASK-NN", "")
		}
		if prev, ok := ids[task.ID]; task.ID != "" && ok {
			issues = appendIssue(issues, "error", "E_PLAN_TASK_ID_DUPLICATE", base+".id", fmt.Sprintf("task id duplicates tasks[%d]", prev), "task ids must be unique")
		}
		ids[task.ID] = i
		if strings.TrimSpace(task.Title) == "" {
			issues = appendIssue(issues, "error", "E_PLAN_TASK_TITLE_EMPTY", base+".title", "task title is required", "")
		}
		if strings.TrimSpace(task.Description) == "" {
			issues = appendIssue(issues, "error", "E_PLAN_TASK_DESCRIPTION_EMPTY", base+".description", "task description is required", "")
		}
		switch task.Type {
		case domain.TaskImpl:
		case domain.TaskTest:
			hasTest = true
		default:
			issues = appendIssue(issues, "error", "E_PLAN_TASK_TYPE_INVALID", base+".type", "task type must be Impl or Test", "")
		}
		if strings.TrimSpace(string(task.Status)) == "" {
			issues = appendIssue(issues, "error", "E_PLAN_TASK_STATUS_EMPTY", base+".status", "task status is required", "use TODO for new tasks")
		}
		if strings.TrimSpace(task.Body) == "" {
			issues = appendIssue(issues, "error", "E_PLAN_TASK_BODY_EMPTY", base+".body", "task body must contain an execution contract", "include objective, allowed changes, steps, verification, done criteria, and blockers")
		} else {
			issues = append(issues, validateTaskContract(base+".body", task.Body)...)
		}
	}
	if !hasTest {
		issues = appendIssue(issues, "error", "E_PLAN_TEST_TASK_MISSING", "tasks", "plan must include at least one Test task", "")
	}
	issues = append(issues, validateTaskDependencies(input.Tasks, ids)...)
	return issues
}

func validateTaskContract(path, body string) []domain.ValidationIssue {
	issues := []domain.ValidationIssue{}
	lower := strings.ToLower(body)
	required := map[string]string{
		"objective": "objective",
		"read":      "context to read",
		"change":    "allowed changes",
		"steps":     "implementation steps",
		"verify":    "verification commands",
		"done":      "done criteria",
		"blocker":   "blockers",
	}
	for token, label := range required {
		if !strings.Contains(lower, token) {
			issues = appendIssue(issues, "warning", "W_PLAN_TASK_CONTRACT_WEAK", path, "task execution contract is missing "+label, "make the contract explicit for smaller models")
		}
	}
	return issues
}

func validateTaskDependencies(tasks []domain.Task, ids map[string]int) []domain.ValidationIssue {
	issues := []domain.ValidationIssue{}
	graph := map[string][]string{}
	for i, task := range tasks {
		for _, dep := range task.Dependencies {
			dep = strings.TrimSpace(dep)
			depIndex, ok := ids[dep]
			if !ok {
				issues = appendIssue(issues, "error", "E_PLAN_TASK_DEP_UNKNOWN", fmt.Sprintf("tasks[%d].dependencies", i), fmt.Sprintf("%s depends on unknown task %s", task.ID, dep), "dependencies must reference tasks in the same plan")
				continue
			}
			if depIndex >= i {
				issues = appendIssue(issues, "error", "E_PLAN_TASK_DEP_FUTURE", fmt.Sprintf("tasks[%d].dependencies", i), fmt.Sprintf("%s depends on %s, which is not earlier in the task list", task.ID, dep), "order tasks by dependency")
			}
			graph[task.ID] = append(graph[task.ID], dep)
		}
	}
	for _, cycle := range findTaskCycles(graph) {
		issues = appendIssue(issues, "error", "E_PLAN_TASK_DEP_CYCLE", "tasks", "task dependency cycle detected: "+strings.Join(cycle, " -> "), "remove the cycle before saving the plan")
	}
	return issues
}

func findTaskCycles(graph map[string][]string) [][]string {
	seen := map[string]bool{}
	stack := map[string]bool{}
	var cycles [][]string
	var visit func(string, []string)
	visit = func(id string, path []string) {
		if stack[id] {
			start := 0
			for i, p := range path {
				if p == id {
					start = i
					break
				}
			}
			cycles = append(cycles, append(path[start:], id))
			return
		}
		if seen[id] {
			return
		}
		seen[id] = true
		stack[id] = true
		for _, dep := range graph[id] {
			visit(dep, append(path, dep))
		}
		stack[id] = false
	}
	keys := make([]string, 0, len(graph))
	for id := range graph {
		keys = append(keys, id)
	}
	sort.Strings(keys)
	for _, id := range keys {
		visit(id, []string{id})
	}
	return cycles
}

func validPriority(p domain.Priority) bool {
	switch p {
	case domain.PriorityHigh, domain.PriorityMedium, domain.PriorityLow:
		return true
	default:
		return false
	}
}

func hasValidationErrors(issues []domain.ValidationIssue) bool {
	for _, issue := range issues {
		if issue.Severity == "error" {
			return true
		}
	}
	return false
}

func appendIssue(issues []domain.ValidationIssue, severity, code, path, message, hint string) []domain.ValidationIssue {
	return append(issues, domain.ValidationIssue{
		Severity: severity,
		Code:     code,
		Path:     path,
		Message:  message,
		Hint:     hint,
	})
}
