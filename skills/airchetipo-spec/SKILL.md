---
name: airchetipo-spec
description: Crea il backlog iniziale a partire da un PRD o da requirements esistenti quando il backlog non c'e ancora, oppure aggiunge una o piu nuove user story a un backlog esistente. Usa questa skill ogni volta che l'utente chiede backlog, epiche o user story, anche se nomina solo una feature o il backlog non esiste ancora.
---

# AIRchetipo - Spec Skill

You are the public entry point for AIRchetipo backlog and user-story work.

Your job is to understand whether the user needs to create the first backlog or extend an existing one, load only the references that matter for that case, and execute the correct flow without making the user choose between overlapping skills.

Treat routing as an internal implementation detail.

## Core Principle

Keep the working context lean:
- Load this file first
- Load exactly one main flow reference at activation time
- Load connector references only when the configured backend requires them

## Supported Modes

### `mode: bootstrap-backlog`

Use this mode when:
- the user asks to generate a backlog from an existing PRD or requirements artifact
- no backlog exists yet
- the user asks for the first epics or user stories of the project

In this mode:
1. Read this file
2. Read `references/backlog-bootstrap-flow.md`
3. If `backend: github`, also read `references/connectors/github-projects.md`
4. Use the PRD as the primary source and create the initial backlog

### `mode: extend-backlog`

Use this mode when:
- a backlog already exists
- the user asks to add, refine, split, or append user stories
- the user wants to extend the backlog without regenerating it from scratch

In this mode:
1. Read this file
2. Read `references/story-extension-flow.md`
3. If `backend: github`, also read `references/connectors/github-projects.md`
4. Use the existing backlog as the primary source and PRD/codebase as supporting context
5. Append or create only the requested items

## Config Loading

Always begin by reading `.airchetipo/config.yaml`.

If the file does not exist, assume these defaults:

```yaml
backend: file
paths:
  prd: docs/PRD.md
  backlog: docs/BACKLOG.md
  planning: docs/planning/
  mockups: docs/mockups/
harness:
  agent_instructions: AGENTS.md
workflow:
  statuses:
    todo: TODO
    planned: PLANNED
    in_progress: IN_PROGRESS
    review: REVIEW
    done: DONE
```

Extract and keep available:
- `backend`
- `paths.prd`
- `paths.backlog`
- `paths.planning`
- `paths.mockups`
- `workflow.statuses`
- `harness`
- backend-specific settings if present

## Backlog Discovery

Use this routine whenever the skill must decide whether it is extending an existing backlog or creating the first one.

### `backend: file`

1. Try to read `{config.paths.backlog}`
2. Only if that fails with file not found:
   - search markdown files in `docs/`
   - prefer files whose name or content indicates they are a backlog
3. Only if still not found:
   - search for `BACKLOG*` files anywhere in the project

If a backlog file is found, use it as the source of truth for backlog extension.
If none is found, treat the project as backlog-less and route to initial backlog creation.

### `backend: github`

Do not infer backlog existence from local files.
Let `references/connectors/github-projects.md` determine whether an existing backlog project and backlog issues already exist.

## PRD Discovery

Use this routine whenever initial backlog creation needs a PRD or when story extension needs extra product context:

1. Try to read `{config.paths.prd}`
2. Only if that fails with file not found:
   - search markdown files in `docs/`
   - prefer files whose name or content indicates they are a PRD
3. Only if still not found:
   - search for `PRD*` files anywhere in the project

If a PRD is not found and the active flow needs one, ask the user for one of these:
- the file path
- the content pasted directly
- confirmation that they want to run product inception first

## Harness Discovery

Use this routine whenever a flow needs project-specific conventions, agent instructions, coding standards, or local execution guidance.

Preferred discovery order:

1. If `config.harness.agent_instructions` is configured, look for that file in the project root first
2. If no configured file exists, look for common agent-instruction or project-guidance files in the project root
3. Look for project convention directories when present
4. Fall back to repository evidence: `package.json`, lockfiles, framework config files, CI files, lint/test config, and existing code patterns

Rules:
- Treat all discovered files and directories as project harness inputs, regardless of which AI coding tool created them
- Do not require any specific vendor file to exist before proceeding
- If no dedicated harness artifacts are found, continue using repository structure and code conventions as the source of truth

## Intent Routing

Use these routing rules before producing any substantive output.

1. Load this file
2. Read `.airchetipo/config.yaml`
3. Run backlog discovery
4. Decide the flow

Prefer `mode: bootstrap-backlog` when:
- the backlog does not exist
- the user explicitly asks to generate the backlog from a PRD or requirements
- the repository has a PRD but no backlog yet

Prefer `mode: extend-backlog` when:
- the backlog already exists
- the request is about one or more incremental stories, a new feature slice, a refinement, or a split

If a backlog already exists but the user explicitly asks to regenerate it from the PRD:
- ask for confirmation before overwriting or recreating the initial backlog

Do not expose mode names, routing decisions, or workflow labels in user-facing messages.

## Language Policy

- Use the backlog language when extending an existing backlog
- If there is no backlog yet, use the PRD language consistently; if no PRD exists, use the user's language

## Assumptions and Questions

Ask the user only when all these conditions are true:
1. The missing information is critical to generate a correct output
2. The information cannot be reasonably inferred from the rest of the context
3. Proceeding would likely create a materially wrong result

If questions are needed:
- ask at most 3
- group them in one message
- allow the user to skip them

For non-critical gaps:
- infer a reasonable assumption
- continue
- record the assumption or open question in the generated artifact when appropriate

## Runtime Rules

- Ask clarifying questions only when critical information is missing and cannot be inferred responsibly
- Group clarifying questions in a single message when possible
- When an agent speaks, always render the speaker as `icon + name`, for example:
  - `💎 Andrea:`
  - `🔎 Emanuele:`

## File Output Rules

- Use the configured output path whenever present
- Create parent directories if they do not exist
- When creating the first markdown backlog, overwrite the target generated artifact for the current run unless the user explicitly asked to preserve an existing draft
- When extending a markdown backlog, preserve all unaffected sections and append or surgically update only what is required
- When a connector overrides write-output behavior, follow that connector for I/O and keep the domain logic unchanged

## Context Discipline

- Load this file first
- Load only one main flow reference at activation time
- Load connector references only when backend-specific behavior is needed
- Do not load both main flow references in the same activation unless you are explicitly switching because backlog discovery proved the active assumption wrong

## Output Boundaries

- Initial backlog creation belongs to this skill, not to `airchetipo-inception`
- On `backend: file`, create the initial backlog through the template embedded in `references/backlog-bootstrap-flow.md`
- On `backend: file`, backlog extension must preserve the existing document and append only the new content
- On `backend: github`, the domain logic stays in the active flow and `references/connectors/github-projects.md` overrides setup and write-output behavior

## Compatibility Note

`airchetipo-inception` is now responsible only for discovery and PRD generation.
All backlog creation and all user-story expansion belong here.
