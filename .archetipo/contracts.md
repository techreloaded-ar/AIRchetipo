# ARchetipo Connector Contracts

This file is the single entry point for connector operations. Skills read this file to know **how to invoke the CLI** that performs every operation deterministically.

## How It Works

1. Read `.archetipo/config.yaml` to determine the active `connector` (`file` or `github`) and the target paths.
2. Invoke the CLI binary at `.archetipo/bin/archetipo`.
3. Parse the JSON envelope written to stdout. On failure, the JSON envelope on stderr describes the error.

> **Context discipline:** Load this file once at the start of the skill. Do not re-read it unless the skill explicitly requires a refresh.

## Protocol

### Stdout envelope (success)

```json
{"schema":"archetipo/v1","kind":"<kind>","data":{...}}
```

### Stderr envelope (failure)

```json
{"schema":"archetipo/v1","kind":"error","error":{"code":"E_*","message":"...","hint":"..."}}
```

The skill should never branch on `message` (free-text); branch on `code`.

### Exit codes

| Code | Meaning |
|---|---|
| `0` | success |
| `1` | generic error |
| `2` | invalid input (bad flag, malformed stdin JSON) |
| `3` | connector failure (auth, network, gh, fs) |
| `4` | precondition missing (e.g. backlog absent) |

### Error codes (`error.code`)

`E_INVALID_INPUT`, `E_AUTH_SCOPE`, `E_NETWORK`, `E_CONNECTOR`, `E_PRECONDITION`, `E_NOT_FOUND`, `E_CONFLICT`, `E_INTERNAL`.

### Configuration

The CLI reads `.archetipo/config.yaml` from the project root (walks up if invoked from a subdir). Defaults: `connector: file` with the canonical paths.

```yaml
connector: file | github
paths:
  prd: docs/PRD.md
  backlog: docs/BACKLOG.md
  planning: docs/planning/
  mockups: docs/mockups/
  test_results: docs/test-results/
workflow:
  statuses:
    todo: TODO
    planned: PLANNED
    in_progress: IN PROGRESS
    review: REVIEW
    done: DONE
```

---

## Operation Catalog

