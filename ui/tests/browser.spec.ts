import { expect, test, type Page } from '@playwright/test';

const extensionName = 'inari-ext-argocd';
interface ObservedActionRequest {
  url: string;
  method: string;
  payload: unknown;
}

async function installActionDouble(page: Page): Promise<void> {
  await page.addInitScript((name) => {
    const browserWindow = window as typeof window & {
      __argocdActionRequests: ObservedActionRequest[];
    };
    browserWindow.__argocdActionRequests = [];
    const originalFetch = window.fetch.bind(window);
    window.fetch = async (input, init) => {
      const rawUrl = typeof input === 'string' ? input : input instanceof URL ? input.href : input.url;
      const url = new URL(rawUrl, window.location.href);
      const paths = [
        `/api/extensions/${name}/actions/argocd.sync`,
        `/api/extensions/${name}/actions/argocd.refresh`,
      ];
      if (url.origin === window.location.origin && paths.includes(url.pathname)) {
        browserWindow.__argocdActionRequests.push({
          url: url.href,
          method: (init?.method ?? (input instanceof Request ? input.method : 'GET')).toUpperCase(),
          payload: typeof init?.body === 'string' ? JSON.parse(init.body) : null,
        });
        return new Response(JSON.stringify({ commandId: 'browser-command-1', outcome: 'applied' }), {
          status: 200,
          headers: { 'content-type': 'application/json' },
        });
      }
      return originalFetch(input, init);
    };
  }, extensionName);
}

async function mountRemote(page: Page): Promise<void> {
  const remoteEntry = process.env.INARI_REMOTE_ENTRY;
  await page.goto('/acme/extensions');
  const existing = page
    .locator('tr', { hasText: extensionName })
    .filter({ has: page.getByRole('button', { name: 'Remove', exact: true }) });
  await expect(existing).toHaveCount(1);
  await existing.getByRole('button', { name: 'Remove', exact: true }).click();
  await expect(existing).toHaveCount(0);
  await page.locator('#ext-name').fill(extensionName);
  await page.locator('#ext-entry').fill(remoteEntry!);
  await page.getByRole('button', { name: 'Add remote', exact: true }).click();
  await expect(page.locator('tr', { hasText: extensionName }).getByText('ready', { exact: true })).toBeVisible();
}

async function navigateWithinHost(page: Page, path: string): Promise<void> {
  // Full reloads reset the host's in-memory MSW registry and discard the test remote.
  await page.evaluate((nextPath) => {
    window.history.pushState(null, '', nextPath);
    window.dispatchEvent(new PopStateEvent('popstate'));
  }, path);
  await expect(page).toHaveURL(new RegExp(`${path}$`));
}

async function openInstanceActions(page: Page): Promise<void> {
  await navigateWithinHost(page, '/acme/deploys/ri-orders-db');
  await page.getByRole('button', { name: 'Actions', exact: true }).click();
}

async function assertActionSuccess(
  page: Page,
  action: 'argocd.sync' | 'argocd.refresh',
  label: string,
  payload: unknown,
): Promise<void> {
  const path = `/api/extensions/${extensionName}/actions/${action}`;
  await expect.poll(() => page.evaluate((requestPath) => {
    const requests = (window as typeof window & { __argocdActionRequests: ObservedActionRequest[] })
      .__argocdActionRequests;
    return requests.filter((request) => new URL(request.url).pathname === requestPath).length;
  }, path)).toBe(1);
  const request = await page.evaluate((requestPath) => {
    const requests = (window as typeof window & { __argocdActionRequests: ObservedActionRequest[] })
      .__argocdActionRequests;
    return requests.find((candidate) => new URL(candidate.url).pathname === requestPath);
  }, path);
  const origin = await page.evaluate(() => window.location.origin);
  expect(request).toEqual({ url: `${origin}${path}`, method: 'POST', payload });
  await expect(page.getByRole('button', { name: label, exact: true })).toBeEnabled();
  // Hosts render informational alerts on the instance page (for example the Phase-3
  // workflow-render stop notice and unsigned-observability-metadata notice); only the
  // host's danger-styled alerts indicate an action failure.
  await expect(page.locator('[role="alert"].u-danger')).toHaveCount(0);
}

let pageErrors: string[];

test.beforeEach(async ({ page }) => {
  pageErrors = [];
  page.on('pageerror', (error) => pageErrors.push(error.message));
  await installActionDouble(page);
  await mountRemote(page);
});

test.afterEach(async () => {
  expect(pageErrors).toEqual([]);
});

test('slot cluster-tab:argocd-health', async ({ page }) => {
  await navigateWithinHost(page, '/acme/clusters/cl-kind-dev');
  await page.getByRole('tab', { name: 'ArgoCD', exact: true }).click();
  await expect(page.getByText(/ArgoCD health —/).first()).toBeVisible();
  await expect(page.getByText('ri-orders-db')).toBeVisible();
});

test('slot catalog-card:argocd-badge', async ({ page }) => {
  await navigateWithinHost(page, '/acme/catalog');
  await expect(page.getByText('GitOps · ArgoCD', { exact: true }).first()).toBeVisible();
});

test('slot instance-action:argocd-sync', async ({ page }) => {
  await openInstanceActions(page);
  await page.getByRole('button', { name: 'Sync (ArgoCD)', exact: true }).click();
  await assertActionSuccess(page, 'argocd.sync', 'Sync (ArgoCD)', {
    clusterId: 'cl-eks-prod',
    app: { name: 'orders-db', namespace: 'argocd', project: 'inari' },
  });
});

test('slot instance-action:argocd-refresh', async ({ page }) => {
  await openInstanceActions(page);
  await page.getByRole('button', { name: 'Refresh (ArgoCD)', exact: true }).click();
  await assertActionSuccess(page, 'argocd.refresh', 'Refresh (ArgoCD)', {
    clusterId: 'cl-eks-prod',
    app: { name: 'orders-db', namespace: 'argocd', project: 'inari' },
    hard: false,
  });
});

test('slot instance-action:argocd-rollback', async ({ page }) => {
  await openInstanceActions(page);
  await page.getByRole('button', { name: 'Rollback (ArgoCD)', exact: true }).click();
  await expect(page.getByRole('alert').filter({ hasText: 'has no recorded ArgoCD revision' })).toBeVisible();
});
