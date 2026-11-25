# argocd (platform app)

Argo CD hosts tenant-local GitOps root applications and native AppProjects. The source-agnostic renderer emits ArgoCD `Application`/`ApplicationSet` and native resources through the first-class addon flow; no inari-operator or provider controller is required.

The package points to the first-party OCI chart `oci://ghcr.io/7k-hiroba/argocd/charts/argocd`, which wraps the exact upstream `argo-cd` chart `7.6.12` (Argo CD `v2.12.6`) through `Chart.lock`. Defaults preserve HA server/repo-server/application-set, redis-ha, metrics/ServiceMonitors, default read-only RBAC, and disabled Dex/notifications. SSO and secret material remain delivered through the existing Inari baseline and ESO machinery.

The package carries a day-one values schema and an exact compatible `observability-pack` reference. Inari-managed deployments leave chart-native observability disabled and select a central or local companion explicitly.

```sh
helm dependency build charts/argocd
helm template argocd charts/argocd \
  --namespace argocd \
  --api-versions monitoring.coreos.com/v1 \
  --values packages/argocd/values-defaults.yaml
```
