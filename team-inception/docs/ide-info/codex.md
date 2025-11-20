# AIRchetipo Method - Codex Instructions

## Activating Agents

AIRchetipo agents, tasks and workflows are installed as custom prompts in
`$CODEX_HOME/prompts/air-*.md` files. If `CODEX_HOME` is not set, it
defaults to `$HOME/.codex/`.

### Examples

```
/air-aim-agents-dev - Activate development agent
/air-aim-agents-architect - Activate architect agent
/air-aim-workflows-dev-story - Execute dev-story workflow
```

### Notes

Prompts are autocompleted when you type /
Agent remains active for the conversation
Start a new conversation to switch agents
