# Fresh import: inari-ext-argocd

- Source repository: `https://github.com/7K-Inari/inari-ext-argocd`
- Source path: repository root
- Last migrated source SHA: `941b0af20ec43a07977645125c97c7759a259a68`
- Import mode: fresh file import without source Git history; no `filter-repo`, subtree, mirror, or re-pin tooling

## Destination and content inventory

- `cmd/`, `internal/`, `e2e/`, and `scripts/` → `extension/` with the same relative paths
- `Dockerfile`, `go.mod`, `go.sum`, `Makefile`, `README.md`, `CHANGELOG.md`, `LICENSE`, `CONTRIBUTING.md`, and `SECURITY.md` → `extension/`
- `ui/` → `ui/`
- `extension.yaml` → `extension.yaml`

## Deliberately excluded legacy inputs

- Source `.git/` and `.github/` directories
- Source release-please manifests/configuration, Renovate configuration, and root agent guidance superseded by the Hiroba app contract
- The source GitHub-Release UI publication path and mutable backend release tags; these are not carried forward as artifacts or release behavior

Workaround-only files present in this snapshot are removed or replaced in subsequent TDD commits. This document is an attribution record, not an assertion that the imported code was already SDK-only, released, admitted, rendered, or deployed.
