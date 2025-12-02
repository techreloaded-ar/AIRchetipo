---
name: implement-story
agent: developer-agent
subtask: true
---

@.opencode/context/core/essential-patterns.md
@.opencode/context/requirements/backlog-format-guide.md
@.opencode/context/development/coding-standards.md

**Implementation Request:** $ARGUMENTS


## Implementation Workflow

### Phase 1: Initialization

#### 1. Parse Backlog and Select Story

**Read backlog structure:**
- Read `docs/backlog.md` to identify epics and stories
- Look for TODO stories: `- [ ] [US-XXX](stories/US-XXX-slug.md)`

**Story Selection:**
- **If user provided story ID** (e.g., "US-005"): Use that story
- **If no story specified**: Auto-select first TODO story in backlog
- **If no TODO stories**: Report "No TODO stories found. All done! 🎉" and exit

#### 2. Read and Validate Story File

**Read story file:**
- Parse `docs/stories/US-XXX-slug.md`
- Extract sections:
  - User Story (As a... I want... So that...)
  - Acceptance Criteria (GHERKIN scenarios)
  - Architecture Notes (technical guidance)
  - Tasks list (implementation tasks)

**Validate story:**
- Ensure story has at least one task
- **If no tasks**: Create task following the backlog-format-guide.md"

#### 3. Ask Execution Mode

Present execution mode choice to user:

```
📋 Ready to implement US-XXX: [Story Title]

Tasks to implement: X tasks

How do you want to proceed?

1. 🚀 YOLO mode (default): Implement all tasks automatically in sequence
2. 🐢 Step-by-step mode: Implement one task at a time, wait for confirmation

Your choice (press Enter for YOLO):
```

**Default behavior:**
- If user presses Enter or doesn't respond: Use YOLO mode
- YOLO mode: Implement all tasks sequentially without stopping
- Step-by-step mode: After each task, ask "Continue to next task? (y/n)"

#### 4. Create Feature Branch

**Branch naming:** `feature/US-XXX-slug`

```bash
git checkout -b feature/US-XXX-slug
```

**If branch exists:**
- Ask user: "Branch feature/US-XXX-slug already exists. Use it? (y/n)"
- If no: Ask for alternative branch name

**Report to user:**
```
✅ Created branch: feature/US-XXX-slug
📍 Ready to implement X tasks
```

---

### Phase 2: Task Implementation Loop

**For each task in the story file:**

#### Step 1: Update Task Status to IN PROGRESS

**Read-Modify-Write:**
1. Read story file content
2. Find task line: `- [ ] TK-XXX: description (role, duration, @assignee)`
3. Update to: `- [~] TK-XXX: description (role, duration, @assignee)`
4. Write story file

**Report to user:**
```
🔨 Starting TK-XXX: [task description]
```

#### Step 2: Analyze Task Requirements

**Gather context:**
- Read task description to understand what needs to be implemented
- Review Architecture Notes for:
  - Components to create/modify
  - APIs and endpoints
  - Data models and schemas
  - Technologies and libraries to use
- Review Acceptance Criteria to understand:
  - Expected behavior (happy path)
  - Error handling requirements
  - Edge cases to cover
- Check PRD (from context) for:
  - Tech stack details
  - Project structure conventions
  - Framework-specific patterns

**Identify work scope:**
- List files to create
- List files to modify
- List dependencies to add (if any)

#### Step 3: Implement Code

**Align with Architecture Notes:**
- Use components/services suggested by @architect-agent
- Use suggested libraries and patterns

**Cover Acceptance Criteria:**
- Ensure implementation satisfies all GHERKIN scenarios
- Handle happy path (normal successful flow)
- Handle validation errors (invalid input, business rules)
- Handle edge cases (boundary conditions, timeouts, null values)

**Implementation approach:**
- Work in small increments
- Implement one feature/component at a time








### Phase 3: Git Commit (After Successful Task)

Follow the Conventional Commits format (Conventional Commits) and create a commit for the new implementation

**Report to user:**
```
💾 Committed: TK-XXX - [brief description]
```

---







### Phase 4: Testing

#### Step 1: Write Tests

**1. Report to user:**
```
🧪 I'm now calling the @tester-agent to write tests for  $ARGUMENTS
```

**2. Execute the command:**
Execute the command `/write-tests $ARGUMENTS`

**3. Run tests**

Run the tests, if you don't know what command to launch, ask the @tester-agent.


### Phase 5: Test Outcome Handling

#### Case A: Tests PASS ✅

**Actions:**
1. **Update task checkbox:**
   - Read story file
   - Find task: `- [~] TK-XXX: ...`
   - Update to: `- [x] TK-XXX: ... ✅ YYYY-MM-DD` (use current date)
   - Write story file

