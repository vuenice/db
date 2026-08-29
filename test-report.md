# ChatDB / VueNiceDB — Playwright E2E Test Report
**Date:** 2026-06-13
**Method:** Static analysis of Vue source + cross-referencing against all test selectors and assertions. Playwright browser download is blocked by sandbox network policy; Go compiler unavailable (Windows-only `chatdb.exe` binary committed). Analysis traces every selector and assertion against the actual Vue/TS source code. Actual previous test failure artifacts (`error-context.md`) also confirmed.

---

## Environment Summary

| Item | Status |
|---|---|
| Frontend (Vue 3 + Vite) | Working tree has corrected `vite.config.ts`; **committed HEAD has wrong `@` alias** |
| Backend (Go) | `chatdb.exe` — Windows PE32+ only; no Linux binary; no Go compiler available |
| Database | Not confirmed running (MySQL credentials are placeholders) |
| Playwright config (`tests/`) | `baseURL: http://localhost:5173`, `webServer: null` — no auto-start |
| Playwright browsers | Not installed; download blocked |
| `tests/package.json` script | `"test": "echo error && exit 1"` — broken, cannot `npm test` |

---

## Test File Inventory

| File | Location | Tests |
|---|---|---|
| `app.spec.ts` | `tests/e2e/` | 3 |
| `dashboard.spec.ts` | `tests/e2e/` | 11 |
| `mysql-connection.spec.ts` | `tests/e2e/` | 2 |
| `import-export.spec.ts` | `frontend/tests/e2e/` (unstaged) | 4 |
| `operations.spec.ts` | `frontend/tests/e2e/` (unstaged) | 1 |
| **Total** | | **21** |

---

## Test Results

### `tests/e2e/app.spec.ts` — Basic Page Load Tests

| # | Test | Result | Reason |
|---|---|---|---|
| 1 | `register page loads` | ✅ PASS | `RegisterView.vue` renders without any API call on mount. `<h1>ChatDB</h1>` matches, `input[placeholder="e.g. production"]` matches `placeholder*="production"` |
| 2 | `login page loads` | ⚠️ FLAKY | `LoginView.vue` starts with `pageLoading = true` (form hidden). After `onMounted` resolves the connection-label fetch, if backend has no registered connections, router redirects to `/register` before the `h1`/button check runs → FAIL. If backend is down (network error), error is caught and form stays → PASS. |
| 3 | `register then redirect to login` | ⚠️ FLAKY | Clicking "Already have an account? Sign in" navigates to `/login`. Login page `onMounted` runs the same risky redirect-to-register path as test 2. URL/h1 assertions may then fail. |

---

### `tests/e2e/dashboard.spec.ts` — Dashboard Tests (11 tests)

**All 11 FAIL.** Every test runs `login()` in `beforeEach`. Login requires a running backend at `http://127.0.0.1:6366` (unavailable — Windows-only binary) and real credentials (hardcoded placeholders: `laravel_user` / `your_password`). `waitForURL('**/workbench', { timeout: 15000 })` times out in every test.

