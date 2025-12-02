# Coding Standards

Questo file definisce gli standard di sviluppo universali per l'implementazione di user stories. Gli standard sono project-agnostic e si applicano a qualsiasi linguaggio o framework.

---

## 1. Code Quality Standards

### Clean Code Principles

**Nomenclatura:**
- Nomi significativi e auto-esplicativi per variabili, funzioni, classi
- Evitare abbreviazioni criptiche (usa `userRepository` non `usrRepo`)
- Nomi di funzioni devono essere verbi (`getUserById`, `calculateTotal`)
- Nomi di classi devono essere sostantivi (`User`, `OrderService`)

**Funzioni:**
- Funzioni brevi: massimo 20-30 righe
- Una responsabilità per funzione (Single Responsibility)
- Massimo 3-4 parametri per funzione
- Evitare side effects nascosti

**Complessità:**
- Evitare annidamenti profondi (max 3 livelli di if/for)
- Preferire early returns per ridurre complessità
- Estrarre logica complessa in funzioni separate

**DRY (Don't Repeat Yourself):**
- No duplicazione di codice
- Estrarre logica ripetuta in funzioni/utility
- Riutilizzare componenti esistenti quando possibile

### SOLID Principles (riferimento)

- **S**ingle Responsibility: Una classe/funzione = una responsabilità
- **O**pen/Closed: Aperto all'estensione, chiuso alla modifica
- **L**iskov Substitution: Le sottoclassi devono sostituire le classi base
- **I**nterface Segregation: Interfacce specifiche, non monolitiche
- **D**ependency Inversion: Dipendere da astrazioni, non implementazioni concrete

---

## 2. Test Patterns

### Test Structure (AAA Pattern)

Ogni test deve seguire la struttura **Arrange-Act-Assert**:

```
// Arrange - Setup del contesto
const user = createTestUser({ role: 'admin' });
const repository = new UserRepository();

// Act - Esecuzione dell'azione da testare
const result = await repository.findById(user.id);

// Assert - Verifica del risultato
expect(result).toBeDefined();
expect(result.id).toBe(user.id);
expect(result.role).toBe('admin');
```

### Coverage Expectations

- **Unit tests**: Testare singole funzioni/metodi in isolamento
- **Integration tests**: Testare interazioni tra componenti
- **E2E tests**: Testare flussi completi utente

**Priorità:**
- Coprire tutti gli acceptance criteria GHERKIN della story
- Testare happy path (scenario principale)
- Testare error handling (validazione input, errori previsti)
- Testare edge cases (boundary conditions, valori limite)

**Target coverage:** Minimo 80% per nuovo codice (quando specificato dal progetto)

### Test Naming Conventions

**Pattern:** `should_ExpectedBehavior_When_Condition`

Esempi:
- `should_ReturnUser_When_ValidIdProvided`
- `should_ThrowError_When_UserNotFound`
- `should_RejectInvalidEmail_When_CreatingUser`

**Alternative (BDD style):**
- `it('returns user when valid ID is provided')`
- `it('throws error when user not found')`

---

## 3. Git Commit Format

### Conventional Commits

Seguire lo standard [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): TK-XXX - brief description

- Implementation details (1-3 bullet points)
- Files: list of main files changed
- Tests: ✅ passing / ❌ failing (if applicable)

🤖 Generated with Claude Code
Co-Authored-By: Claude <noreply@anthropic.com>
```

### Commit Types

- **feat**: Nuova funzionalità (la maggior parte dei task)
- **fix**: Bug fix
- **refactor**: Ristrutturazione codice senza cambiare funzionalità
- **test**: Aggiunta o modifica di test
- **docs**: Modifiche alla documentazione
- **chore**: Task di manutenzione (build, config, dependencies)
- **perf**: Miglioramenti performance
- **style**: Formattazione codice (whitespace, formatting, no logic change)

### Scope

Lo scope è la user story di riferimento (es: `US-005`, `US-012`)

### Esempi Completi

**Feature implementation:**
```
feat(US-005): TK-012 - Implement ISBN cataloging API endpoint

- Created POST /api/books/isbn endpoint with OpenLibrary integration
- Added validation for ISBN-10 and ISBN-13 formats
- Implemented error handling for API failures
- Files: src/books/books.controller.ts, src/books/isbn.service.ts, src/books/dto/isbn.dto.ts
- Tests: ✅ passing (12 tests, 4 scenarios)

🤖 Generated with Claude Code
Co-Authored-By: Claude <noreply@anthropic.com>
```

**Bug fix:**
```
fix(US-008): TK-025 - Fix date validation in book metadata form

- Corrected regex pattern for ISO date format
- Added edge case handling for leap years
- Files: src/books/validators/date.validator.ts
- Tests: ✅ passing (8 tests)

🤖 Generated with Claude Code
Co-Authored-By: Claude <noreply@anthropic.com>
```

**Refactoring:**
```
refactor(US-003): TK-018 - Extract IP validation logic to utility

- Moved IP whitelist validation from controller to utility
- Improved testability and reusability
- Files: src/auth/guards/ip.guard.ts, src/utils/ip-validator.util.ts
- Tests: ✅ passing (6 tests)

🤖 Generated with Claude Code
Co-Authored-By: Claude <noreply@anthropic.com>
```

---

## 4. Documentation Standards

### Quando Aggiungere Commenti

**Aggiungere commenti quando:**
- Logica non ovvia o algoritmi complessi
- Workaround per bug di librerie esterne
- Decisioni architetturali non evidenti dal codice
- Regex complesse o formule matematiche
- Vincoli di business non ovvi

**NON aggiungere commenti quando:**
- Il codice è auto-esplicativo con buoni nomi
- Si ripete semplicemente cosa fa il codice (ridondante)
- Si può riscrivere il codice per renderlo più chiaro

**Esempio buono:**
```typescript
// Workaround: OpenLibrary API sometimes returns 503 under load.
// Retry with exponential backoff (max 3 attempts)
const book = await retryWithBackoff(() => openLibraryApi.fetchByISBN(isbn), 3);
```

**Esempio cattivo (non fare):**
```typescript
// Increment counter by 1
counter++;
```

### Dev Notes Format

Dopo ogni task completato, aggiungere una entry nella sezione **Dev Notes** del file story:

```markdown
### TK-XXX Implementation (YYYY-MM-DD)

**What was done:**
- Breve descrizione (1-3 punti) dell'implementazione

**Files changed:**
- path/to/file1.ts (created/modified/deleted)
- path/to/file2.ts (modified)

**Tests:** ✅ All tests passing / ❌ Tests failing (with details)

**Notes:** (opzionale)
- Decisioni tecniche importanti
- Problemi risolti
- Miglioramenti futuri suggeriti
```

---

## 5. Error Handling

### Graceful Degradation

- Gestire sempre gli errori prevedibili
- Fornire fallback quando possibile
- Non lasciare l'applicazione in stato inconsistente

**Esempio:**
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

- Messaggi chiari e comprensibili per l'utente finale
- Evitare stack traces o dettagli tecnici nell'UI
- Includere suggerimenti su come risolvere (quando possibile)

**Buono:**
```
"Impossibile trovare il libro. Verifica che l'ISBN sia corretto (formato: 978-0-123456-78-9)."
```

**Cattivo:**
```
"Error: ECONNREFUSED 127.0.0.1:3000"
```

### Logging Practices

**Livelli di log:**
- **ERROR**: Errori che richiedono attenzione immediata
- **WARN**: Situazioni anomale ma gestite
- **INFO**: Eventi importanti (startup, config changes, major operations)
- **DEBUG**: Informazioni dettagliate per debugging

**Cosa loggare:**
- Errori con context (user ID, request ID, operation)
- Operazioni critiche (login, purchase, data deletion)
- Chiamate API esterne (request/response)
- Performance metrics per operazioni lente

**Cosa NON loggare:**
- Password o token di autenticazione
- Dati sensibili (numeri carta di credito, dati sanitari)
- PII (Personally Identifiable Information) non necessaria

---

## 6. Implementation Workflow

Quando implementi un task, segui questo workflow:

1. **Leggi Architecture Notes** della story per capire componenti e pattern
2. **Leggi Acceptance Criteria** per capire comportamento atteso
3. **Identifica file da creare/modificare** basandoti su architecture
4. **Implementa in piccoli incrementi** testabili
5. **Scrivi/aggiorna test** per coprire acceptance criteria
6. **Esegui test** e verifica che passino
7. **Commit** con messaggio conventional commits
8. **Aggiorna Dev Notes** nel file story

---

## 7. Project-Specific Context

Questo file fornisce standard universali. Per dettagli specifici del progetto (framework, librerie, architettura, API), fai riferimento al file **`docs/prd.md`** che viene iniettato nel contesto.

Il PRD contiene:
- Stack tecnologico (frontend/backend/database)
- Struttura delle directory
- Pattern architetturali specifici
- Librerie e dipendenze utilizzate
- Convenzioni specifiche del progetto

**Quando c'è conflitto tra questo file e il PRD, il PRD ha priorità** in quanto contiene le specifiche del progetto corrente.

---

## 8. Quality Checklist

Prima di marcare un task come completato, verifica:

- [ ] Il codice segue i clean code principles
- [ ] I nomi di variabili/funzioni sono significativi
- [ ] Non c'è duplicazione di codice (DRY)
- [ ] Gli acceptance criteria GHERKIN sono coperti da test
- [ ] I test passano (happy path + error cases + edge cases)
- [ ] Gli errori sono gestiti gracefully
- [ ] I messaggi di errore sono user-friendly
- [ ] Il commit segue conventional commits format
- [ ] Le Dev Notes sono aggiornate
- [ ] Non ci sono warning del linter/compiler

---

**Nota finale:** Questi standard sono linee guida, non regole rigide. Usa il buon senso e adatta quando necessario per il contesto specifico del task. L'obiettivo è codice pulito, manutenibile e testabile.
