# Connector: GitHub Projects

Load this connector only when `.airchetipo/config.yaml` has `backend: github`.

This connector overrides the setup and write-output phases of `airchetipo-spec` while keeping backlog decomposition, story generation, prioritization, and quality rules in the active flow unchanged.

Use it for both:
- initial backlog creation
- incremental story addition

## Common Setup

### Step 1 - Auth and Repository Discovery

Detect repository owner and repository name in one command:

```bash
gh repo view --json owner,name,nameWithOwner --jq '{owner: .owner.login, name: .name, repo: .nameWithOwner}'
```

Save:
- `$OWNER`
- `$REPO_NAME`
- `$REPO_SLUG`

Then discover available projects:

```bash
gh project list --owner "$OWNER" --format json
```

If this fails because of missing scopes, stop and show:

```text
Non ho i permessi necessari per accedere ai GitHub Projects.

Esegui:
gh auth refresh -s read:project -s project

Poi rilancia la richiesta di backlog o user story.
```

### Step 2 - Detect Whether a Backlog Already Exists

Infer which project is linked to the current repository.

A project counts as an existing backlog only when both conditions are true:
1. its title suggests backlog work, preferring exact title `$REPO_NAME Backlog`, then titles containing `Backlog`
2. it already contains issues labeled `airchetipo-backlog` for `$REPO_SLUG`

Behavior:
- if an existing backlog project with backlog issues is found, save it as the active backlog project
- if a project exists but has no backlog issues yet, treat it as an empty backlog container
- if no project exists, create one only when you are creating the initial backlog

Project creation:

```bash
gh project create --owner "$OWNER" --title "$REPO_NAME Backlog"
```

Save the project number as `$PROJECT_NUMBER`.

### Step 3 - Existing-vs-Initial Routing for GitHub

If the active flow is extending the backlog:
- and an existing backlog project with backlog issues is found, continue with incremental story addition
- and no backlog issues are found, do not fail; tell the user that no GitHub backlog exists yet and switch to initial backlog creation

If the active flow is creating the initial backlog:
- use the discovered project if one already exists
- otherwise create it

### Step 4 - Custom Fields and Status Setup

Read existing fields once:

```bash
gh project field-list $PROJECT_NUMBER --owner "$OWNER" --format json
```

Extract:
- `$PROJECT_NODE_ID`
- `$STATUS_FIELD_ID`
- existing status option IDs
- `$PRIORITY_FIELD_ID`
- priority option IDs if present
- `$SP_FIELD_ID`
- `$EPIC_FIELD_ID` if present

Create missing fields when needed:

```bash
gh project field-create $PROJECT_NUMBER --owner "$OWNER" --name "Priority" --data-type "SINGLE_SELECT" --single-select-options "HIGH,MEDIUM,LOW"
gh project field-create $PROJECT_NUMBER --owner "$OWNER" --name "Story Points" --data-type "NUMBER"
```

Ensure the configured workflow statuses exist:

```bash
gh api graphql -f query='mutation {
  updateProjectV2Field(input: {
    projectId: "<PROJECT_NODE_ID>",
    fieldId: "<STATUS_FIELD_ID>",
    name: "Status",
    singleSelectOptions: [
      {name: "{config.workflow.statuses.todo}", color: GRAY},
      {name: "{config.workflow.statuses.planned}", color: BLUE},
      {name: "{config.workflow.statuses.in_progress}", color: YELLOW},
      {name: "{config.workflow.statuses.review}", color: PURPLE},
      {name: "{config.workflow.statuses.done}", color: GREEN}
    ]
  }) {
    projectV2Field {
      ... on ProjectV2SingleSelectField { id options { id name } }
    }
  }
}'
```

Save the returned status option IDs directly from the mutation response.

## Startup Variants

### When creating the initial backlog

Send this before issue creation:

```text
AIRCHETIPO - BACKLOG INITIALIZATION (GitHub Projects)

Emanuele e Andrea sono pronti a costruire il backlog iniziale.

PRD trovato: [file path]
GitHub Project: [project title] (#N)
Owner: [owner]

Analisi dei requisiti in corso...
```

### When extending an existing backlog

Send this after context discovery:

```text
AIRCHETIPO - BACKLOG EXTENSION (GitHub Projects)

Andrea ed Emanuele sono pronti ad aggiungere nuove storie al backlog.

GitHub Project: [project title] (#N)
Backlog esistente: [N issue rilevate]
Prossimo codice disponibile: US-XXX
```

## Initial Backlog Output

Use this section when the active flow is creating the first backlog.

### Step 1 - Idempotency Check

Search for existing backlog issues:

```bash
gh issue list --label "airchetipo-backlog" --state all --json number,title --limit 200
```

If issues already exist, ask whether to:
- skip existing
- recreate
- abort

### Step 2 - Create Labels

Create the shared backlog label and one label per epic in a single shell call.

### Step 3 - Create or Update Epic Field

Create the `Epic` single-select field after the epic list is finalized.
If the field already exists, update it while preserving existing options.

### Step 4 - Create Issues

Create one GitHub issue per story.

Each issue body must contain:
- story
- demonstrates
- acceptance criteria
- epic
- priority
- story points
- blocked by
- scope

Use the `airchetipo-backlog` label plus the epic label.

### Step 5 - Backfill Dependencies

After all issues are created, replace story-code dependencies with actual GitHub issue references for stories that have blockers.

### Step 6 - Collect Node IDs

Fetch node IDs in one GraphQL query for the created issues.

### Step 7 - Add to Project and Set Fields

1. Add all issues to the project via a batched GraphQL mutation
2. Then set all required fields in a second batched GraphQL mutation:
   - Status
   - Priority
   - Story Points
   - Epic

Split the field-update mutation into chunks only if the backlog is very large.

### Final Summary for Initial Backlog

```text
Backlog generated successfully on GitHub Projects.

Project: [project URL]

Summary:
- Epics: N
- User Stories (Issues): N
- Total Story Points: N
- HIGH priority: N stories
- MEDIUM priority: N stories
- LOW priority: N stories
```

Then list the created issues concisely.

## Incremental Story Output

Use this section when the active flow is extending an existing backlog.

### Step 1 - Read Existing Backlog Issues

Read backlog issues once:

```bash
gh issue list --label "airchetipo-backlog" --state all --json number,title,labels,body --limit 200
```

Extract:
- existing epics
- last `US-XXX` code used
- current issue numbers for dependency backfilling

### Step 2 - Create Missing Labels and Epic Options

If the new stories touch a new epic:
- create the missing epic label
- add the missing option to the `Epic` field while preserving existing ones

### Step 3 - Create Only the Confirmed New Issues

For each confirmed story:
- create the issue
- use the same body shape as initial backlog creation
- attach `airchetipo-backlog` and the epic label

### Step 4 - Backfill `Blocked by`

When a new story depends on another story from the same epic, replace story codes with the corresponding GitHub issue references after issue creation.

### Step 5 - Add New Issues to Project and Set Fields

For every newly created issue:
- add it to the project
- set Status, Priority, Story Points, and Epic

### Final Summary for Incremental Stories

```text
Storie aggiunte al backlog GitHub.

Project: [project URL]

Aggiunto:
- #NN US-XXX: [titolo] (EP-XXX | PRIORITY | Npt)
- #NN US-XXX: [titolo] (EP-XXX | PRIORITY | Npt)
```

## Performance Rules

- Minimize API round-trips
- Prefer single `gh` shell calls with loops for label and issue creation
- Prefer batched GraphQL mutations for project item creation and field updates
- Avoid rereading project fields unless the previous mutation response is insufficient
