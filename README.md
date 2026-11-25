# ArgoCD for Hiroba

Self-contained ArgoCD application repository: public-SDK Go extension, Module Federation UI, locked OCI chart, catalog packages, native Application scaffold, canonical observability, and artifact-bound runtime verification. It is a **non-stack** addon/extension repository and deliberately has no stack lifecycle lane.

## Layout

- `extension/`: ArgoCD sync, refresh, rollback, and resource actions using the public Agent Gateway and typed re-auth SDKs.
- `ui/`: existing cluster tab, catalog badge, and instance actions built with `inari-ui-ext` and OCI-published.
- `charts/argocd`: thin OCI wrapper around locked upstream `argo-cd` 7.6.12; no forked controller templates.
- `packages/argocd`: platform-app descriptor and values schema for that chart.
- `packages/observability-pack`: generated v2 KRO companion for central/local assets.
- `templates/argocd-application`: pure native Application scaffold; no deployment at scaffold time.
- `observability/`: canonical dashboards, rules, metadata, and reviewed post-Alloy fixtures.
- `tests/smoke`: actual Argo CD health/version journey with read-only ci/stage/prod declarations.

See `docs/architecture.md`, `docs/observability.md`, `docs/verification.md`, and `docs/release-process.md` for component contracts. Exact source and candidate dependency provenance is in `docs/migration/` and `docs/dependencies.md`. Current candidate inputs are not represented as merged, signed, rendered, or deployed evidence.

## Develop

```sh
make test-backend
make test-ui
make validate
make check-generated
make check-workarounds
```

`WF` must point at the exact workflow checkout used for local generation/validation. Hosted thin callers are pinned to full SHAs and are authoritative for package, chart, compatibility, browser, and runtime lanes. Generated observability and indexes must remain diff-free.

Deployable artifacts always use explicit versions/digests and source-qualified identity:

```sh
inari deploy <source>/argocd --cluster <cluster-id> --file values.yaml
```

No GitHub-Release UI remote, classic Helm index, mutable latest/edge deployment, SDK postinstall, private gateway proto, custom webpack build, or handwritten re-auth mapping is supported.

## Release and dogfood gates

All merges, GHCR rehearsals, signed releases, registry activation, renderer/baseline changes, deployed stage/prod runs, and production actions require their applicable separate approvals. Sibling Phase 3 must be independently confirmed before companion or Workflow rendering. Current open verification work is tracked in `docs/reference-promotion.md`; no such evidence is claimed by this README.