2. **Append to Dev Notes:**
   - Find `## Dev Notes` section
   - If section contains only `_(Sezione da compilare in sviluppo)_`, replace with:
   ```markdown
   ## Dev Notes

   ### TK-XXX Implementation (YYYY-MM-DD)

   **What was done:**
   - Brief description of implementation (1-3 bullet points)
   ```
   - If section already has content, append the new entry

3. **Report to user:**
   ```
   ✅ TK-XXX completed successfully
   📁 Files: file1.ts, file2.ts
   🧪 Tests: All passing
   ```

4. **If step-by-step mode:**
   - Ask user: "Continue to next task? (y/n)"
   - Wait for response
   - If "n": Stop and report "Paused. Run `/implement-story US-XXX` again to continue."

5. **Proceed to next task**

---

#### Case B: Tests FAIL ❌ (First Time)

**Actions:**
1. **Keep task as IN PROGRESS:**
   - Task remains: `- [~] TK-XXX: ...`

2. **Append to Dev Notes:**
   ```markdown
   ### TK-XXX Implementation Attempt (YYYY-MM-DD)

   **Status:** ❌ Tests failing

   **What was implemented:**
   - Brief description of what was done

   **Files changed:**
   - path/to/file1.ext (created/modified)

   **Error output:**
   ```
   <paste full test output here>
   ```

   **Fix attempted:**
   - Description of what auto-fix tried

   **User guidance needed.**
   ```

3. **Report to user:**
   ```
   ❌ Task TK-XXX: Tests failing after auto-fix attempt

   Error summary: [brief description of error]

   Full test output has been logged in Dev Notes section of the story file.

   How would you like to proceed?
   1. Let me try a different implementation approach
   2. You'll fix it manually (I'll move to next task)
   3. Skip this task for now (mark as blocked)

   Your choice:
   ```

4. **Wait for user decision:**

   **Choice 1 - Try different approach:**
   - Ask user: "Please describe the alternative approach you'd like me to try:"
   - Implement based on user guidance
   - Re-run tests
   - Continue based on outcome

   **Choice 2 - Manual fix:**
   - Report: "Moving to next task. You can fix TK-XXX manually and commit."
   - Proceed to next task

   **Choice 3 - Skip/Block:**
   - Update task: `- [!] TK-XXX: ...` (blocked)
   - Append to Dev Notes: "**Status:** ⚠️ Blocked - awaiting resolution"
   - Proceed to next task

---


### Phase 4: Story Completion

**Trigger:** All tasks in story have `[x]` checkbox (DONE)

#### Actions:

**1. Update Story File Status:**
- Read story file
- Find metadata line: `**Epic:** EP-XXX | **Priority:** HIGH | **Estimate:** 5pt | **Status:** TODO`
- Update Status: `TODO` → `DONE` (or `IN PROGRESS` → `DONE`)
- Write story file

**2. Update Backlog Index:**
- Read `docs/backlog.md`
- Find story line: `- [ ] [US-XXX](stories/US-XXX-slug.md) - Story title | **HIGH** | 5pt`
- Update checkbox: `- [ ]` → `- [x]`
- Write `docs/backlog.md`

**3. Commit Backlog Updates:**
```
chore(US-XXX): Mark story as DONE in backlog
```

**Report to user:**
```
✅ Story US-XXX: All tasks completed!
📊 Summary:
   - Tasks implemented: X
   - Commits: X
```

---

## Error Handling

### Error Categories and Responses

#### 1. Story Not Found
```
❌ Error: Story US-XXX not found in docs/stories/

Available TODO stories:
- US-001: Story title 1
- US-002: Story title 2
- US-005: Story title 5

Please specify which story to implement.
```

#### 2. No TODO Stories
```
🎉 No TODO stories found in backlog. All done!

If you need to implement a specific story, specify its ID:
/implement-story US-XXX
```

#### 3. Story Has No Tasks
```
❌ Error: Story US-XXX has no tasks defined

Please add tasks to the story before implementing.
You can edit: docs/stories/US-XXX-slug.md
```

#### 4. Test Framework Not Detected
```
⚠️ I couldn't auto-detect the test command for this project.

Searched in:
- package.json, pom.xml, build.gradle, Cargo.toml, go.mod, etc.
- README.md, CONTRIBUTING.md, Makefile

Please specify how to run tests:
Examples: "npm test", "pytest", "gradle test", "make test"

Test command:
```

#### 5. Git Operation Failed
```
❌ Git operation failed: <error message>

Please resolve this manually and then:
- Continue: /implement-story US-XXX (will resume from where it stopped)
- Or fix git issue and retry
```

#### 6. Branch Already Exists
```
⚠️ Branch feature/US-XXX-slug already exists

Options:
1. Use existing branch (y)
2. Specify different branch name (n)

Your choice:
```

#### 7. File Write Error
```
❌ Couldn't update file <path>: <error>

Retrying once...

<If retry fails>
❌ File write failed after retry. Please check file permissions.
```


