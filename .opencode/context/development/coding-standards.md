# Coding Standards

This file defines the development standards. The standards are project-agnostic and apply to any language or framework.

---

## Behavioral Guidelines

- When you find discrepancies between testing and implementation, stop and always ask me which approach I prefer (modify testing vs. implementation) before proceeding.
- Implement only what is requested, without adding extra features.
- Always avoid testing and implementation at the same time. If you are modifying project files, do not modify test files, and vice versa.
- Before modifying any files (test or implementation), always present the available options and wait for my approval.


## 1. Code Quality Standards

### Clean Code Principles

**Naming:**
- Use meaningful and self-explanatory names for variables, functions, and classes
- Avoid cryptic abbreviations (use `userRepository`, not `usrRepo`)
- Function names must be verbs (`getUserById`, `calculateTotal`)
- Class names must be nouns (`User`, `OrderService`)

**Functions:**
- Keep functions short: maximum 20-30 lines
- One responsibility per function (Single Responsibility)
- Maximum 3-4 parameters per function
- Avoid hidden side effects

**Complexity:**
- Avoid deep nesting (max 3 levels of if/for)
- Prefer early returns to reduce complexity
- Extract complex logic into separate functions

**DRY (Don't Repeat Yourself):**
- Do not duplicate code
- Extract repeated logic into helper functions/utilities
- Reuse existing components whenever possible

### SOLID Principles (reference)

- **S**ingle Responsibility: One class/function = one responsibility
- **O**pen/Closed: Open for extension, closed for modification
- **L**iskov Substitution: Subclasses must substitute their base classes
- **I**nterface Segregation: Prefer specific interfaces, not monolithic ones
- **D**ependency Inversion: Depend on abstractions, not concrete implementations

---

## 2. Test Patterns

### Test Structure (AAA Pattern)

Every test must follow the **Arrange-Act-Assert** structure:

```
// Arrange - Set up the context
const user = createTestUser({ role: 'admin' });
const repository = new UserRepository();

// Act - Execute the action under test
const result = await repository.findById(user.id);

// Assert - Verify the result
expect(result).toBeDefined();
expect(result.id).toBe(user.id);
expect(result.role).toBe('admin');
```

### Coverage Expectations

**Priority:**
- Cover all story acceptance criteria written in GHERKIN
- Test the happy path (primary scenario)
- Test error handling (input validation, expected errors)
- Test edge cases (boundary conditions, limit values)

**Target coverage:** Minimum 80% for code describing business logic.

### Test Naming Conventions

**Pattern:** `should_ExpectedBehavior_When_Condition`

Examples:
- `should_ReturnUser_When_ValidIdProvided`
- `should_ThrowError_When_UserNotFound`
- `should_RejectInvalidEmail_When_CreatingUser`

**Alternative (BDD style):**
- `it('returns user when valid ID is provided')`
- `it('throws error when user not found')`

---

## 3. Git Commit Format

### Conventional Commits

Follow the [Conventional Commits](https://www.conventionalcommits.org/) standard:

```
<type>(<scope>): TK-XXX - brief description

- Implementation details (1-3 bullet points)
```

### Commit Types

- **feat**: New feature (most tasks)
- **fix**: Bug fix
- **refactor**: Code restructuring without functionality changes
- **test**: Adding or modifying tests
- **docs**: Documentation changes
- **chore**: Maintenance tasks (build, config, dependencies)
- **perf**: Performance improvements
- **style**: Code formatting (whitespace, formatting, no logic change)

## 4. Documentation Standards

### When to Add Comments

**Add comments when:**
- Logic is not obvious or algorithms are complex
- Workarounds are needed for external library bugs
- Architectural decisions are not evident from the code
- Regexes or mathematical formulas are complex
- Business constraints are not obvious

**DO NOT add comments when:**
- Code is self-explanatory thanks to good naming
- The comment merely restates what the code does (redundant)
- The code can be rewritten to make it clearer

**Good example:**
```typescript
// Workaround: OpenLibrary API sometimes returns 503 under load.
// Retry with exponential backoff (max 3 attempts)
const book = await retryWithBackoff(() => openLibraryApi.fetchByISBN(isbn), 3);
```

**Bad example (do not do this):**
```typescript
// Increment counter by 1
counter++;
```

---

## 5. Error Handling

### Graceful Degradation

- Always handle predictable errors
- Provide fallbacks whenever possible
- Do not leave the application in an inconsistent state

**Example:**
```typescript
try {
  const bookData = await externalApi.fetchBookByISBN(isbn);
  return bookData;
} catch (error) {
  logger.error('Failed to fetch book from external API', { isbn, error });

  // Fallback: try local database
  const localBook = await db.books.findByISBN(isbn);
  if (localBook) {
    return localBook;
  }

  // If no fallback available, throw meaningful error
  throw new NotFoundException(`Book with ISBN ${isbn} not found`);
}
```

### User-Friendly Error Messages

- Provide clear, understandable messages for end users
- Avoid stack traces or technical details in the UI
- Include hints on how to resolve the issue whenever possible

**Good:**
```
"Unable to find the book. Verify the ISBN is correct (format: 978-0-123456-78-9)."
```

**Bad:**
```
"Error: ECONNREFUSED 127.0.0.1:3000"
```

### Logging Practices

**Log levels:**
- **ERROR**: Errors requiring immediate attention
- **WARN**: Anomalies that are handled
- **INFO**: Important events (startup, config changes, major operations)
- **DEBUG**: Detailed information for debugging

**What to log:**
- Errors with context (user ID, request ID, operation)
- Critical operations (login, purchase, data deletion)
- External API calls (request/response)
- Performance metrics for slow operations

**What NOT to log:**
- Passwords or authentication tokens
- Sensitive data (credit card numbers, medical data)
- Unnecessary PII (Personally Identifiable Information)

---

**Final note:** These standards are guidelines, not rigid rules. Use good judgment and adapt them to the specific task context. The goal is clean, maintainable, and testable code.
