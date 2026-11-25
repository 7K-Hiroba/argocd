# Architecture

## Components

- `extension/`: Go backend with subprocess and HTTP dial-mode entrypoints. Actions use the public `host/agentgateway` client and generated inari-api Connect contract. `oidc-sso-session` supplies per-user ArgoCD credentials through the existing vault hop; no tenant credential is stored by this repository.
- `ui/`: Module Federation remote built by `inari-ui-ext`. It contributes the existing cluster tab, catalog badge, and sync/refresh/rollback instance actions through the public SDK.
- `charts/argocd`: thin first-party OCI wrapper around locked upstream `argo-cd` 7.6.12. Generated observability output is separate from scrape topology.
- `packages/argocd`: platform-app descriptor for the OCI chart, schema validation, and exact observability reference.
- `packages/observability-pack`: generated v2 KRO RGD for folder metadata, dashboards, and central/local rule selection.
- `templates/argocd-application`: pure native Application scaffold with no deploy-time side effects.
- `tests/smoke`: one artifact-bound read-only Argo Workflows suite for CI and separately gated environments.

## Action request flow

Browser → authenticated extension proxy → backend action → typed Agent Gateway → existing outbound agent session → tenant-local ArgoCD. The tenant is derived only from `AuthContext`, configured projects are enforced on both extension and agent paths, and disconnected agents fail closed. Raw downstream tokens never enter command payloads, logs, audits, or errors. Typed SDK re-auth integrates with the host's retry-once behavior.

## GitOps boundaries

Catalog and companion rendering use the existing source-agnostic renderer and resolved baseline repository after sibling Phase 3. Native `Application`/`AppProject` resources replace any operator facade. This repository does not implement a resolver, renderer, admission verifier, Grafana client, environment verifier, or single-repo fallback.
