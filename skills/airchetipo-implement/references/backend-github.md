# Backend: GitHub Projects v2

> Load this file only when `.airchetipo/config.yaml` has `backend: github`.
> This reference overrides the backend-specific I/O phases of the implement skill.
> It does **not** override the core execution contract, autonomy policy, review policy, or fix loop defined in `SKILL.md`.

---

## Setup

### Step 1 - Auth check & current repository detection

1. Detect repository owner:
   ```bash
   gh repo view --json owner --jq '.owner.login'
   ```
   Save as `$OWNER`.
2. Detect the current repository name and slug:
   ```bash
   gh repo view --json name,nameWithOwner --jq '{name: .name, repo: .nameWithOwner}'
   ```
   Save as `$REPO` and `$REPO_SLUG`.
3. Test GitHub Projects auth:
   ```bash
   gh project list --owner "$OWNER" --limit 1 --format json
   ```
4. If Projects auth fails with a missing-scope error, stop and show:

```text
🔧 **Ugo:** Non ho i permessi necessari per accedere ai GitHub Projects.

Esegui questo comando per abilitare lo scope necessario:
- `gh auth refresh -s read:project -s project`

Poi rilancia la skill.
```

### Step 2 - Project discovery

1. Find the project linked to the current git repository:
   ```bash
   gh project list --owner "$OWNER" --format json
   ```
2. Treat a project as linked only if its items contain issues whose `content.repository.nameWithOwner` matches `$REPO_SLUG`.
3. If multiple linked projects exist, prefer:
   - exact title `$REPO Backlog`
   - otherwise a title containing `Backlog`
   - otherwise the linked project with the lowest project number
4. If no linked project is found, fall back to an exact title match `$REPO Backlog`.
5. If still not found, stop and show:

```text
🔧 **Ugo:** Non trovo un GitHub Project collegato al repository corrente.

Esegui prima `airchetipo-inception` chiedendo di generare il backlog dal PRD su GitHub Projects.
```

6. Save the project number and fetch field metadata:
   ```bash
   gh project field-list $PROJECT_NUMBER --owner "$OWNER" --format json
   ```
7. Record field IDs and option IDs for Status, Priority, Story Points, and Epic.

---

## Read Backlog (Story Source)

With `backend: github`, GitHub is the source of truth. The implementation plan lives in:
- the parent issue body for strategic plan and test strategy
- the sub-issues for executable tasks

### Step 3 - Fetch eligible items

1. Fetch all items:
   ```bash
   gh project item-list $PROJECT_NUMBER --owner "$OWNER" --format json -L 200
   ```
2. Filter to items where Status equals `{config.workflow.statuses.planned}`.
3. If no eligible items are found, stop and show:

```text
🔧 **Ugo:** Non ci sono story pronte per l'implementazione.

Per essere implementabile, una story deve essere in stato "{config.workflow.statuses.planned}" nel project.

Puoi:
- Eseguire `/airchetipo-plan` per pianificare una story
- Specificare una story diversa come argomento
```

### Step 4 - Story selection

1. If a story code was passed as argument, search for it among the eligible items by title prefix.
2. If no argument was passed:
   - choose the highest priority eligible item
   - break ties with the lowest story number
3. Read the full issue body:
   ```bash
   gh issue view <NUMBER> --json body,title,labels,number,url
   ```

### Step 4b - Load implementation plan from GitHub

Read both sources:

1. **Strategic plan from parent issue body**
   - Parse `## 📋 Piano di Implementazione`
   - Read the sections describing the technical solution and test strategy

2. **Task list from sub-issues**
   ```bash
   gh api /repos/$OWNER/$REPO/issues/$PARENT_NUMBER/sub_issues \
     -H "X-GitHub-Api-Version: 2026-03-10"
   ```

3. For each open sub-issue, extract when present:
   - stable identity: GitHub issue number
   - task ID from title, if present
   - type from `**Tipo:**`
   - dependencies from `**Dipendenze:**`
   - prose description from the body
   - completion criteria from `**Completamento:**`

4. Build the task list with enough structure to schedule waves when possible.
   - if task identifiers and dependencies are sufficiently reliable, use them for dependency-aware wave planning
   - if they are too weak for confident graph-based wave planning, switch to sequential scheduling

### Validation policy

Do **not** invent certainty. The point of this validation is to prevent a fake dependency graph.

- If `Tipo` is missing but the body clearly describes an implementation task or a test task, infer it and log a warning
- If `Tipo` is missing and the task cannot be classified confidently, treat that task as sequential-only
- If `Dipendenze` is missing or malformed, do **not** assume independent scheduling; treat the task as sequential unless another trustworthy source makes the dependency clear
- If `TASK-XX` is missing, keep the GitHub issue number as the task identity; do **not** invent a new ordered task code when the ordering would affect planning decisions
- If multiple malformed tasks prevent a trustworthy execution order, stop and tell the user that the planning artifacts need repair

### Sequential scheduling for GitHub tasks

