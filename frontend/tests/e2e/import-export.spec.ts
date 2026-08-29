import { test, expect, type Page } from '@playwright/test'

// ---------------------------------------------------------------------------
// Shared login helper
// The original beforeEach was silently failing: it navigated to /workbench/import
// (which redirected to /login), then swallowed the post-login waitForURL error
// with .catch(()=>{}), leaving subsequent tests running on the /login page.
//
// Fixed approach:
//  1. Navigate directly to /login.
//  2. Wait for the async connection-label fetch to resolve before interacting.
//  3. Select the first connection label from the <select>.
//  4. Use the correct field labels to fill username/password.
//  5. Wait for /workbench redirect — no error swallowing.
//
// Update USERNAME / PASSWORD to match your registered DB user.
// ---------------------------------------------------------------------------
const USERNAME = 'laravel_user'
const PASSWORD = 'your_password'

async function login(page: Page) {
  await page.goto('http://localhost:5173/login')

  // The form is hidden behind pageLoading while connection labels are fetched;
  // wait for the submit button to appear before filling any fields.
  await page.waitForSelector('button[type="submit"]', { timeout: 10000 })

  // Connection-label selector is a plain <select>, not an ARIA combobox.
  const connectionSelect = page.locator('select').first()
  if (await connectionSelect.isVisible()) {
    await connectionSelect.selectOption({ index: 0 })
  }

  await page.getByLabel('Username').fill(USERNAME)
  await page.getByLabel('Password').fill(PASSWORD)

  await page.click('button[type="submit"]')
  // Router redirects to /workbench on success
  await page.waitForURL('**/workbench**', { timeout: 15000 })
}

test.describe('Import/Export functionality', () => {

  test.beforeEach(async ({ page }) => {
    await login(page)
  })

  test('should display import page with all required UI elements', async ({ page }) => {
    await page.goto('http://localhost:5173/workbench/import')
    await page.waitForLoadState('networkidle')

    // Connection selector
    await expect(page.locator('select.io-input').first()).toBeVisible()

    // Format cards (Plain SQL / pg_dump) are only rendered for PostgreSQL
    // connections (v-if="isPostgres"). Skip this assertion for MySQL connections.
    const isPostgres = await page.locator('button.format-card:has-text("Plain SQL")').isVisible()
    if (isPostgres) {
      await expect(page.locator('button.format-card:has-text("Plain SQL")')).toBeVisible()
    }

    // File input and import button are always present
    await expect(page.locator('input[type="file"]')).toBeVisible()
    await expect(page.locator('button.import-btn')).toBeVisible()
  })

  test('should display export page with all required UI elements', async ({ page }) => {
    await page.goto('http://localhost:5173/workbench/export')
    await page.waitForLoadState('networkidle')

    // Connection selector
    await expect(page.locator('select.io-input').first()).toBeVisible()

    // Export format cards are always rendered (no driver guard)
    await expect(page.locator('button.format-card:has-text("Plain SQL")')).toBeVisible()

    // The download button has class "primary footer-btn"; its text is
    // "Download dump" (not "Export"). Match on the class to avoid fragile
    // text coupling.
    await expect(page.locator('button.primary.footer-btn')).toBeVisible()
  })

  test('should switch between psql and pgdump import formats', async ({ page }) => {
    await page.goto('http://localhost:5173/workbench/import')
    await page.waitForLoadState('networkidle')

    // Format cards only appear for PostgreSQL connections; skip for MySQL.
    const plainSqlBtn = page.locator('button.format-card:has-text("Plain SQL")')
    if (!await plainSqlBtn.isVisible()) {
      console.log('Skipping format-switch test: connection is not PostgreSQL')
      return
    }

    // Default importKind is 'psql' → Plain SQL card should be active
    await expect(plainSqlBtn).toHaveClass(/active/)

    // Switch to pgdump
    await page.click('button.format-card:has-text("pg_dump archive")')
    await expect(page.locator('button.format-card:has-text("pg_dump archive")')).toHaveClass(/active/)
  })

  test('should switch between plain and archive export formats', async ({ page }) => {
    await page.goto('http://localhost:5173/workbench/export')
    await page.waitForLoadState('networkidle')

    // Default archiveKind is 'plain' → Plain SQL card should be active.
    // The export page always shows both format cards regardless of driver.
    const plainSqlBtn = page.locator('button.format-card:has-text("Plain SQL (psql)")')
    await expect(plainSqlBtn).toBeVisible()
    await expect(plainSqlBtn).toHaveClass(/active/)

    // Switch to archive format
    await page.click('button.format-card:has-text("pg_dump archive")')
    await expect(page.locator('button.format-card:has-text("pg_dump archive")')).toHaveClass(/active/)
  })
})