| # | Test | Result | Failure point |
|---|---|---|---|
| 1 | `view tables list in sidebar` | ❌ FAIL | `login()` → timeout waiting for `/workbench` |
| 2 | `switch database connection` | ❌ FAIL | Same + wrong selector (see Bug #7) |
| 3 | `view table structure by clicking table` | ❌ FAIL | `login()` timeout |
| 4 | `search tables functionality` | ❌ FAIL | `login()` timeout |
| 5 | `navigate to SQL chat interface` | ❌ FAIL | `login()` timeout |
| 6 | `navigate to History` | ❌ FAIL | `login()` timeout |
| 7 | `navigate to Queries` | ❌ FAIL | `login()` timeout |
| 8 | `navigate to Users` | ❌ FAIL | `login()` timeout |
| 9 | `view DB roles dropdown` | ❌ FAIL | `login()` timeout |
| 10 | `view DB pool dropdown` | ❌ FAIL | `login()` timeout |
| 11 | `user menu dropdown` | ❌ FAIL | `login()` timeout |

**Additional bugs that would cause failures even with a working login:**

- **Test 9:** `page.getByLabel('DB roles')` targets a block inside `v-if="selectedConnId"`. If no connection is selected, the block is hidden and `isVisible()` guard silently skips — the test validates nothing.
- **Test 10:** `page.getByLabel('DB pool')` is inside `v-if="auth.isEngineer"`. A non-engineer login hides it entirely.

---

### `tests/e2e/mysql-connection.spec.ts` — MySQL Connection Tests (2 tests)

**Both FAIL** — identical root cause as dashboard tests.

| # | Test | Result | Reason |
|---|---|---|---|
| 1 | `login with existing connection and view tables` | ❌ FAIL | `login()` timeout; no backend; placeholder creds |
| 2 | `view table structure` | ❌ FAIL | Same |

---

### `frontend/tests/e2e/import-export.spec.ts` — Import/Export Tests (4 tests, unstaged)

**All 4 FAIL.** Confirmed by actual error artifacts in `frontend/test-results/`:

> `Error: expect(locator).toBeVisible() failed — Locator: locator('select.io-input').first() — navigated to "http://localhost:5173/login"`

The `beforeEach` login helper failed (bad credentials + no backend), leaving the page on `/login`. All subsequent assertions target workbench-specific elements that are never rendered.

| # | Test | Result | Status |
|---|---|---|---|
| 1 | `should display import page with all required UI elements` | ❌ FAIL | Error artifact confirmed (`error-context.md` present) |
| 2 | `should display export page with all required UI elements` | ❌ FAIL | Error artifact confirmed |
| 3 | `should switch between psql and pgdump import formats` | ❌ FAIL | Same root cause |
| 4 | `should switch between plain and archive export formats` | ❌ FAIL | Same root cause |

Note: The selectors themselves (`select.io-input`, `button.format-card`, `button.import-btn`, `button.primary.footer-btn`) do match the actual Vue source — the only blocker is authentication.

---

### `frontend/tests/e2e/operations.spec.ts` — Operations Panel (1 test, unstaged)

| # | Test | Result | Reason |
|---|---|---|---|
| 1 | `should display all operation cards in the UI` | ❌ FAIL | `expect(page).toHaveTitle(/frontend\|ChatDB/)` — actual `<title>` is `VueNiceDB` (`frontend/index.html` line 7). Neither word in the regex matches. |

---

## Overall Results

| File | Tests | Pass | Fail | Flaky |
|---|---|---|---|---|
| `app.spec.ts` | 3 | 1 | 0 | 2 |
| `dashboard.spec.ts` | 11 | 0 | 11 | 0 |
| `mysql-connection.spec.ts` | 2 | 0 | 2 | 0 |
| `import-export.spec.ts` | 4 | 0 | 4 | 0 |
| `operations.spec.ts` | 1 | 0 | 1 | 0 |
| **Total** | **21** | **1** | **18** | **2** |

---

## Bugs Found

### 🔴 Critical — Breaks Committed Build

**Bug 1: `vite.config.ts` has wrong `@` alias in committed HEAD**

The committed HEAD version has:
```ts
resolve: {
  alias: {
    '@': 'C:/laragon/www/dail-it',  // wrong: another project on developer's machine
  },
},
```

The working tree already has the correct fix (`path.resolve(__dirname, './src')`), but it is **not committed**. A fresh `git clone` gets the broken alias. This breaks TypeScript resolution for any file that uses `@/` imports.

**Fix:** Commit the working-tree `vite.config.ts`.

---

**Bug 2: `WorkbenchView.vue` imports non-existent components**

```ts
import VueNiceTable from '@/shared/common/CommonTable.vue'  // does not exist
import VueNiceModal from '@/shared/common/CommonModal.vue'  // does not exist
```

`src/shared/` does not exist. Currently dead code (the router imports `workbench/index.vue` instead), so the dev build does not fail. But the moment this file is imported anywhere, the build hard-fails.

**Fix:** Delete `WorkbenchView.vue`, or create `src/shared/common/CommonTable.vue` and `CommonModal.vue`.

---

### 🟠 High — Tests Cannot Run At All

**Bug 3: No cross-platform backend binary**

Only `chatdb.exe` (Windows PE32+) is committed. No Linux/macOS binary, no pre-built Docker image. No one outside Windows can start the backend.

**Fix:** Add `make build-linux` and commit a Linux binary, or add a proper Docker Compose service for the Go app.

---

**Bug 4: `tests/package.json` has a broken test script**

```json
"scripts": { "test": "echo \"Error: no test specified\" && exit 1" }
```

Running `npm test` in `tests/` always fails. This is the standard CI entry point.

**Fix:** `"test": "playwright test"`

---

**Bug 5: `tests/playwright.config.ts` has `webServer: null`**

No server is auto-started before tests run. Developers must manually start Vite and the backend first.

**Fix:** Add a `webServer` block that starts Vite automatically.

---

### 🟡 Medium — Wrong Values / Logic in Tests

**Bug 6: Test credentials are placeholder values**

All auth-dependent test files use `USERNAME = 'laravel_user'` / `PASSWORD = 'your_password'`.

**Fix:** Use env vars: `process.env.TEST_USERNAME ?? 'laravel_user'` / `process.env.TEST_PASSWORD ?? ''`, and document seeding the test DB user.

---

**Bug 7: `dashboard.spec.ts` test 2 checks the wrong `<select>` for connection switching**

In the workbench, `page.locator('select').first()` is the **physical database** switcher, not the connection-label selector (which only exists on the login page). The test asserts an option named `'test'` in the wrong element.

---

**Bug 8: `operations.spec.ts` page title regex is wrong**

`/frontend|ChatDB/` does not match the actual title `VueNiceDB`.

**Fix:** `expect(page).toHaveTitle(/VueNiceDB/)`

---

**Bug 9: Dashboard tests assert MySQL tables; project uses PostgreSQL**

Tests look for `button.table-chip:has-text("users")`, `"migrations"`, `"admins"` — a Laravel MySQL schema. The `docker-compose.yml` provisions PostgreSQL with a `widgets` sample table.

---

### 🟢 Low — Code Quality

**Bug 10: `WorkbenchView.vue` is unreferenced dead code** — not imported by the router or `App.vue`.

**Bug 11: Two conflicting `playwright.config.ts` files** — `frontend/playwright.config.ts` (port 6366, starts Go backend) and `tests/playwright.config.ts` (port 5173, `webServer: null`) have different, incompatible assumptions. No single source of truth.

---

## Recommended Fixes (Priority Order)

1. **Commit the corrected `vite.config.ts`** — the working-tree fix exists, it just needs `git add && git commit`.
2. **Fix `tests/package.json`** test script to `"playwright test"`.
3. **Fix `operations.spec.ts`** title regex to `/VueNiceDB/`.
4. **Replace placeholder credentials** with env-var-driven values across all test files.
5. **Add a `webServer` block** in `tests/playwright.config.ts` to auto-start Vite.
6. **Provide a Linux backend binary** or a Docker Compose service for the Go app.
7. **Fix `dashboard.spec.ts` test 2** — use the correct selector for connection-label switching.
8. **Align table assertions with actual PostgreSQL schema** — use `widgets`, not MySQL tables.
9. **Delete `WorkbenchView.vue`** or create the missing `src/shared/common/` components.

---

## Suggested New Feature (if all tests were passing)

**Query History — Date Range & Table Filter**

The History panel (`nav === 'history'`) already loads past queries and has a basic text search input (`placeholder="Search logs, queries…"`). Adding filter controls for **date range** (today / last 7 days / custom) and **table name** (extracted from SQL via a simple regex on `FROM <table>`) would make the history genuinely useful for teams running high query volumes. This is a pure frontend change targeting the existing `queryHistory` ref and `loadQueryHistory()` function — no backend API changes required.
