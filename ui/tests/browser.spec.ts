import { expect, test, type Page } from '@playwright/test';

const extensionName = 'inari-ext-argocd';
async function mountRemote(page: Page): Promise<void> {
  const remoteEntry = process.env.INARI_REMOTE_ENTRY;
  await page.route(`**/api/extensions/${extensionName}/actions/*`, async (route) => {
    await route.fulfill({ json: { commandId: 'browser-command-1', outcome: 'applied' } });
  });

  await page.goto('/acme/extensions');
  const existing = page.locator('tr', { hasText: extensionName });
  await expect(existing).toHaveCount(1);
  await existing.getByRole('button', { name: 'Remove', exact: true }).click();
  await expect(existing).toHaveCount(0);
  await page.locator('#ext-name').fill(extensionName);
  await page.locator('#ext-entry').fill(remoteEntry!);
  await page.getByRole('button', { name: 'Add remote', exact: true }).click();
  await expect(page.locator('tr', { hasText: extensionName }).getByText('ready', { exact: true })).toBeVisible();
}

async function openInstanceActions(page: Page): Promise<void> {
  await page.goto('/acme/deploys/ri-orders-db');
  await page.getByRole('button', { name: 'Actions', exact: true }).click();
}

let pageErrors: string[];

test.beforeEach(async ({ page }) => {
  pageErrors = [];
  page.on('pageerror', (error) => pageErrors.push(error.message));
  await mountRemote(page);
});

test.afterEach(async () => {
  expect(pageErrors).toEqual([]);
});

test('slot cluster-tab:argocd-health', async ({ page }) => {
  await page.goto('/acme/clusters/cl-kind-dev');
  await page.getByRole('tab', { name: 'ArgoCD', exact: true }).click();
  await expect(page.getByText(/ArgoCD health —/).first()).toBeVisible();
  await expect(page.getByText('ri-orders-db')).toBeVisible();
});

test('slot catalog-card:argocd-badge', async ({ page }) => {
  await page.goto('/acme/catalog');
  await expect(page.getByText('GitOps · ArgoCD', { exact: true }).first()).toBeVisible();
});

test('slot instance-action:argocd-sync', async ({ page }) => {
  await openInstanceActions(page);
  const response = page.waitForResponse(`**/api/extensions/${extensionName}/actions/argocd.sync`);
  await page.getByRole('button', { name: 'Sync (ArgoCD)', exact: true }).click();
  expect((await response).ok()).toBe(true);
});

test('slot instance-action:argocd-refresh', async ({ page }) => {
  await openInstanceActions(page);
  const response = page.waitForResponse(`**/api/extensions/${extensionName}/actions/argocd.refresh`);
  await page.getByRole('button', { name: 'Refresh (ArgoCD)', exact: true }).click();
  expect((await response).ok()).toBe(true);
});

test('slot instance-action:argocd-rollback', async ({ page }) => {
  await openInstanceActions(page);
  await page.getByRole('button', { name: 'Rollback (ArgoCD)', exact: true }).click();
  await expect(page.getByRole('alert').filter({ hasText: 'has no recorded ArgoCD revision' })).toBeVisible();
});
