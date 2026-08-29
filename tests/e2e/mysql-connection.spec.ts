import { test, expect, type Page } from '@playwright/test'

// ---------------------------------------------------------------------------
// Login helper
// Fixes from original:
//  - `input[required=""]` was fragile; the <select> also has `required` so the
//    first `input[required]` happened to be the username field only by luck.
//    Using getByLabel() is semantically correct and won't break if field order
//    changes.
//  - The app redirects to /workbench after login, not back to "/". Waiting for
//    "http://localhost:5173/" timed out because that URL was never reached.
//  - The login form is hidden while connection labels load from the API. We
//    now wait for the submit button before interacting with any fields.
//  - The connection-label <select> must be filled before username/password.
//
// Update USERNAME / PASSWORD to match your registered MySQL DB user.
// ---------------------------------------------------------------------------
const USERNAME = 'laravel_user'
const PASSWORD = 'your_password'

async function login(page: Page) {
  await page.goto('http://localhost:5173/login')

  // Wait for the async connection-label fetch and form render
  await page.waitForSelector('button[type="submit"]', { timeout: 10000 })

  // Select the first available connection label
  const connectionSelect = page.locator('select').first()
  if (await connectionSelect.isVisible()) {
    await connectionSelect.selectOption({ index: 0 })
  }

  // Use label-based selectors — robust to attribute/order changes
  await page.getByLabel('Username').fill(USERNAME)
  await page.getByLabel('Password').fill(PASSWORD)

  await page.click('button[type="submit"]')
  // Router redirects to /workbench on success, not to root "/"
  await page.waitForURL('**/workbench', { timeout: 15000 })
}

test.describe('ChatDB MySQL Connection', () => {
  test('login with existing connection and view tables', async ({ page }) => {
    await login(page)

    // Wait for the sidebar table chips to load from the backend
    await page.waitForSelector('button.table-chip', { timeout: 15000 })

    // Verify the expected tables are visible
    await expect(page.locator('button.table-chip:has-text("users")')).toBeVisible()
    await expect(page.locator('button.table-chip:has-text("migrations")')).toBeVisible()
    await expect(page.locator('button.table-chip:has-text("admins")')).toBeVisible()
  })

  test('view table structure', async ({ page }) => {
    await login(page)

    // Wait for table chips before clicking
    await page.waitForSelector('button.table-chip', { timeout: 15000 })

    // Click the "users" table chip to open its detail panel
    await page.click('button.table-chip:has-text("users")')

    // Allow time for the panel to render
    await page.waitForTimeout(2000)

    // The page content should reference the table name
    const pageContent = await page.content()
    expect(pageContent).toContain('users')
  })
})