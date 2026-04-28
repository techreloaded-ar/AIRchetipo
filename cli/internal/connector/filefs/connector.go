package filefs

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/techreloaded-ar/ARchetipo/cli/internal/config"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/connector"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/domain"
	"github.com/techreloaded-ar/ARchetipo/cli/internal/iox"
)

// Connector is the file-system implementation of connector.Connector.
type Connector struct {
	cfg config.Config
}

// New constructs a Connector. Always succeeds — config validation happens at
// load time.
func New(cfg config.Config) *Connector { return &Connector{cfg: cfg} }

// Register hooks the file connector into the registry under the canonical
// name "file".
func Register() {
	connector.Register(config.ConnectorFile, func(cfg config.Config) (connector.Connector, error) {
		return New(cfg), nil
	})
}

// SETUP

func (c *Connector) InitializeConnector(ctx context.Context) (domain.SetupInfo, error) {
	return domain.SetupInfo{
		Connector: config.ConnectorFile,
		Paths:     c.cfg.Paths,
		Workflow:  c.cfg.Workflow,
	}, nil
}

// READ

func (c *Connector) FetchBacklogItems(ctx context.Context, statusFilter domain.Status) ([]domain.Story, error) {
	stories, err := c.readBacklog()
	if err != nil {
		return nil, err
	}
	if statusFilter == "" {
		return stories, nil
	}
	out := make([]domain.Story, 0, len(stories))
	for _, s := range stories {
		if s.Status == statusFilter {
			out = append(out, s)
		}
	}
	return out, nil
}

func (c *Connector) SelectStory(ctx context.Context, q domain.SelectQuery) (domain.Story, error) {
	stories, err := c.readBacklog()
	if err != nil {
		return domain.Story{}, err
	}
	if q.StoryCode != "" {
		for _, s := range stories {
			if s.Code == q.StoryCode {
				return s, nil
			}
		}
		return domain.Story{}, iox.NewPrecondition(
			fmt.Sprintf("story %s not found in backlog", q.StoryCode),
			"check the backlog or run `archetipo backlog list`", nil)
	}
	eligible := map[domain.Status]struct{}{}
	for _, st := range q.EligibleStatuses {
		eligible[st] = struct{}{}
	}
	candidates := make([]domain.Story, 0, len(stories))
	for _, s := range stories {
		if _, ok := eligible[s.Status]; ok {
			candidates = append(candidates, s)
		}
	}
	if len(candidates) == 0 {
		return domain.Story{}, iox.NewPrecondition(
			"no eligible stories in backlog",
			"check --eligible or status of stories", nil)
	}
	sortByPriorityThenCode(candidates)
	return candidates[0], nil
}

func (c *Connector) ReadStoryDetail(ctx context.Context, ref string) (domain.Story, error) {
	stories, err := c.readBacklog()
	if err != nil {
		return domain.Story{}, err
	}
	for _, s := range stories {
		if s.Code == ref {
			return s, nil
		}
	}
	return domain.Story{}, iox.NewPrecondition(
		fmt.Sprintf("story %s not found in backlog", ref), "", nil)
}

func (c *Connector) ReadStoryTasks(ctx context.Context, parentRef string) ([]domain.Task, error) {
	planPath := c.planPath(parentRef)
	raw, err := os.ReadFile(planPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, iox.NewPrecondition(
				fmt.Sprintf("planning file for %s not found", parentRef),
				"run `archetipo plan save` first", err)
		}
		return nil, fmt.Errorf("reading plan: %w", err)
	}
	_, tasks, err := parsePlan(string(raw))
	return tasks, err
}

func (c *Connector) ReadExistingBacklog(ctx context.Context) (domain.BacklogSummary, error) {
	stories, err := c.readBacklog()
	if err != nil {
		return domain.BacklogSummary{}, err
	}
	out := domain.BacklogSummary{}
	seenEpics := map[string]domain.Epic{}
	for _, s := range stories {
		out.Codes = append(out.Codes, s.Code)
		out.Titles = append(out.Titles, s.Title)
		if s.Epic.Code != "" {
			if existing, ok := seenEpics[s.Epic.Code]; !ok || (existing.Title == "" && s.Epic.Title != "") {
				seenEpics[s.Epic.Code] = s.Epic
			}
		}
	}
	sort.Strings(out.Codes)
	if len(out.Codes) > 0 {
		out.LastCode = highestCode(out.Codes)
	}
	for _, e := range seenEpics {
		out.Epics = append(out.Epics, e)
	}
	sort.Slice(out.Epics, func(i, j int) bool { return out.Epics[i].Code < out.Epics[j].Code })
	return out, nil
}

