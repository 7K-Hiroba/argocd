import { useEffect, useState } from 'react';
import type { ClusterTabSlotProps } from '@7k-inari/ui-plugin-sdk';
import { useTenant } from '@7k-inari/ui-plugin-sdk';
import { listClusterInstances, type ResourceInstance } from './api';
import { healthBadgeClass } from './health';

export function HealthBadge({ health }: { health?: string }) {
  const h = health ?? 'Unknown';
  return (
    <span data-testid="health-badge" className={healthBadgeClass(h)}>
      {h}
    </span>
  );
}

/** ArgoCDHealthTab is the ClusterTab contribution: ArgoCD health of every
 * resource instance on the cluster. Rendered with 7k-design-system classes —
 * the host console loads the stylesheet and tokens. */
export function ArgoCDHealthTab({ cluster }: ClusterTabSlotProps) {
  const { current } = useTenant();
  const scopeKey = current ? `${current.orgId}/${cluster.id}` : null;
  const [result, setResult] = useState<{
    scopeKey: string;
    instances: ResourceInstance[];
    error?: string;
  }>();

  useEffect(() => {
    if (!current || scopeKey === null) return;
    let cancelled = false;
    listClusterInstances(current.orgId, cluster.id)
      .then((list) => {
        if (!cancelled) setResult({ scopeKey, instances: list });
      })
      .catch((e: Error) => {
        if (!cancelled) setResult({ scopeKey, instances: [], error: e.message });
      });
    return () => {
      cancelled = true;
    };
  }, [current, cluster.id, scopeKey]);

  // Results belong to exactly one tenant/cluster scope. Keying them prevents
  // a previous tenant's inventory or error from rendering during a transition.
  const scopedResult = result?.scopeKey === scopeKey ? result : undefined;
  const instances = scopedResult?.instances ?? [];
  const error = scopedResult?.error;
  const loading = scopeKey !== null && scopedResult === undefined;

  return (
    <section className="card">
      <div className="card-header">
        <h3 className="mono-label">ArgoCD health — {cluster.name}</h3>
      </div>
      <div className="card-body">
        {error && (
          <p role="alert" className="alert alert-danger">
            Failed to load instances: {error}
          </p>
        )}
        {!error && loading && (
          <span className="skeleton" aria-busy="true" style={{ display: 'block', height: 48 }} />
        )}
        {!error && !loading && instances.length === 0 && (
          <div className="empty">
            <span className="empty-title">No instances</span>
            <span className="empty-text">No Inari-managed instances on this cluster.</span>
          </div>
        )}
        <ul className="list">
          {instances.map((i) => {
            // Server returns state/statusMessage as top-level fields.
            const top = i as unknown as Record<string, unknown>;
            const state = (top.state as string) ?? '';
            const message = (top.statusMessage as string) ?? '';
            return (
              <li key={i.id} className="list-item argocd-instance-row">
                <strong>{i.name || i.id}</strong>
                <span className="tag">{i.catalogItemId.replace(/^curated:/, '')}</span>
                {i.namespace ? <span className="tag">({i.namespace})</span> : null}
                <HealthBadge health={i.health} />
                {state && state !== 'running' ? <span className="tag">{state}</span> : null}
                {message ? (
                  <span
                    className="badge badge-warning"
                    title={message}
                    style={{ maxWidth: 360, overflow: 'hidden', textOverflow: 'ellipsis', whiteSpace: 'nowrap' }}
                  >
                    {message}
                  </span>
                ) : null}
              </li>
            );
          })}
        </ul>
      </div>
    </section>
  );
}