Use sequential scheduling when:
- task identity is partially usable but not clean enough for graph scheduling
- dependencies are unclear
- the issue bodies are understandable enough to continue sequentially

In sequential scheduling:
- preserve the GitHub issue number as the stable task identifier
- execute tasks sequentially
- do not claim that tasks are independent unless the evidence is explicit
- tell the user briefly that execution is continuing with sequential scheduling because the GitHub task structure is imperfect

### Hard blocker

If no open sub-issues exist, or the sub-issues are too malformed to reconstruct executable work, stop and show:

```text
🔧 **Ugo:** L'issue #{PARENT_NUMBER} non ha un task plan GitHub sufficientemente leggibile per procedere in sicurezza.

Puoi:
- Eseguire `/airchetipo-plan {US-CODE}` per rigenerare il piano
- Correggere manualmente le sub-issues su GitHub
```

### Step 4c - Load mockup references

1. Scan the parent issue body for a mockup section, `🎨` marker, or explicit paths under `{config.paths.mockups}/`
2. If a mockup directory is referenced, check whether it exists locally and list its contents
3. If mockup files are found, treat them as mandatory references for UI implementation tasks
4. If the directory does not exist or is empty, do not block only for that reason

---

## Status Transition: Move to {config.workflow.statuses.in_progress}

### Step 5 - Move the project item

Update the item status to `{config.workflow.statuses.in_progress}`:

```bash
gh project item-edit --project-id "<PROJECT_NODE_ID>" --id "<ITEM_ID>" --field-id "<STATUS_FIELD_ID>" --single-select-option-id "<IN_PROGRESS_OPTION_ID>"
```

The session announcement should include the issue reference:

```text
**Issue:** #NN - spostata a "{config.workflow.statuses.in_progress}" ✅
```

---

## Task Completion During Implementation

When a GitHub task is completed during Phase 2:
- if the task maps cleanly to a sub-issue, close that sub-issue with `gh issue close <SUB_ISSUE_NUMBER>`
- if execution is running with sequential scheduling and task identity is weaker, close only the sub-issues that are clearly complete; do not guess

---

## Write Output (Completion)

After the core completion gate passes:
- no critical review findings remain
- the full required test suite passes
- only non-blocking improvements may remain

### 1. Run the full required test suite

Do the final verification exactly as required by the core skill.

### 2. Close completed sub-issues

1. List the native sub-issues of the parent:
   ```bash
   gh api /repos/$OWNER/$REPO/issues/$PARENT_NUMBER/sub_issues \
     -H "X-GitHub-Api-Version: 2026-03-10" --jq '.[].number'
   ```
2. Close the sub-issues that correspond to completed work.
3. With sequential scheduling, do **not** close ambiguous sub-issues speculatively.

### 3. Move the story to {config.workflow.statuses.review}

```bash
gh project item-edit --project-id "<PROJECT_NODE_ID>" --id "<ITEM_ID>" --field-id "<STATUS_FIELD_ID>" --single-select-option-id "<REVIEW_OPTION_ID>"
```

Do not move the story to `{config.workflow.statuses.done}` from this skill. That transition remains human-only.

### 4. Post a summary comment on the parent issue

Use a completion comment like. If non-blocking `🟡 MIGLIORAMENTO` items remain open, include them in the comment under an explicit optional improvements section:

```text
## ⚡ Implementazione Completata

**Stato:** {config.workflow.statuses.review}

**Riepilogo:**
- Task completati: {N}/{N}
- Sub-issues chiuse: {N}
- Test scritti/eseguiti: {N}
- Code review: superata ✅
- Cicli di review: {N}

**File creati/modificati:**
- `path/to/new-file.ts`
- `path/to/modified-file.ts`
- `path/to/test-file.test.ts`

**Miglioramenti opzionali rimasti aperti:**
- [Titolo miglioramento] - `path/to/file.ts:NN` - [breve suggerimento]

_Implementato da AIRchetipo Implementation Team_
```

### 5. Update labels if the repository workflow uses them

If the repo already uses labels such as `planned` and `in-review`, update them consistently.
Do not introduce a new label taxonomy unless the existing project workflow clearly expects it.

---

## Technical Reference

### Parsing flow

All `item-edit` commands require node IDs. The normal retrieval flow is:
1. `gh project list --owner "$OWNER" --format json`
2. `gh project field-list $N --owner "$OWNER" --format json`
3. `gh project item-list $N --owner "$OWNER" --format json -L 200`

Prefer `--format json` for machine-parseable output.

### Item list limit

Always use `-L 200` with `gh project item-list` to avoid the default limit of 30 items.

### Status transitions

| From | To | When |
|---|---|---|
| {config.workflow.statuses.planned} | {config.workflow.statuses.in_progress} | Implementation starts |
| {config.workflow.statuses.in_progress} | {config.workflow.statuses.review} | Implementation completes |
| {config.workflow.statuses.review} | {config.workflow.statuses.done} | Human reviewer only |