// WRITE

func (c *Connector) SavePRD(ctx context.Context, content string) (domain.WriteResult, error) {
	path := c.cfg.AbsPath(c.cfg.Paths.PRD)
	if err := writeFile(path, []byte(content)); err != nil {
		return domain.WriteResult{}, err
	}
	return domain.WriteResult{OK: true, Refs: []domain.Ref{{Path: path}}}, nil
}

func (c *Connector) SaveInitialBacklog(ctx context.Context, stories []domain.Story) (domain.WriteResult, error) {
	if len(stories) == 0 {
		return domain.WriteResult{}, iox.NewInvalidInput(
			"no stories to write", "stdin must contain a non-empty stories array", nil)
	}
	path := c.cfg.AbsPath(c.cfg.Paths.Backlog)
	if exists(path) {
		// Idempotency: if the existing backlog has stories, refuse to overwrite.
		existing, _ := c.readBacklog()
		if len(existing) > 0 {
			return domain.WriteResult{}, iox.NewConnector(iox.CodeConflict,
				"backlog already exists with stories",
				"use `archetipo backlog append` to add to it, or remove the file to recreate", nil)
		}
	}
	content := renderBacklog(stories)
	if err := writeFile(path, []byte(content)); err != nil {
		return domain.WriteResult{}, err
	}
	return domain.WriteResult{OK: true, Refs: refsFromStories(stories, path)}, nil
}

func (c *Connector) AppendStories(ctx context.Context, stories []domain.Story) (domain.WriteResult, error) {
	if len(stories) == 0 {
		return domain.WriteResult{}, iox.NewInvalidInput(
			"no stories to append", "stdin must contain a non-empty stories array", nil)
	}
	path := c.cfg.AbsPath(c.cfg.Paths.Backlog)
	existing, err := c.readBacklog()
	if err != nil && !errors.Is(err, errBacklogMissing) {
		return domain.WriteResult{}, err
	}
	// Skip stories whose code already exists.
	known := map[string]struct{}{}
	for _, s := range existing {
		known[s.Code] = struct{}{}
	}
	added := make([]domain.Story, 0, len(stories))
	for _, s := range stories {
		if _, ok := known[s.Code]; ok {
			continue
		}
		added = append(added, s)
	}
	merged := append(existing, added...)
	content := renderBacklog(merged)
	if err := writeFile(path, []byte(content)); err != nil {
		return domain.WriteResult{}, err
	}
	return domain.WriteResult{OK: true, Refs: refsFromStories(added, path)}, nil
}

func (c *Connector) SavePlan(ctx context.Context, storyRef string, plan domain.PlanInput) (domain.WriteResult, error) {
	if storyRef == "" {
		return domain.WriteResult{}, iox.NewInvalidInput(
			"missing story ref", "pass --ref US-XXX", nil)
	}
	path := c.planPath(storyRef)
	content := renderPlan(storyRef, plan)
	if err := writeFile(path, []byte(content)); err != nil {
		return domain.WriteResult{}, err
	}
	refs := []domain.Ref{{Code: storyRef, Path: path}}
	for _, t := range plan.Tasks {
		refs = append(refs, domain.Ref{Code: t.ID})
	}
	return domain.WriteResult{OK: true, Refs: refs}, nil
}

func (c *Connector) TransitionStatus(ctx context.Context, storyRef string, newStatus domain.Status) (domain.WriteResult, error) {
	stories, err := c.readBacklog()
	if err != nil {
		return domain.WriteResult{}, err
	}
	idx := -1
	for i := range stories {
		if stories[i].Code == storyRef {
			idx = i
			break
		}
	}
	if idx == -1 {
		return domain.WriteResult{}, iox.NewPrecondition(
			fmt.Sprintf("story %s not found", storyRef), "", nil)
	}
	stories[idx].Status = newStatus
	path := c.cfg.AbsPath(c.cfg.Paths.Backlog)
	if err := writeFile(path, []byte(renderBacklog(stories))); err != nil {
		return domain.WriteResult{}, err
	}
	return domain.WriteResult{OK: true, Refs: []domain.Ref{{Code: storyRef, Path: path}}}, nil
}

