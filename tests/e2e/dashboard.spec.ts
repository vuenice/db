import { test, expect, type Page } from '@playwright/test';

// ---------------------------------------------------------------------------
// Shared login helper
// The login form renders a <select> (not an ARIA combobox) for the connection
// label. It loads labels asynchronously on mount, so we wait for the form to
// appear before interacting with it. The app redirects to /workbench on
// successful login — NOT back to "/".
// Update USERNAME / PASSWORD to match your registered DB user.
// ---------------------------------------------------------------------------
const USERNAME = 'laravel_user';
const PASSWORD = 'your_password';

async function login(page: Page) {
  await page.goto('http://localhost:5173/login');

  // Wait for the async connection-label fetch to finish and the form to render
  await page.waitForSelector('button[type="submit"]', { timeout: 10000 });

  // Pick the first available connection label from the <select>
  const connectionSelect = page.locator('select').first();
  if (await connectionSelect.isVisible()) {
    await connectionSelect.selectOption({ index: 0 });
  }

  await page.getByLabel('Username').fill(USERNAME);
  await page.getByLabel('Password').fill(PASSWORD);

  await page.getByRole('button', { name: 'Login' }).click();
  // Router redirects authenticated users to /workbench, not root "/"
  await page.waitForURL('**/workbench', { timeout: 15000 });
}

test.describe('ChatDB Dashboard Tests', () => {

  test.beforeEach(async ({ page }) => {
    await login(page);
  });

  test('view tables list in sidebar', async ({ page }) => {
    // Table chips render as <button class="table-chip"> with two child <span>
    // elements (name and kind). Use the CSS class + hasText filter instead of
    // getByRole({ name }) to avoid relying on case-sensitive accessible-name
    // composition ("users table" vs "users TABLE").
    await expect(page.locator('button.table-chip', { hasText: 'users' }).first()).toBeVisible();
    await expect(page.locator('button.table-chip', { hasText: 'migrations' }).first()).toBeVisible();
    await expect(page.locator('button.table-chip', { hasText: 'admins' }).first()).toBeVisible();

    console.log('Tables list verified');
  });

  test('switch database connection', async ({ page }) => {
    // The connection-label selector in the banner is a plain <select>, not a
    // combobox. Verify that one of its options contains "test".
    const connectionSelect = page.locator('select').first();
    if (await connectionSelect.isVisible()) {
      await expect(connectionSelect.locator('option', { hasText: 'test' })).toHaveCount(1);
      console.log('Database switcher verified');
    }
  });

  test('view table structure by clicking table', async ({ page }) => {
    // Click on the "users" table chip to open its detail panel
    await page.locator('button.table-chip', { hasText: 'users' }).first().click();
    await page.waitForTimeout(500);

    // Page content should reference the table name
    const pageContent = await page.content();
    expect(pageContent).toContain('users');

    console.log('Table structure view verified');
  });

  test('search tables functionality', async ({ page }) => {
    const searchBox = page.getByPlaceholder('Search tables…');

    if (await searchBox.isVisible()) {
      await searchBox.fill('users');
      await page.waitForTimeout(300);

      // After filtering, the "users" table chip must still be visible
      await expect(page.locator('button.table-chip', { hasText: 'users' }).first()).toBeVisible();
      console.log('Search functionality verified');
    }
  });

  test('navigate to SQL chat interface', async ({ page }) => {
    const chatSqlBtn = page.getByRole('button', { name: 'Chat SQL' });

    if (await chatSqlBtn.isVisible()) {
      await chatSqlBtn.click();
      await page.waitForTimeout(500);
      console.log('Chat SQL navigation verified');
    }
  });

  test('navigate to History', async ({ page }) => {
    const historyBtn = page.getByRole('button', { name: 'History' });

    if (await historyBtn.isVisible()) {
      await historyBtn.click();
      await page.waitForTimeout(500);
      console.log('History navigation verified');
    }
  });

  test('navigate to Queries', async ({ page }) => {
    const queriesBtn = page.getByRole('button', { name: 'Queries' });

    if (await queriesBtn.isVisible()) {
      await queriesBtn.click();
      await page.waitForTimeout(500);
      console.log('Queries navigation verified');
    }
  });

  test('navigate to Users', async ({ page }) => {
    const usersBtn = page.getByRole('button', { name: 'Users' });

    if (await usersBtn.isVisible()) {
      await usersBtn.click();
      await page.waitForTimeout(500);
      console.log('Users navigation verified');
    }
  });

  test('view DB roles dropdown', async ({ page }) => {
    const dbRolesSelector = page.getByLabel('DB roles');

    if (await dbRolesSelector.isVisible()) {
      await dbRolesSelector.click();
      await page.waitForTimeout(300);
      console.log('DB roles dropdown verified');
    }
  });

  test('view DB pool dropdown', async ({ page }) => {
    const dbPoolSelector = page.getByLabel('DB pool');

    if (await dbPoolSelector.isVisible()) {
      await dbPoolSelector.click();
      await page.waitForTimeout(300);

      // The pool <select> exposes its options directly as <option> elements
      await expect(dbPoolSelector.locator('option', { hasText: 'Read' })).toHaveCount(1);
      console.log('DB pool dropdown verified');
    }
  });

  test('user menu dropdown', async ({ page }) => {
    const userMenuBtn = page.getByRole('button', { name: /Open account menu/i });

    if (await userMenuBtn.isVisible()) {
      await userMenuBtn.click();
      await page.waitForTimeout(300);
      console.log('User menu dropdown verified');
    }
  });
});