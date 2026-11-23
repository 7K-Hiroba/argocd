// Control-plane API access for the extension's UI remote. Extension actions
// use the public SDK client; the host supplies authentication and typed
// re-auth handling. Only the tenant inventory read below remains an
// ArgoCD-specific REST projection because it is not an extension invocation.

import { ApiClient, getAuthState } from '@7k-inari/ui-plugin-sdk';

export const EXTENSION_NAME = 'inari-ext-argocd';

function apiClient(): ApiClient {
  return new ApiClient({ baseUrl: '', getToken: () => getAuthState()?.getToken() });
}

export type { ResourceInstance } from '@7k-inari/ui-plugin-sdk';
import type { ResourceInstance } from '@7k-inari/ui-plugin-sdk';

export interface ActionResult {
  commandId: string;
  outcome: 'applied' | 'accepted';
  message?: string;
}

/** invokeAction calls one backend plugin action through the public SDK. */
export function invokeAction(action: string, payload: unknown): Promise<ActionResult> {
  return apiClient().invokeExtension<ActionResult>(EXTENSION_NAME, `/actions/${encodeURIComponent(action)}`, {
    method: 'POST',
    body: payload,
  });
}

/** listClusterInstances reads the host's tenant-scoped resource inventory. */
export async function listClusterInstances(
  org: string,
  clusterId: string,
): Promise<ResourceInstance[]> {
  const token = await getAuthState()?.getToken();
  const headers: Record<string, string> = { 'content-type': 'application/json' };
  if (token) headers.authorization = `Bearer ${token}`;
  const path = `/api/v1/tenants/${encodeURIComponent(org)}/instances?clusterId=${encodeURIComponent(clusterId)}`;
  const res = await fetch(path, { method: 'GET', headers });
  if (!res.ok) throw new Error(`GET ${path} failed: ${res.status}`);
  const body = (await res.json()) as { instances: ResourceInstance[] | null };
  return body.instances ?? [];
}

/** argocdAppRef extracts the ArgoCD Application identity recorded in status. */
export function argocdAppRef(instance: ResourceInstance): {
  name: string;
  namespace: string;
  project: string;
} {
  const argo = (instance.status?.argocd ?? {}) as Record<string, unknown>;
  return {
    name: (argo.name as string) ?? instance.name,
    namespace: (argo.namespace as string) ?? 'argocd',
    project: (argo.project as string) ?? 'inari',
  };
}
