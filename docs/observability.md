# Observability

`observability/` is the only editable source for dashboards, rules, metadata, and fixtures. `argocd-overview` uses the registered UID/folder and `${DS_PROMETHEUS}`; link variables are limited to `namespace`. Central rules aggregate degraded and out-of-sync applications while preserving the post-Alloy labels. Local rules check Argo CD cluster connection health.

Run `make generate` after canonical changes. Never hand-edit generated `charts/argocd/inari/observability.yaml`, `charts/argocd/templates/observability/*`, `packages/argocd/observability.yaml`, or `packages/observability-pack/rgd.yaml`. `make check-generated` must be diff-free. Release checks compare packaged output as well.

The fixtures in `observability/tests/alloy/` use the exact monitoring `inari.monitoring.post-alloy-labels/v1` contract: required `org`, `cluster`, `job`, and `instance` labels, with optional `namespace`. They are representative inputs, not a claim of deployed telemetry or readiness. Promtool tests cover degraded/out-of-sync and no-series paths.

For Inari-managed deployments, chart-native observability is disabled to avoid duplicate resources. After sibling Phase 3, an exact compatible companion selects `dashboards + rules-central` for central or `dashboards + rules-local` for local. Renderer-created baseline PRs and target-side verification are required before links appear. No embedded dashboard, folder ACL/team sync, direct Grafana write, raw observability bundle, puller, or PR bot exists here.
