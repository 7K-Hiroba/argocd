# Argo CD chart

This first-party chart wraps the exact upstream `argo-cd` chart `7.6.12` (Argo CD `v2.12.6`) through `Chart.lock`; it does not fork or mirror upstream templates. It publishes generated Inari observability metadata and optional generated dashboard/rule resources from the canonical `observability/` tree.

The upstream values live under the `argo-cd` key. Defaults preserve the imported HA, redis-ha, metrics/ServiceMonitor, read-only RBAC, and disabled Dex/notification choices. ServiceMonitors are deployment topology, not canonical observability assets.

`observability.enabled` and `observability.rules.enabled` default to `false` so Inari-managed deployments do not double-provision resources when a central or local companion is selected. Standalone Helm users may opt in after installing the matching provisioning dependencies.

The chart is published only as an OCI artifact. No classic Helm index or mutable deployment tag is used.