func (c *Connector) CompleteTask(ctx context.Context, parentRef, taskRef string) (domain.WriteResult, error) {
	if parentRef == "" || taskRef == "" {
		return domain.WriteResult{}, iox.NewInvalidInput(
			"missing parent or task ref", "pass --parent US-XXX --ref TASK-NN", nil)
	}
	planPath := c.planPath(parentRef)
	raw, err := os.ReadFile(planPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return domain.WriteResult{}, iox.NewPrecondition(
				fmt.Sprintf("planning file for %s not found", parentRef),
				"run `archetipo plan save` first", err)
		}
		return domain.WriteResult{}, fmt.Errorf("reading plan: %w", err)
	}
	body, tasks, err := parsePlan(string(raw))
	if err != nil {
		return domain.WriteResult{}, err
	}
	hit := false
	for i := range tasks {
		if tasks[i].ID == taskRef {
			tasks[i].Status = domain.StatusDone
			hit = true
			break
		}
	}
	if !hit {
		return domain.WriteResult{}, iox.NewPrecondition(
			fmt.Sprintf("task %s not found in plan %s", taskRef, parentRef), "", nil)
	}
	updated := renderPlan(parentRef, domain.PlanInput{PlanBody: body, Tasks: tasks})
	if err := writeFile(planPath, []byte(updated)); err != nil {
		return domain.WriteResult{}, err
	}
	return domain.WriteResult{OK: true, Refs: []domain.Ref{{Code: taskRef, Path: planPath}}}, nil
}

func (c *Connector) PostComment(ctx context.Context, storyRef, body string) (domain.WriteResult, error) {
	// File connector: comments are not modeled. Skill should skip; return ok.
	return domain.WriteResult{OK: true}, nil
}

// helpers

var errBacklogMissing = errors.New("backlog missing")

func (c *Connector) readBacklog() ([]domain.Story, error) {
	path := c.cfg.AbsPath(c.cfg.Paths.Backlog)
	raw, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return nil, iox.NewPrecondition(
				fmt.Sprintf("backlog not found at %s", path),
				"run `archetipo backlog save` first or `archetipo-spec` skill", errBacklogMissing)
		}
		return nil, fmt.Errorf("reading backlog: %w", err)
	}
	return parseBacklog(string(raw))
}

func (c *Connector) planPath(storyRef string) string {
	dir := c.cfg.AbsPath(c.cfg.Paths.Planning)
	return filepath.Join(dir, storyRef+".md")
}

func writeFile(path string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("creating dir: %w", err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}

func exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func refsFromStories(stories []domain.Story, path string) []domain.Ref {
	out := make([]domain.Ref, 0, len(stories))
	for _, s := range stories {
		out = append(out, domain.Ref{Code: s.Code, Path: path})
	}
	return out
}

func sortByPriorityThenCode(s []domain.Story) {
	rank := map[domain.Priority]int{
		domain.PriorityHigh:   0,
		domain.PriorityMedium: 1,
		domain.PriorityLow:    2,
	}
	sort.SliceStable(s, func(i, j int) bool {
		ri, rj := rank[s[i].Priority], rank[s[j].Priority]
		if ri != rj {
			return ri < rj
		}
		return numericTail(s[i].Code) < numericTail(s[j].Code)
	})
}

func numericTail(code string) int {
	idx := strings.LastIndex(code, "-")
	if idx == -1 || idx == len(code)-1 {
		return 0
	}
	n, err := strconv.Atoi(code[idx+1:])
	if err != nil {
		return 0
	}
	return n
}

// highestCode returns the lexically last US-XXX code (assumes zero-padded).
func highestCode(codes []string) string {
	max := ""
	maxN := -1
	for _, c := range codes {
		n := numericTail(c)
		if n > maxN {
			maxN = n
			max = c
		}
	}
	return max
}
