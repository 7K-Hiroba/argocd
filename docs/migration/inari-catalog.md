# Fresh import: inari-catalog ArgoCD package

- Source repository: `https://github.com/7K-Inari/inari-catalog`
- Source path: `packages/argocd`
- Last migrated source SHA: `43e80dcd24de12ba7d22f6f0a91e452f40b39c58`
- Imported package version: `1.3.0`
- Import mode: fresh file import without source Git history; no `filter-repo`, subtree, mirror, old-channel maintenance, or re-pin tooling

## Destination and content inventory

- `packages/argocd/package.yaml` → `packages/argocd/package.yaml`
- `packages/argocd/chart.yaml` → `packages/argocd/chart.yaml`
- `packages/argocd/values-defaults.yaml` → `packages/argocd/values-defaults.yaml`
- `packages/argocd/README.md` → `packages/argocd/README.md`
- `packages/argocd/CHANGELOG.md` → `packages/argocd/CHANGELOG.md`
- `packages/argocd/tests/run.sh` → `packages/argocd/tests/run.sh`

The source package had no `values.schema.json`, UI hints, canonical observability manifest, first-party OCI chart, or scaffold template. Those are added as new Hiroba content in subsequent commits rather than being presented as migrated source content. The source package remains unchanged in `inari-catalog`; removal and deprecation belong to the later W5 migration task.
