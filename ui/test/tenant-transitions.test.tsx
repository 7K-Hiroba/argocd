import { beforeEach, describe, expect, it, vi } from 'vitest';
import { act, createElement } from 'react';
import { createRoot, type Root } from 'react-dom/client';
import {
  AuthProvider,
  TenantProvider,
  type AuthState,
  type ClusterTabSlotProps,
  type TenantState,
} from '@7k-inari/ui-plugin-sdk';
import { ArgoCDHealthTab } from '../src/ClusterTab';

const auth: AuthState = { principal: null, getToken: () => 'test-token' };
const cluster: ClusterTabSlotProps['cluster'] = {
  id: 'cluster-1',
  name: 'prod-1',
  tenantId: 'tenant-acme',
  state: 'Active',
};
const acme: TenantState = {
  current: { orgId: 'acme', orgName: 'Acme' },
  available: [{ orgId: 'acme', orgName: 'Acme' }],
  switchTenant: () => {},
  onTenantChange: () => () => {},
};
const instance = {
  id: 'ri-acme',
  catalogItemId: 'ci-1',
  clusterId: 'cluster-1',
  name: 'acme-private-instance',
  namespace: 'shop',
  health: 'Healthy',
  status: {},
};

function renderTree(tenant: TenantState, currentCluster = cluster) {
  return createElement(AuthProvider, {
    value: auth,
    children: createElement(TenantProvider, {
      value: tenant,
      children: createElement(ArgoCDHealthTab, { cluster: currentCluster }),
    }),
  });
}

describe('ArgoCDHealthTab tenant and cluster transitions', () => {
  let element: HTMLDivElement;
  let root: Root;

  beforeEach(() => {
    vi.restoreAllMocks();
    element = document.createElement('div');
    document.body.appendChild(element);
    root = createRoot(element);
    return () => {
      root.unmount();
      element.remove();
    };
  });

  it('does not expose the previous tenant instances after the tenant is cleared', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ instances: [instance] }), { status: 200 }),
    );
    await act(async () => root.render(renderTree(acme)));
    expect(element.innerHTML).toContain('acme-private-instance');

    await act(async () => root.render(renderTree({ ...acme, current: null })));
    expect(element.innerHTML).not.toContain('acme-private-instance');
    expect(element.innerHTML).toContain('No instances');
  });

  it('isolates inventory while a direct tenant transition is loading', async () => {
    let completeSecondLoad!: (response: Response) => void;
    vi.spyOn(globalThis, 'fetch')
      .mockResolvedValueOnce(new Response(JSON.stringify({ instances: [instance] }), { status: 200 }))
      .mockReturnValueOnce(new Promise<Response>((resolve) => { completeSecondLoad = resolve; }));
    await act(async () => root.render(renderTree(acme)));
    expect(element.innerHTML).toContain('acme-private-instance');

    const otherTenant = {
      ...acme,
      current: { orgId: 'other', orgName: 'Other' },
      available: [...acme.available, { orgId: 'other', orgName: 'Other' }],
    };
    await act(async () => root.render(renderTree(otherTenant)));
    expect(element.innerHTML).not.toContain('acme-private-instance');
    expect(element.innerHTML).toContain('aria-busy="true"');

    const otherInstance = { ...instance, id: 'ri-other', name: 'other-instance' };
    await act(async () => {
      completeSecondLoad(new Response(JSON.stringify({ instances: [otherInstance] }), { status: 200 }));
    });
    expect(element.innerHTML).toContain('other-instance');
    expect(element.innerHTML).not.toContain('acme-private-instance');
  });

  it('clears an earlier error when a subsequent cluster load succeeds', async () => {
    const fetchSpy = vi
      .spyOn(globalThis, 'fetch')
      .mockRejectedValueOnce(new Error('inventory unavailable'))
      .mockResolvedValueOnce(new Response(JSON.stringify({ instances: [instance] }), { status: 200 }));
    await act(async () => root.render(renderTree(acme)));
    expect(element.innerHTML).toContain('inventory unavailable');

    await act(async () => root.render(renderTree(acme, { ...cluster, id: 'cluster-2', name: 'prod-2' })));
    expect(fetchSpy).toHaveBeenCalledTimes(2);
    expect(element.innerHTML).not.toContain('inventory unavailable');
    expect(element.innerHTML).toContain('acme-private-instance');
  });
});
