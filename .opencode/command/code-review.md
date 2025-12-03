---
name: code-review
agent: developer-agent
---

@.opencode/context/development/coding-standards.md


## Code Review Workflow

### Phase 1: Determine Review Scope

#### 1. Identify Current Branch and Review Mode

**Get current branch:**
- Run `git rev-parse --abbrev-ref HEAD` to get the current branch name
- Store the branch name for later use

**Determine review mode:**

**Mode A: User Story Branch Review**
- **Condition:** Branch matches pattern `feature/*`, `bugfix/*`, or contains a story ID (e.g., `US-XXX`)
- **Scope:** All changes on this branch compared to main
- **Command:** `git diff main...HEAD`
- **Additional context:** Try to identify the user story from branch name and `docs/backlog.md`

**Mode B: Main Branch Review (Pending Changes)**
- **Condition:** Current branch is `main` or `master`
- **Scope:** All uncommitted changes (staged, unstaged, untracked)
- **Commands:**
  - `git diff HEAD` for tracked file changes
  - `git status --short` for all pending files
  - Read untracked files directly if they exist

**Mode C: Other Branch Review**
- **Condition:** Any other branch that doesn't match story patterns
- **Scope:** Same as Mode B (pending uncommitted changes)
- **Commands:** Same as Mode B

#### 2. Collect Changed Files

**For Mode A (Story Branch):**
- Run `git diff main...HEAD --name-status` to get list of changed files
- Categorize files by status:
  - `M` = Modified
  - `A` = Added (new files)
  - `D` = Deleted
  - `R` = Renamed
- Store file list for detailed analysis

**For Mode B & C (Pending Changes):**
- Run `git status --short` to get status of all files
- Parse output:
  - `M ` = Modified (staged)
  - ` M` = Modified (unstaged)
  - `A ` = Added (staged)
  - `??` = Untracked
  - `D ` = Deleted
- For untracked files, read them directly using file read operations
- For modified files, use `git diff HEAD` to get changes

---

### Phase 2: Detailed Code Analysis

#### 4. Review Each Changed File

**For each file in the changed files list:**

**Read the full diff or file content:**
- If Mode A: use `git diff main...HEAD -- <filepath>` for specific file diff
- If Mode B/C: use `git diff HEAD -- <filepath>` or read untracked files directly

**Analyze for the following aspects:**

**A. Security Issues** 🔐
- XSS vulnerabilities (unsanitized user input, innerHTML usage, dangerouslySetInnerHTML)
- CSRF protection (missing tokens, improper session handling)
- SQL Injection (raw queries, improper parameterization, Prisma misuse)
- Authentication/Authorization bypass (missing auth checks, insecure session management)
- Sensitive data exposure (secrets in code, logging sensitive info, weak encryption)
- NextAuth session security (improper callbacks, session token exposure)
- Input validation gaps (missing or weak validation, unvalidated file uploads)
- File upload security (unrestricted file types, path traversal, size limits)

**B. Next.js Best Practices** ⚛️
- Server/Client Component usage (check `'use client'` directive placement)
- API Routes implementation (proper HTTP methods, Edge Runtime compatibility)
- Metadata and SEO (missing or incorrect metadata, og:tags, title tags)
- Data fetching patterns (Server Components vs Client Components, caching strategy)
- Route structure (App Router conventions, dynamic routes, route groups)
- Loading and error states (loading.tsx, error.tsx boundaries)
- Image optimization (next/image usage vs raw img tags)
- Font optimization (next/font usage)

**C. Database & Prisma Issues** 🗄️
- N+1 query problems (missing includes, sequential queries in loops)
- Include/select optimization (over-fetching, missing select clauses)
- Transaction handling (missing transactions for related writes, improper rollback)
- Connection management (connection leaks, missing disconnect)
- Query efficiency (missing indexes implied by queries, full table scans)
- Relation handling (missing cascading deletes, orphaned records)

**D. React & TypeScript Quality** 📝
- React Hooks rules (hooks in loops/conditions, dependency arrays)
- TypeScript type safety (any usage, type assertions, missing types)
- Component structure (prop drilling, missing composition, large components)
- State management (unnecessary state, derived state, stale closures)
- Effect dependencies (missing dependencies, infinite loops)
- Code organization (Clean Architecture adherence, separation of concerns)
- Naming conventions (preferire nomi italiani chiari, meaningful names)

**E. Performance Issues** ⚡
- Unnecessary re-renders (missing React.memo, useMemo, useCallback)
- Large bundle sizes (unoptimized imports, unnecessary dependencies)
- Inefficient algorithms (O(n²) where O(n) is possible)
- Memory leaks (uncleared intervals/timeouts, event listener cleanup)
- Leaflet map rendering (unnecessary map redraws, marker clustering)
- GPX file parsing efficiency (large file handling, streaming vs loading)