> Every command emits an envelope with `schema: archetipo/v1`. The `kind` of each envelope is listed below; `data.*` fields follow the schemas in [domain types](#domain-types).

### SETUP

#### `archetipo init` — initialize_connector

Authenticate, detect repo/project, load metadata.

- **Args:** none
- **Stdin:** none
- **Stdout kind:** `setup` — `data` is a `SetupInfo`.
- **Errors:** `E_AUTH_SCOPE` (gh missing scopes), `E_PRECONDITION` (no project linked).

```bash
.archetipo/bin/archetipo init
```

### READ

#### `archetipo backlog list` — fetch_backlog_items

- **Args:** `--status <STATUS>` (optional) filter by workflow status.
- **Stdout kind:** `stories` — `data.items: Story[]`.

#### `archetipo story select` — select_story

- **Args:** `--story US-XXX` (specific story) **or** `--auto` with `--eligible TODO,PLANNED` (comma-separated). `--story` and `--auto` are mutually exclusive; default is auto-select with `--eligible TODO`.
- **Stdout kind:** `story` — `data` is a `Story`.
- **Errors:** `E_PRECONDITION` (no eligible stories or US-XXX not found).

#### `archetipo story read` — read_story_detail

- **Args:** `--ref US-XXX` (required).
- **Stdout kind:** `story`.

#### `archetipo tasks read` — read_story_tasks

- **Args:** `--ref US-XXX` (required, parent story).
- **Stdout kind:** `tasks` — `data.items: Task[]`.
- **Errors:** `E_PRECONDITION` (no plan saved yet).

#### `archetipo backlog existing` — read_existing_backlog

Idempotency metadata for extending an existing backlog.

- **Stdout kind:** `backlog_summary` — `data: BacklogSummary` with `codes`, `last_code`, `epics`, `titles`.

### WRITE

All write operations emit `kind: write_result` with `data: {ok: boolean, refs: Ref[]}`.

#### `archetipo prd save` — save_prd

- **Stdin:** raw markdown body.
- **Errors:** filesystem errors as `E_CONNECTOR`.

#### `archetipo backlog save` — save_initial_backlog

- **Stdin JSON:** `{"stories":[Story, ...]}`.
- **Errors:** `E_CONFLICT` if a non-empty backlog already exists. Use `backlog append` instead.

#### `archetipo backlog append` — append_stories

- **Stdin JSON:** `{"stories":[Story, ...]}`. Stories whose `code` already exists are skipped.

#### `archetipo plan save` — save_plan

- **Args:** `--ref US-XXX` (parent story).
- **Stdin JSON:** `{"plan_body":"<markdown>","tasks":[Task, ...]}`.
- **Effect (file):** writes `{paths.planning}/{US-XXX}.md` with the canonical layout (preamble marker, plan body, tasks marker, GFM table).
- **Effect (github):** appends the plan body to the parent issue, creates one sub-issue per task, links sub-issues to parent.

#### `archetipo status set` — transition_status

- **Args:** `--ref US-XXX --to <STATUS>` (status from `workflow.statuses`).

#### `archetipo task complete` — complete_task

- **Args:** `--parent US-XXX --ref TASK-NN`.

#### `archetipo comment post` — post_comment

- **Args:** `--ref US-XXX`.
- **Stdin:** raw markdown body.
- **Note:** no-op for the file connector (returns `ok: true`).

---

## Domain types

All field names in JSON are `snake_case`.

### `Story`

```jsonc
{
  "code": "US-001",
  "title": "Login utente",
  "epic": {"code": "EP-001", "title": "Auth Foundations"},
  "priority": "HIGH",            // HIGH | MEDIUM | LOW
  "story_points": 3,
  "status": "TODO",              // value from workflow.statuses
  "blocked_by": ["US-002"],      // optional, strings
  "scope": "MVP",                // optional
  "body": "## Story\n\n...",     // markdown body — produced by the skill
  "ref": "US-001",               // connector-local id (issue number for github)
  "url": "https://..."           // populated when the connector has one
}
```

### `Task`

```jsonc
{
  "id": "TASK-01",
  "title": "Schema DB",
  "description": "Create the users schema",
  "type": "Impl",                // Impl | Test
  "status": "TODO",              // value from workflow.statuses
  "dependencies": ["TASK-00"],   // optional
  "body": "...",                 // optional markdown body (read on github)
  "ref": "TASK-01"               // connector-local id (sub-issue number for github)
}
```

### `Ref`

```jsonc
{"code": "US-001", "number": 42, "path": "docs/BACKLOG.md", "url": "https://..."}
```

`number`, `path`, `url` are populated only when the connector has one.

### `SetupInfo`

```jsonc
{
  "connector": "file",
  "paths": { ... },              // mirrors config.yaml paths
  "workflow": { "statuses": { ... } },
  "repo": { "owner": "...", "name": "...", "slug": "owner/name", "node_id": "..." },     // github only
  "project": { "number": 4, "node_id": "...", "url": "...", "fields": { ... } }          // github only
}
```

### `BacklogSummary`

```jsonc
{
  "codes": ["US-001", "US-002"],
  "last_code": "US-002",
  "epics": [{"code": "EP-001", "title": "..."}],
  "titles": ["Login", "Logout"]
}
```

---

## Notes for skill authors

- **Call only what you need.** Not every skill uses every command. Unused commands have zero cost.
- **Content templates belong to the skill, not to the CLI.** The skill produces the markdown body of stories, plans, comments and PRDs. The CLI persists what the skill emits and adds machine-readable markers around it.
- **Branch on error `code`, not on `message`.** The CLI guarantees stable codes; messages are human-readable and may change.
- **No-op operations are explicit.** `comment post` returns `ok: true` even when the file connector has no comment store. Skills never need to suppress those calls.
- **Compose with stdin/stdout.** Every command that takes content reads it from stdin; every command that returns data writes a single JSON envelope to stdout. Pipe and parse.