**F. Domain-Specific Issues** 🏍️
*Note: This is a Next.js + Prisma + NextAuth project (RideAtlas - motorcycle trip management)*
- GPX file handling (validation, parsing errors, coordinate precision)
- Trip approval workflow (state transitions, authorization checks)
- Media upload handling (image optimization, storage limits, MIME types)
- Geospatial data handling (coordinate validation, bounding box queries)

**G. General Code Quality** ✨
- Code duplication (repeated logic, missing abstractions)
- Function complexity (functions > 50 lines, deep nesting, cyclomatic complexity)
- Naming clarity (nomi italiani preferiti, descriptive variable/function names)
- Error handling completeness (try-catch coverage, error boundaries, user feedback)
- Loading states (missing spinners, skeleton screens)
- User feedback (success/error toasts, form validation messages)

**H. Testing Gaps** 🧪
- Missing unit tests for new business logic
- Missing integration tests for API routes
- Fragile tests (hardcoded IDs, timing dependencies)
- Edge cases not covered (null/undefined, empty arrays, boundary values)

#### 5. Cross-File Analysis

**Check for architectural issues:**
- Circular dependencies
- Improper layer violations (UI calling DB directly, skipping service layer)
- Inconsistent patterns across similar files
- Breaking changes that affect other parts of the codebase

---

### Phase 3: Generate Review Report

#### 6. Create Code Review Report

**Display a concise, scannable report in chat with the following structure:**

---

**� Code Review Complete**

**📊 Stats:** X files modified, +XXX/-XXX lines  
**Assessment:** <🟢 Excellent | 🟡 Needs Review | 🔴 Critical Issues>

<1-2 sentence overall summary>

---

**🔴 Critical Issues** (X found)  
<If none: "✅ None detected">

<For each critical issue - keep to 3-5 max, most important ones:>
**X.** `File:Line` - **Brief Title**  
   - **Issue:** <1 sentence description>  
   - **Fix:** <1 sentence solution>  
   - **Risk:** <Why this matters in 1 sentence>

---

**🟡 Important Issues** (X found)  
<If none: "✅ None detected">

<For each important issue - keep to 3-5 max:>
**X.** `File:Line` - **Brief Title**  
   - **Impact:** <Performance/Maintainability/Scalability>  
   - **Suggestion:** <1 sentence recommendation>

---

**🔵 Minor Suggestions** (X found)  
<If none: "✅ None detected">

<List format - max 5 items:>
- `File:Line` - Brief suggestion
- `File:Line` - Brief suggestion

---

**✅ Strengths**

<List 2-3 specific positive observations:>
- ✅ <What was done well>
- ✅ <Another good thing>

---

**⚠️ Pre-Commit Checks**

<Only show items that FAILED or need attention. If all pass, show:>
✅ All pre-commit checks passed

<Otherwise, show only failed items from:>
- [ ] Unit tests updated
- [ ] No build warnings
- [ ] No debug code
- [ ] Input validation complete
- [ ] Error handling present
- [ ] TypeScript strict mode
- [ ] Server/Client Components correct

---

**🎯 Next Steps**

<2-3 most important actionable recommendations>

---

**Guidelines for the report:**
- **Be concise:** Each issue description should be 1-2 sentences max
- **Prioritize:** Show max 5 critical, 5 important, 5 minor (most important ones)
- **Be specific:** Always include file path and line number
- **Be actionable:** Every issue should have a clear fix/suggestion
- **Skip verbosity:** Avoid repeating section headers if empty, just use checkmarks
- **Code examples:** Only include for critical security issues, and keep them to 3-5 lines max



### Phase 4: Interactive Feedback


#### 7. Fix Issues

**IF the user chose to fix issues**, **ALWAYS** make sure that the tests are passing before moving to the next issue and the project is building.




```
What do you want to do?

1. 🔧 Start fixing issues
2. ❌ Nothing, just wanted the review

Choice (default 4):
```

**Handle user choice:**
- **Choice 1:** Start fixing issues based on the priority list
- **Choice 2:** End the command gracefully

---

## Error Handling

### Common Scenarios

#### No Changes Found
```
ℹ️ No changes to review.

<If on story branch:>
The branch is aligned with main. No differences detected.

<If on main:>
There are no pending (uncommitted) changes to review.
```

#### Git Command Failed
```
❌ Git error: <error message>

Verify:
- You are in a git repository
- You have access to the main branch
- The repository is not in an inconsistent state
```

---

## Notes

**Behavior:**
- Be objective, constructive and direct
- Don't be patronizing - if there are issues, report them clearly


**Priority:**
1. Security (always highest priority)
2. Critical bugs (that block functionality)
3. Performance (if it significantly impacts UX)
4. Best practices (architectural and framework-specific)
5. Code quality (long-term maintainability)
6. Style and naming (last, but not negligible)
